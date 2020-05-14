package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ravil23/usebot/telegrambot/collection/entity"
)

const (
	initializationRetryPeriod     = time.Second
	initializationMaxRetriesCount = 30
	timeoutSeconds                = 60
	listenersPoolSize             = 10
)

type Bot struct {
	hostName string
	api      *tgbotapi.BotAPI
	database *entity.Database
}

func NewBot(database *entity.Database) *Bot {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown_host"
	}
	return &Bot{
		hostName: hostName,
		database: database,
	}
}

func (b *Bot) Init() {
	log.Printf("Bot is initializing...")
	botToken := getBotTokenOrPanic()
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= initializationMaxRetriesCount; i++ {
		if api, err := tgbotapi.NewBotAPI(botToken); err != nil {
			log.Printf("Attempt %d failed: %v", i, err)
			time.Sleep(initializationRetryPeriod)
		} else {
			b.api = api
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
	log.Printf("Bot is running...")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = timeoutSeconds
	updates := b.api.GetUpdatesChan(updateConfig)

	listener := func() {
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
	}

	for i := 0; i < listenersPoolSize; i++ {
		go listener()
	}

	b.serve()
}

func (b *Bot) serve() {
	b.sendAlert(fmt.Sprintf("@%s started", botName))
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	b.sendAlert(fmt.Sprintf("@%s stopped", botName))
}

func (b *Bot) handleMessage(tgMessage *tgbotapi.Message) {
	chatID := tgMessage.Chat.ID

	userChat[tgMessage.From.ID] = chatID

	if tgMessage.Command() == commandStart {
		b.sendWithAlertOnError(b.getStartMenu(chatID))
		b.sendWithAlertOnError(b.getSubjectsList(chatID))
		b.sendAlert(fmt.Sprintf("%s started conversation with @%s", formatUserString(tgMessage.From), botName))
	} else if tgMessage.Text == commandSelectSubject {
		b.sendWithAlertOnError(b.getSubjectsList(chatID))
	} else {
		b.sendNextTask(chatID, tgMessage.From.ID)
	}
}

func (b *Bot) handleCallbackQuery(tgCallbackQuery *tgbotapi.CallbackQuery) {
	chatID := tgCallbackQuery.Message.Chat.ID

	userChat[tgCallbackQuery.From.ID] = chatID

	if tgCallbackQuery.Message.Text == textSelectSubject {
		if b.selectSubject(tgCallbackQuery) {
			b.sendNextTask(chatID, tgCallbackQuery.From.ID)
		}
	} else if b.updateInlineQuestion(tgCallbackQuery) {
		b.sendNextTask(chatID, tgCallbackQuery.From.ID)
	}
}

func (b *Bot) handlePollAnswer(tgPollAnswer *tgbotapi.PollAnswer) {
	userID := tgPollAnswer.User.ID
	chatID := userChat[userID]

	b.sendNextTask(chatID, userID)
}

func (b *Bot) selectSubject(callbackQuery *tgbotapi.CallbackQuery) bool {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	userSelectedSubject[callbackQuery.From.ID] = callbackQuery.Data

	alreadyAnswered := callbackQuery.Data == labelAnswered
	callbackText := ""
	if alreadyAnswered {
		callbackText = fmt.Sprintf(`Для смены предмета, воспользуйтесь кнопкой "%s"`, commandSelectSubject)
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
			b.sendAlert(err.Error())
		}
	}

	tgCallback := tgbotapi.NewCallback(callbackQuery.ID, callbackText)
	if _, err := b.api.Request(tgCallback); err != nil {
		b.sendAlert(err.Error())
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
	} else if strings.HasPrefix(callbackQuery.Data, entity.ExplanationPrefix) {
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
					entity.ExplanationPrefix+correctOptionText,
				),
			),
		)
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
		tgKeyboardUpdate.ReplyMarkup = &tgKeyboard
		if _, err := b.api.Send(tgKeyboardUpdate); err != nil {
			b.sendAlert(err.Error())
		}
	}

	tgCallback := tgbotapi.NewCallback(callbackQuery.ID, callbackText)
	if _, err := b.api.Request(tgCallback); err != nil {
		b.sendAlert(err.Error())
		return false
	}
	return !alreadyAnswered
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

func (b *Bot) getSubjectsList(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(chatID, textSelectSubject)
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(subjectRussian, subjectRussian)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(subjectHistory, subjectHistory)},
	)
	return &tgMessage
}

func (b *Bot) sendNextTask(chatID int64, userID int) {
	for {
		subject := userSelectedSubject[userID]
		var task tgbotapi.Chattable
		switch subject {
		case subjectRussian:
			task = b.getNextTask(b.database.Russian.Tasks, chatID)
		case subjectHistory:
			task = b.getNextTask(b.database.History.Tasks, chatID)
		default:
			task = b.getSubjectsList(chatID)
		}
		if b.sendWithAlertOnError(task) {
			break
		}
	}
}

func (b *Bot) getNextTask(tasks []*entity.Task, chatID int64) tgbotapi.Chattable {
	task := tasks[rand.Intn(len(tasks))]
	if tgPoll := task.MakeTelegramPoll(chatID); tgPoll != nil {
		return tgPoll
	}
	return task.MakeTelegramMessage(chatID)
}

func (b *Bot) sendWithAlertOnError(tgChattable tgbotapi.Chattable) bool {
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(err.Error())
		return false
	}
	return true
}

func (b *Bot) sendAlert(text string) {
	log.Print(text)
	tgMessage := tgbotapi.NewMessage(alertsChatID, fmt.Sprintf("[%s] %s", b.hostName, text))
	_, err := b.api.Send(tgMessage)
	if err != nil {
		log.Printf("Error on sending alert: %s", err)
	}
}
