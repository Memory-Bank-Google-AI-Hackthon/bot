package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

type MessageSave struct {
	Urls    []string `json:"urls"`
	Images  []string `json:"images"`
	Message string   `json:"message"`
}

var records = NewRecords()

func NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(message.Content, "!save"):

		discord.ChannelTyping(message.ChannelID)

		msg := GetSaveMessage(message.Content)

		if message.Attachments != nil && len(message.Attachments) > 0 {
			images := make([]string, len(message.Attachments))
			for i, attachment := range message.Attachments {
				images[i] = attachment.URL
			}
			msg.Images = images
		}

		_, err := json.Marshal(msg)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error saving message")
			return
		}

		summary := getSummary(msg)

		discord.ChannelMessageSend(message.ChannelID, summary)
	case strings.HasPrefix(message.Content, "!take"):
		discord.ChannelTyping(message.ChannelID)
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	case strings.HasPrefix(message.Content, "!records"):
		records := records.GetRecords()
		if len(records) == 0 {
			discord.ChannelMessageSend(message.ChannelID, "No records found")
			return
		}

		for _, record := range records {
			discord.ChannelMessageSend(message.ChannelID, record.Message)
		}
	case strings.HasPrefix(message.Content, "!clear"):
		records.Clear()
		discord.ChannelMessageSend(message.ChannelID, "Records cleared")
	case strings.HasPrefix(message.Content, "!summarize"):
		discord.ChannelTyping(message.ChannelID)
		summaries, err := GetGeminiSummary(records.GetRecords())
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error summarizing")
			return
		}

		records.Clear()

		for _, summary := range summaries {
			discord.ChannelMessageSend(message.ChannelID, summary)
		}

	default:
		records.AddRecord(Record{
			UserId:   message.Author.ID,
			UserName: message.Author.GlobalName,
			Message:  message.Content,
		})

		return
	}
}

func getUrls(message string) []string {

	rxStrict := xurls.Strict()
	urls := rxStrict.FindAllString(message, -1)

	return urls
}

func GetSaveMessage(message string) *MessageSave {
	msg := &MessageSave{
		Message: message,
	}

	urls := getUrls(msg.Message)

	for _, url := range urls {
		msg.Message = strings.ReplaceAll(msg.Message, url, "")
	}

	msg.Message = strings.TrimPrefix(msg.Message, "!save")

	msg.Urls = urls

	return msg
}

func getSummary(message *MessageSave) string {
	jsonString, err := json.Marshal(message)
	if err != nil {
		return "Error parsing message"
	}

	res, err := http.Post("https://us-central1-memory-bank-423810.cloudfunctions.net/chatbot-test", "application/json", bytes.NewReader(jsonString))
	if err != nil {
		return "Error getting summary"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "Error reading response"
	}

	return string(body)
}
