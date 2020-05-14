package entity

import (
	"fmt"
	"math/rand"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ExplanationPrefix = "Правильный ответ: "

	LevelLow    = 1
	LevelMedium = 2
	LevelHigh   = 3
)

type Task struct {
	ID      int               `json:"id"`
	Level   int               `json:"level"`
	Text    string            `json:"text"`
	Doc     string            `json:"doc"`
	Answer  string            `json:"answer"`
	Options map[string]string `json:"options"`
	Themes  []string          `json:"themes"`
}

func (t *Task) MakeTelegramPoll(chatID int64) *tgbotapi.SendPollConfig {
	if len(t.Text) > 255 {
		return nil
	}
	var correctOptionID int64 = -1
	tgOptions := make([]string, 0, len(t.Options))
	for i, key := range t.shuffledOptionKeys() {
		if key == t.Answer {
			correctOptionID = int64(i)
		}
		option := t.Options[key]
		if len(option) > 100 {
			return nil
		}
		tgOptions = append(tgOptions, option)
	}

	tgPoll := tgbotapi.NewPoll(chatID, fmt.Sprintf("#%d\n%s", t.ID, t.Text), tgOptions...)
	tgPoll.Explanation = ExplanationPrefix + t.Options[t.Answer]
	tgPoll.CorrectOptionID = correctOptionID
	tgPoll.Type = "quiz"
	tgPoll.IsAnonymous = false
	return &tgPoll
}

func (t *Task) MakeTelegramMessage(chatID int64) *tgbotapi.MessageConfig {
	tgMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("<b>#%d\n%s</b>\n", t.ID, t.Text))
	if t.Doc != "" {
		tgMessage.Text += "\n" + t.Doc + "\n"
	}
	tgMessage.ParseMode = tgbotapi.ModeHTML
	tgButtons := make([]tgbotapi.InlineKeyboardButton, len(t.Options))
	for i, key := range t.shuffledOptionKeys() {
		option := t.Options[key]
		index := i + 1
		tgMessage.Text += fmt.Sprintf("\n%d. %s", index, option)
		tgButtons[i] = tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(index), fmt.Sprintf("%d:%t", index, key == t.Answer))
	}
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgButtons)
	return &tgMessage
}

func (t *Task) shuffledOptionKeys() []string {
	keys := make([]string, 0, len(t.Options))
	for key := range t.Options {
		keys = append(keys, key)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}
