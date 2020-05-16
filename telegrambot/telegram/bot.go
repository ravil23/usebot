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

	"github.com/ravil23/usebot/telegrambot/collection"
)

const (
	initializationRetryPeriod     = time.Second
	initializationMaxRetriesCount = 30
	timeoutSeconds                = 60
	listenersPoolSize             = 10
	delayBetweenQuestions         = time.Second
)

type Bot struct {
	hostName string
	api      *tgbotapi.BotAPI
	database *collection.Database
}

func NewBot(database *collection.Database) *Bot {
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

func (b *Bot) TestAllTasks(chatID int64) {
	log.Printf("Bot is sending all tasks to chat %d", chatID)

	for name, subject := range b.database.Subjects {
		log.Printf("%s: %s", name, subject)
		for _, task := range subject.Tasks {
			var tgChattable tgbotapi.Chattable
			if task.SendAsPoll {
				tgChattable = task.MakeTelegramPoll(chatID)
			} else {
				tgChattable = task.MakeTelegramMessage(chatID)
			}
			if _, err := b.api.Send(tgChattable); err != nil {
				b.sendAlert(fmt.Sprintf("Error on sending %v: %s", tgChattable, err))
				panic(err)
			}
		}
	}
}

func (b *Bot) serve() {
	b.sendAlert(fmt.Sprintf("@%s started", Bot11Name))
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	b.sendAlert(fmt.Sprintf("@%s stopped", Bot11Name))
}

func (b *Bot) handleMessage(tgMessage *tgbotapi.Message) {
	chatID := tgMessage.Chat.ID

	userChat[tgMessage.From.ID] = chatID

	if tgMessage.Command() == commandStart {
		b.sendWithAlertOnError(b.getStartMenu(chatID))
		b.sendWithAlertOnError(b.getSubjectsList(chatID))
		b.sendAlert(fmt.Sprintf("%s started conversation with @%s", formatUserString(tgMessage.From), Bot11Name))
	} else if tgMessage.Text == commandSelectSubject {
		b.sendWithAlertOnError(b.getSubjectsList(chatID))
	} else if tgMessage.Text == commandSelectLevel {
		b.sendWithAlertOnError(b.getLevelsList(chatID))
	} else if tgMessage.Text == commandNext {
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
	} else if tgCallbackQuery.Message.Text == textSelectLevel {
		if b.selectLevel(tgCallbackQuery) {
			b.sendNextTask(chatID, tgCallbackQuery.From.ID)
		}
	} else {
		b.updateInlineQuestion(tgCallbackQuery)
	}
}

func (b *Bot) handlePollAnswer(_ *tgbotapi.PollAnswer) {
}

func (b *Bot) selectSubject(callbackQuery *tgbotapi.CallbackQuery) bool {
	userSelectedSubject[callbackQuery.From.ID] = callbackQuery.Data

	popupIfSucceeded := fmt.Sprintf(`Выбран предмет "%s"`, callbackQuery.Data)
	popupIfAlreadyAnswered := fmt.Sprintf(`Для смены предмета, воспользуйтесь кнопкой "%s"`, commandSelectSubject)

	return b.updateMessageAfterSelect(callbackQuery, popupIfSucceeded, popupIfAlreadyAnswered, "📖️")
}

func (b *Bot) selectLevel(callbackQuery *tgbotapi.CallbackQuery) bool {
	userSelectedLevel[callbackQuery.From.ID] = callbackQuery.Data

	popupIfSucceeded := fmt.Sprintf(`Выбрана сложность "%s"`, callbackQuery.Data)
	popupIfAlreadyAnswered := fmt.Sprintf(`Для смены сложности, воспользуйтесь кнопкой "%s"`, commandSelectLevel)

	return b.updateMessageAfterSelect(callbackQuery, popupIfSucceeded, popupIfAlreadyAnswered, "🎓️")
}

func (b *Bot) updateMessageAfterSelect(callbackQuery *tgbotapi.CallbackQuery, popupIfSucceeded, popupIfAlreadyAnswered, marker string) bool {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	alreadyAnswered := callbackQuery.Data == labelAnswered
	var popupText string
	if alreadyAnswered {
		popupText = popupIfAlreadyAnswered
	} else {
		popupText = popupIfSucceeded
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
					tgButton.Text += " " + marker
				}
				tgButtons = append(tgButtons, tgButton)
			}
			tgRows = append(tgRows, tgbotapi.NewInlineKeyboardRow(tgButtons...))
		}
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
		tgKeyboardUpdate.ReplyMarkup = &tgKeyboard
		b.sendWithAlertOnError(tgKeyboardUpdate)
	}

	b.sendCallback(callbackQuery.ID, popupText)
	return !alreadyAnswered
}

