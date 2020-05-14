package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ravil23/usebot/telegrambot/collection/entity"
)

const (
	retryPeriod     = time.Second
	maxRetriesCount = 30
	timeoutSeconds  = 60
	listenersPoolSize = 10

	commandSelectSubject = "Выбрать предмет"
	commandSkip = "Пропустить"
	commandStart = "start"

	labelAnswered = "answered"

	textSelectSubject = "Что поизучаем?"

	subjectRussian = "Русский язык"
	subjectHistory = "История"
)

var userSelectedSubject = map[int]string{}
var userChat = map[int]int64{}

func GetBotTokenOrPanic() string {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Panic("bot token is empty")
	}
	return botToken
}

type Bot struct {
	api *tgbotapi.BotAPI
	database *entity.Database
}

func NewBot(database *entity.Database) *Bot {
	return &Bot{
		database: database,
	}
}

func (b *Bot) Init() {
	log.Printf("Bot is initializing...")
	botToken := GetBotTokenOrPanic()
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= maxRetriesCount; i++ {
		if api, err := tgbotapi.NewBotAPI(botToken); err != nil {
			log.Printf("Attempt %d failed: %v", i, err)
			time.Sleep(retryPeriod)
		} else {
			b.api = api
			log.Printf("Bot successfully initialized")
			return
		}
	}
	log.Panic("max retries count exceeded")
}

func (b *Bot) HealthCheck() {
	go func() {
		address := ":8080"
		path := "/healthcheck"
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s request to %s%s with User-Agent: %s", r.Method, r.Host, r.URL, r.UserAgent())
			_, _ = fmt.Fprint(w, `{"status": "ok"}`)
		})
		log.Printf("Listening health check on address %s%s", address, path)
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Panic(err)
		}
	}()
}

func (b *Bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = timeoutSeconds
	updates := b.api.GetUpdatesChan(updateConfig)

	wg:= sync.WaitGroup{}
	listener := func () {
		for update := range updates {
			if update.Message != nil {
				b.handleMessage(update.Message)
			} else if update.CallbackQuery != nil {
				b.handleCallbackQuery(update.CallbackQuery)
			} else if update.PollAnswer != nil {
				b.handlePollAnswer(update.PollAnswer)
			} else {
				continue
			}
		}
		wg.Done()
	}

	for i := 0; i < listenersPoolSize; i++ {
		wg.Add(1)
		go listener()
	}
	log.Printf("Telegram listeners pool size: %d", listenersPoolSize)
	wg.Wait()
}

func (b *Bot) handleMessage(tgMessage *tgbotapi.Message) {
	chatID := tgMessage.Chat.ID
	messageID := tgMessage.MessageID

	userChat[tgMessage.From.ID] = chatID

	if tgMessage.Command() == commandStart {
		b.sendStartMenu(chatID, messageID)
		b.sendSubjectsList(chatID, messageID)
	} else if tgMessage.Text == commandSelectSubject{
		b.sendSubjectsList(chatID, messageID)
	} else {
		b.sendNextTask(chatID, messageID, tgMessage.From.ID)
	}
}

func (b *Bot) sendStartMenu(chatID int64, messageID int) {
	tgChattable := b.getStartMenu(chatID)
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
	}
}

func (b *Bot) getStartMenu(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(
		chatID,
		fmt.Sprintf(`Выбери предмет, чтобы начать подготовку. Его всегда можно сменить, нажав на кнопку "%s".`, commandSelectSubject),
	)
	tgButtons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(commandSelectSubject),
		tgbotapi.NewKeyboardButton(commandSkip),
	}
	tgMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgButtons)
	return &tgMessage
}

func (b *Bot) handleCallbackQuery(tgCallbackQuery *tgbotapi.CallbackQuery) {
	chatID := tgCallbackQuery.Message.Chat.ID
	messageID := tgCallbackQuery.Message.MessageID

	userChat[tgCallbackQuery.From.ID] = chatID

	if tgCallbackQuery.Message.Text == textSelectSubject {
		if b.selectSubject(tgCallbackQuery) {
			b.sendNextTask(chatID, messageID, tgCallbackQuery.From.ID)
		}
	} else if b.updateInlineQuestion(tgCallbackQuery) {
		b.sendNextTask(chatID, messageID, tgCallbackQuery.From.ID)
	}
}

func (b *Bot) selectSubject(callbackQuery *tgbotapi.CallbackQuery) bool {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	userSelectedSubject[callbackQuery.From.ID] = callbackQuery.Data

	alreadyAnswered := callbackQuery.Data == labelAnswered
	callbackText := ""
	if alreadyAnswered {
		callbackText = fmt.Sprintf(`Для смены предмета, воспользуйся кнопкой "%s"`, commandSelectSubject)
	} else {
		tgKeyboardUpdate := tgbotapi.NewEditMessageText(chatID, messageID, callbackQuery.Message.Text)
		tgRows := make([][]tgbotapi.InlineKeyboardButton, 0, len(callbackQuery.Message.ReplyMarkup.InlineKeyboard))
		for _, row := range callbackQuery.Message.ReplyMarkup.InlineKeyboard {
			tgButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(row))
			for _, button := range row {
				if button.CallbackData == nil {
					continue
				}
				tgButton := tgbotapi.NewInlineKeyboardButtonData(button.Text, labelAnswered)
				if *button.CallbackData == callbackQuery.Data {
					tgButton.Text += " 📖️"
				}
				tgButtons = append(tgButtons, tgButton)
			}
			tgRows = append(tgRows, tgbotapi.NewInlineKeyboardRow(tgButtons...))
		}
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
		tgKeyboardUpdate.ReplyMarkup = &tgKeyboard
		if _, err := b.api.Send(tgKeyboardUpdate); err != nil {
			b.sendAlert(chatID, err.Error(), messageID)
		}
	}

	tgCallback := tgbotapi.NewCallback(callbackQuery.ID, callbackText)
	if _, err := b.api.Request(tgCallback); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
		return false
	}
	return !alreadyAnswered
}

