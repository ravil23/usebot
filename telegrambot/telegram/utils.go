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
	commandNext          = "Продолжить"
	commandStart         = "start"

	labelAnswered = "answered"

	textExplanation   = "Пояснение"
	textSelectSubject = "Список доступных предметов"
	textSelectLevel   = "Варианты уровней сложности"

	AlertsChatID = -1001436548831

	Bot11Name = "GIA11Bot"
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

func getWelcomeText(tgUser *tgbotapi.User) string {
	userString := formatUserStringPretty(tgUser)
	var welcomeText string
	if userString == "" {
		welcomeText = "Привет, любознательный незнакомец!"
	} else {
		welcomeText = fmt.Sprintf("Привет, %s!", userString)
	}
	welcomeText += "\nВыбери предмет, чтобы начать подготовку."
	return welcomeText
}

func formatUserStringPretty(tgUser *tgbotapi.User) string {
	userString := ""
	if tgUser.LastName != "" {
		userString = tgUser.LastName + " " + userString
	}
	if tgUser.FirstName != "" {
		userString = tgUser.FirstName + " " + userString
	}
	return userString
}

func formatUserStringVerbose(tgUser *tgbotapi.User) string {
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
