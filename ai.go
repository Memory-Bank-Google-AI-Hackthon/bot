package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const tmpl = `{{range .}}
User: {{ .UserName }}, Message: {{ .Message }}\n
{{end }}`

func GetGeminiSummary(record []Record) ([]string, error) {
	ctx := context.Background()
	key := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.0-pro")
	cs := model.StartChat()

	send := func(msg string) (*genai.GenerateContentResponse, error) {
		fmt.Printf("== Me: %s\n== Model:\n", msg)
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	t, err := template.New("example").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, records.GetRecords())
	if err != nil {
		return nil, err
	}

	result := buf.String()

	res, err := send(fmt.Sprintf(`
		Here is some records of the conversation so far:
		The format is a JSON string with the structure userName and message.
		[%s]
		Please summarize the above information and response with traditional chinese.
	`, result))

	if err != nil {
		return nil, err
	}

	var summaries []string

	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				t, ok := part.(genai.Text)
				if ok {
					summaries = append(summaries, string(t))
				}
			}
		}
	}

	return summaries, nil
}