func (b *Bot) updateInlineQuestion(callbackQuery *tgbotapi.CallbackQuery) bool {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	alreadyAnswered := false
	callbackText := ""
	if callbackQuery.Data == labelAnswered {
		alreadyAnswered = true
		callbackText = "К сожалению изменить ответ нельзя"
	} else if strings.HasPrefix(callbackQuery.Data, entity.ExplanationPrefix){
		alreadyAnswered = true
		callbackText = callbackQuery.Data
	} else {
		tgKeyboardUpdate := tgbotapi.NewEditMessageText(chatID, messageID, callbackQuery.Message.Text)
		endQuestionIndex := strings.Index(tgKeyboardUpdate.Text, "\n\n")
		if endQuestionIndex >= 0 {
			tgKeyboardUpdate.Text = "<b>" + strings.Replace(tgKeyboardUpdate.Text, "\n\n", "</b>\n\n", 1)
		} else {
			tgKeyboardUpdate.Text = "<b>" + tgKeyboardUpdate.Text + "</b>"
		}
		tgKeyboardUpdate.ParseMode = tgbotapi.ModeHTML

		tgRows := make([][]tgbotapi.InlineKeyboardButton, 0, len(callbackQuery.Message.ReplyMarkup.InlineKeyboard))
		correctOptionText := "?"
		hasMistake := false
		for _, row := range callbackQuery.Message.ReplyMarkup.InlineKeyboard {
			tgButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(row))
			for _, button := range row {
				if button.CallbackData == nil {
					continue
				}
				data := strings.Split(*button.CallbackData, ":")
				if len(data) != 2 {
					continue
				}
				tgButton := tgbotapi.NewInlineKeyboardButtonData(button.Text, labelAnswered)
				if data[1] == "true" {
					correctOptionText = tgButton.Text
					tgButton.Text += " ✅"
				} else if *button.CallbackData == callbackQuery.Data {
					tgButton.Text += " ❌"
					hasMistake = true
				}
				tgButtons = append(tgButtons, tgButton)
			}
			tgRows = append(tgRows, tgbotapi.NewInlineKeyboardRow(tgButtons...))
		}
		if hasMistake {
			callbackText = entity.ExplanationPrefix + correctOptionText
		}
		tgRows = append(
			tgRows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Пояснение",
					entity.ExplanationPrefix + correctOptionText,
				),
			),
		)
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
		tgKeyboardUpdate.ReplyMarkup = &tgKeyboard
		if _, err := b.api.Send(tgKeyboardUpdate); err != nil {
			b.sendAlert(chatID, err.Error(), messageID)
		}
	}

	tgCallback := tgbotapi.NewCallback(callbackQuery.ID,  callbackText)
	if _, err := b.api.Request(tgCallback); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
		return false
	}
	return !alreadyAnswered
}

func (b *Bot) handlePollAnswer(tgPollAnswer *tgbotapi.PollAnswer) {
	userID := tgPollAnswer.User.ID
	chatID := userChat[userID]

	b.sendNextTask(chatID, 0, userID)
}

func (b *Bot) sendSubjectsList(chatID int64, messageID int) {
	tgChattable := b.getSubjectsList(chatID)
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
	}
}

func (b *Bot) getSubjectsList(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(chatID, textSelectSubject)
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(subjectRussian, subjectRussian)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(subjectHistory, subjectHistory)},
	)
	return &tgMessage
}

func (b *Bot) sendNextTask(chatID int64, messageID int, userID int) {
	subject := userSelectedSubject[userID]
	switch subject {
	case subjectRussian:
		b.sendNextRussianTask(chatID, messageID)
	case subjectHistory:
		b.sendNextHistoryTask(chatID, messageID)
	default:
		b.sendSubjectsList(chatID, messageID)
	}
}

func (b *Bot) sendNextRussianTask(chatID int64, messageID int) {
	tgChattable := b.getNextTask(b.database.Russian.Tasks, chatID)
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
	}
}

func (b *Bot) sendNextHistoryTask(chatID int64, messageID int) {
	tgChattable := b.getNextTask(b.database.History.Tasks, chatID)
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(chatID, err.Error(), messageID)
	}
}

func (b *Bot) getNextTask(tasks []*entity.Task, chatID int64) tgbotapi.Chattable {
	task := tasks[rand.Intn(len(tasks))]
	tgPoll := task.MakeTelegramPoll(chatID)
	if tgPoll != nil {
		return tgPoll
	}
	return task.MakeTelegramMessage(chatID)
}

func (b *Bot) sendAlert(chatID int64, text string, messageID int) {
	log.Print(text)
	tgMessage := tgbotapi.NewMessage(chatID, text)
	tgMessage.ReplyToMessageID = messageID
	_, err := b.api.Send(tgMessage)
	if err != nil {
		log.Printf("Error on sending alert: %s", err)
	}
}
