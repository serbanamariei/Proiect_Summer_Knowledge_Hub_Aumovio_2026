package main

import (
	"bytes"
	"comun"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kkdai/youtube/v2"
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
			go proceseazaRutare(bot, update.Message)
		}
	}
}

func proceseazaRutare(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if strings.Contains(message.Text, "youtube.com") || strings.Contains(message.Text, "youtu.be") {
		proceseazaYouTube(bot, message)
		return
	}

	proceseazaScraping(bot, message)
}

func proceseazaYouTube(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msgAsteptare := tgbotapi.NewMessage(message.Chat.ID, "Se descarca video ul...")
	bot.Send(msgAsteptare)

	client := youtube.Client{}
	video, err := client.GetVideo(message.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare la citirea video ului: "+err.Error()))
		return
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Nu s a gasit un format valid cu audio"))
		return
	}

	format := formats[0]

	stream, _, err := client.GetStream(video, &format)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare la preluarea stream ului"))
		return
	}
	defer stream.Close()

	numeFisier := fmt.Sprintf("video_temp_%d.mp4", message.Chat.ID)
	file, err := os.Create(numeFisier)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare la salvarea temp ului"))
		return
	}

	_, err = io.Copy(file, stream)
	file.Close()

	defer os.Remove(numeFisier)

	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare la salvarea video ului pe disc"))
		return
	}

	info, err := os.Stat(numeFisier)
	if err == nil && info.Size() > 49*1024*1024 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Video ul depaseste 50 mb si deci nu poate fi trimis pe Telegram"))
		return
	}

	msgVideo := tgbotapi.NewVideo(message.Chat.ID, tgbotapi.FilePath(numeFisier))
	msgVideo.Caption = "Titlu video:" + video.Title

	_, err = bot.Send(msgVideo)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Eroare la trimiterea video ului pe Telegram"))
	}
}

func proceseazaScraping(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
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

		limitaAfsare := 15
		if len(raspuns.Linkuri) < limitaAfsare {
			limitaAfsare = len(raspuns.Linkuri)
		}

		for i := 0; i < limitaAfsare; i++ {
			textMesaj += "- " + raspuns.Linkuri[i] + "\n"
		}

		if len(raspuns.Linkuri) > 15 {
			textMesaj += "\n... mai sunt si alte link uri"
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, textMesaj)
	bot.Send(msg)
}
