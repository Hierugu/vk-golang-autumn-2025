package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var BotToken = "XXX"
var WebhookURL = "https://525f2cb5.ngrok.io"
var bot *tgbotapi.BotAPI

func processUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	matches := commandRe.FindStringSubmatch(update.Message.Text)
	if len(matches) > 0 {
		cmd := matches[0]
		if handler, ok := handlers[cmd]; ok {
			handler(update.Message)
		}
	}
}

func startTaskBot(ctx context.Context) error {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	panicIf(err)
	log.Println("Bot:", bot)
	bot.Debug = true
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates := bot.GetUpdatesChan(u)
	// updates := bot.ListenForWebhook(WebhookURL)

	log.Printf("Authorized on account %s", bot.Self.UserName)
	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}
	updates := bot.ListenForWebhook("/")

	port := "8081"
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	for update := range updates {
		processUpdate(&update)
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	panicIf(err)
}
