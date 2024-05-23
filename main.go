package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	discord, err := discordgo.New("Bot " + config.DISCORD_BOT_SECRET)
	if err != nil {
		log.Fatal(err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(NewMessage)
	defer discord.Close()

	bot, err := linebot.New(
		config.LineClientID,
		config.LineClientSecret,
	)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					msg := GetSaveMessage(message.Text)
					if msg == nil {
						return
					}

					msgJson, err := json.Marshal(msg)
					if err != nil {
						log.Print(err)
						return
					}

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(string(msgJson))).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
