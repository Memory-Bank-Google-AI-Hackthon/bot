package main

import (
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

type MessageSave struct {
	Urls    []string `json:"urls"`
	Images  []string `json:"images"`
	Message string   `json:"message"`
}

func NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == discord.State.User.ID {
		return
	}

	discord.ChannelTyping(message.ChannelID)

	switch {
	case strings.HasPrefix(message.Content, "!save"):
		msg := GetSaveMessage(message.Content)

		if message.Attachments != nil && len(message.Attachments) > 0 {
			images := make([]string, len(message.Attachments))
			for i, attachment := range message.Attachments {
				images[i] = attachment.URL
			}
			msg.Images = images
		}

		jsonString, err := json.Marshal(msg)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error saving message")
			return
		}

		discord.ChannelMessageSend(message.ChannelID, string(jsonString))
	case strings.HasPrefix(message.Content, "!take"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	default:
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
