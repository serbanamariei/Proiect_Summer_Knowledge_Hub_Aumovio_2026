package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"comun"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	webhookConfig, err := tgbotapi.NewWebhook(os.Getenv("WEBHOOK_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.Request(webhookConfig)
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe("0.0.0.0:8080", nil)

	for update := range updates {
		if update.Message != nil {
			go proceseazaMesaj(bot, update.Message)
		}
	}
}

func proceseazaMesaj(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	cerere := comun.CerereInitala{
		ChatID: message.Chat.ID,
		URL:    message.Text,
	}

	body, _ := json.Marshal(cerere)
	resp, err := http.Post("http://dispatcher_aplicatie:8082/dispatch", "application/json", bytes.NewBuffer(body))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare: Dispatcher offline."))
		return
	}
	defer resp.Body.Close()

	var raspuns comun.RaspunsCrawler
	json.NewDecoder(resp.Body).Decode(&raspuns)

	textMesaj := raspuns.Mesaj
	if len(raspuns.Linkuri) > 0 {
		textMesaj += "\n\nLink-urile gasite:\n"
		for _, link := range raspuns.Linkuri {
			textMesaj += "- " + link + "\n"
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, textMesaj)
	bot.Send(msg)
}