func (b *Bot) updateInlineQuestion(callbackQuery *tgbotapi.CallbackQuery) bool {
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	alreadyAnswered := false
	var popupText string
	if callbackQuery.Data == labelAnswered {
		alreadyAnswered = true
		popupText = "К сожалению изменить ответ нельзя"
	} else if strings.HasPrefix(callbackQuery.Data, collection.ExplanationPrefix) {
		alreadyAnswered = true
		tgMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("<b>%s</b>\n%s", textExplanation, callbackQuery.Data))
		tgMessage.ReplyToMessageID = messageID
		tgMessage.ParseMode = tgbotapi.ModeHTML
		b.sendWithAlertOnError(tgMessage)
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
			popupText = collection.ExplanationPrefix + correctOptionText
		}
		tgRows = append(
			tgRows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					textExplanation,
					collection.ExplanationPrefix+correctOptionText,
				),
			),
		)
		tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRows...)
		tgKeyboardUpdate.ReplyMarkup = &tgKeyboard
		b.sendWithAlertOnError(tgKeyboardUpdate)
	}

	b.sendCallback(callbackQuery.ID, popupText)
	return !alreadyAnswered
}

func (b *Bot) getStartMenu(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(
		chatID,
		fmt.Sprintf(`Выбери предмет, чтобы начать подготовку. Его всегда можно сменить, нажав на кнопку "%s".`, commandSelectSubject),
	)
	tgButtons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(commandSelectSubject),
		tgbotapi.NewKeyboardButton(commandSelectLevel),
		tgbotapi.NewKeyboardButton(commandNext),
	}
	tgMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgButtons)
	return &tgMessage
}

func (b *Bot) getSubjectsList(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(chatID, textSelectSubject)
	tgRows := make([][]tgbotapi.InlineKeyboardButton, 0, len(collection.AllSubjectNames))
	tgButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(collection.AllSubjectNames))
	for _, subjectName := range collection.AllSubjectNames {
		if subject, found := b.database.Subjects[subjectName]; found && len(subject.Tasks) == 0 {
			continue
		}
		tgButtons = append(tgButtons, tgbotapi.NewInlineKeyboardButtonData(subjectName, subjectName))
	}
	tgRow := make([]tgbotapi.InlineKeyboardButton, 0)
	for i := range tgButtons {
		tgRow = append(tgRow, tgButtons[i])
		if (i+1) == len(tgButtons) || (i+1)%2 == 0 {
			tgRows = append(tgRows, tgRow)
			tgRow = make([]tgbotapi.InlineKeyboardButton, 0)
		}
	}
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgRows...)
	return &tgMessage
}

func (b *Bot) getLevelsList(chatID int64) tgbotapi.Chattable {
	tgMessage := tgbotapi.NewMessage(chatID, textSelectLevel)
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(collection.LevelLow.String(), collection.LevelLow.String())},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(collection.LevelMedium.String(), collection.LevelMedium.String())},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(collection.LevelHigh.String(), collection.LevelHigh.String())},
	)
	return &tgMessage
}

func (b *Bot) sendCallback(callbackID, callbackText string) bool {
	tgCallback := tgbotapi.NewCallback(callbackID, callbackText)
	if _, err := b.api.Request(tgCallback); err != nil {
		b.sendAlert(err.Error())
		return false
	}
	return true
}

func (b *Bot) sendNextTaskWithDelay(chatID int64, userID int) {
	go func() {
		time.Sleep(delayBetweenQuestions)
		b.sendNextTask(chatID, userID)
	}()
}

func (b *Bot) sendNextTask(chatID int64, userID int) {
	for {
		subjectName := userSelectedSubject[userID]
		level := userSelectedLevel[userID]
		var task tgbotapi.Chattable
		if subject, found := b.database.Subjects[subjectName]; found {
			task = b.getNextTaskByLevel(subject, chatID, level)
		} else {
			task = b.getSubjectsList(chatID)
		}
		if b.sendWithAlertOnError(task) {
			break
		}
	}
}

func (b *Bot) getNextTaskByLevel(subject *collection.Subject, chatID int64, level string) tgbotapi.Chattable {
	switch level {
	case collection.LevelHigh.String():
		if len(subject.HighLevelTasks) > 0 {
			return b.getNextTask(subject.HighLevelTasks, chatID)
		}
		fallthrough
	case collection.LevelMedium.String():
		if len(subject.MediumLevelTasks) > 0 {
			return b.getNextTask(subject.MediumLevelTasks, chatID)
		}
		fallthrough
	case collection.LevelLow.String():
		if len(subject.LowLevelTasks) > 0 {
			return b.getNextTask(subject.LowLevelTasks, chatID)
		}
		fallthrough
	default:
		return b.getNextTask(subject.Tasks, chatID)
	}
}

func (b *Bot) getNextTask(tasks []*collection.Task, chatID int64) tgbotapi.Chattable {
	task := tasks[rand.Intn(len(tasks))]
	if task.SendAsPoll {
		return task.MakeTelegramPoll(chatID)
	} else {
		return task.MakeTelegramMessage(chatID)
	}
}

func (b *Bot) sendWithAlertOnError(tgChattable tgbotapi.Chattable) bool {
	if _, err := b.api.Send(tgChattable); err != nil {
		b.sendAlert(fmt.Sprintf("Error on sending %v: %s", tgChattable, err))
		return false
	}
	return true
}

func (b *Bot) sendAlert(text string) {
	log.Print(text)
	tgMessage := tgbotapi.NewMessage(AlertsChatID, fmt.Sprintf("[%s] %s", b.hostName, text))
	_, err := b.api.Send(tgMessage)
	if err != nil {
		log.Printf("Error on sending alert: %s", err)
	}
}
