package main

import (
	"fmt"
	"os"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
	WebhookURL    string
}

func IncarcaConfig() (*Config, error) {
	token := os.Getenv("TELEGRAM_TOKEN")
	dbUrl := os.Getenv("DATABASE_URL")
	webhookUrl := os.Getenv("WEBHOOK_URL")

	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_TOKEN nu este setat in mediul curent")
	}

	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL nu este setat in mediul curent")
	}

	if webhookUrl == "" {
		return nil, fmt.Errorf("WEBHOOK_URL nu este setat in mediul curent")
	}

	return &Config{
		TelegramToken: token,
		DatabaseURL:   dbUrl,
		WebhookURL:    webhookUrl,
	}, nil
}
