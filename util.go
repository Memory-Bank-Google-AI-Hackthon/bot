package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LineClientID       string
	LineClientSecret   string
	DISCORD_BOT_SECRET string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client_id := os.Getenv("LINE_CLIENT_ID")
	client_secret := os.Getenv("LINE_CLIENT_SECRET")
	discord_secret := os.Getenv("DISCORD_BOT_SECRET")

	if client_id == "" || client_secret == "" || discord_secret == "" {
		return nil, fmt.Errorf("LINE_CLIENT_ID, LINE_CLIENT_SECRET, DISCORD_BOT_SECRET must be set")
	}

	return &Config{
		LineClientID:       client_id,
		LineClientSecret:   client_secret,
		DISCORD_BOT_SECRET: discord_secret,
	}, nil
}
