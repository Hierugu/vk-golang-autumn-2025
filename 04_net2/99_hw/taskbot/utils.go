package main

import (
	"log"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func panicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func replyMsg(msg *tgbotapi.Message, text string) tgbotapi.MessageConfig {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID
	return reply
}

func sendMsg(msg tgbotapi.MessageConfig) {
	_, err := bot.Send(msg)
	panicIf(err)
}
