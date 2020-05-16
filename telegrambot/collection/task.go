package collection

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ExplanationPrefix = "Правильный ответ: "
)

type Level int

const (
	LevelLow    = Level(1)
	LevelMedium = Level(2)
	LevelHigh   = Level(3)
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "Базовая"
	case LevelMedium:
		return "Повышенная"
	case LevelHigh:
		return "Высокая"
	default:
		return ""
	}
}

type Task struct {
	ID           int               `json:"id"`
	Level        Level             `json:"level"`
	Text         string            `json:"text"`
	Doc          string            `json:"doc"`
	Answer       string            `json:"answer"`
	Options      map[string]string `json:"options"`
	Themes       []string          `json:"themes"`
	Requirements []string          `json:"requirements"`
	SendAsPoll   bool              `json:"sendAsPoll"`
	SubjectName  string            `json:"subjectName"`
}

func (t *Task) String() string {
	if data, err := json.Marshal(t); err == nil {
		return string(data)
	} else {
		return fmt.Sprintf("Task{ID: %d}", t.ID)
	}
}

func (t *Task) MakeTelegramPoll(chatID int64) *tgbotapi.SendPollConfig {
	log.Printf("Make poll to chat %d from task: %s", chatID, t)
	var correctOptionID int64 = -1
	tgOptions := make([]string, 0, len(t.Options))
	for i, key := range t.shuffledOptionKeys() {
		if key == t.Answer {
			correctOptionID = int64(i)
		}
		option := t.Options[key]
		tgOptions = append(tgOptions, option)
	}
	tgPoll := tgbotapi.NewPoll(chatID, t.getTextWithSubject(), tgOptions...)
	tgPoll.Explanation = ExplanationPrefix + t.Options[t.Answer]
	tgPoll.CorrectOptionID = correctOptionID
	tgPoll.Type = "quiz"
	tgPoll.IsAnonymous = false
	return &tgPoll
}

func (t *Task) MakeTelegramMessage(chatID int64) *tgbotapi.MessageConfig {
	log.Printf("Make message to chat %d from task: %s", chatID, t)
	tgMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("<b>%s</b>\n", t.getTextWithSubject()))
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

func (t *Task) getTextWithSubject() string {
	return fmt.Sprintf("%s\n%s", t.SubjectName, t.Text)
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
