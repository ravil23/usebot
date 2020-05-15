package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandSelectSubject = "Предмет"
	commandSelectLevel   = "Сложность"
	commandSkip          = "Пропустить"
	commandStart         = "start"

	labelAnswered = "answered"

	textSelectSubject = "Выберите предмет"
	textSelectLevel   = "Выберите сложность"

	alertsChatID = -1001436548831

	botName = "GIA11Bot"
)

var userSelectedSubject = map[int]string{}
var userSelectedLevel = map[int]string{}
var userChat = map[int]int64{}

func getBotTokenOrPanic() string {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Panic("bot token is empty")
	}
	return botToken
}

func formatUserString(tgUser *tgbotapi.User) string {
	userString := fmt.Sprintf("[id=%d]", tgUser.ID)
	if tgUser.UserName != "" {
		userString = "@" + tgUser.UserName + " " + userString
	}
	if tgUser.LastName != "" {
		userString = tgUser.LastName + " " + userString
	}
	if tgUser.FirstName != "" {
		userString = tgUser.FirstName + " " + userString
	}
	return userString
}
