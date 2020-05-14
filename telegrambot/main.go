package main

import (
	"flag"
	"os"

	"github.com/ravil23/usebot/telegrambot/collection/entity"
	"github.com/ravil23/usebot/telegrambot/telegram"
)

var russianSubjectPath string
var historySubjectPath string

func init() {
	flag.StringVar(&russianSubjectPath, "rus", "", "Database of Russian subject tasks")
	flag.StringVar(&historySubjectPath, "history", "", "Database of History subject tasks")
}

func parseArguments() {
	flag.Parse()
	if russianSubjectPath == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	parseArguments()

	database := entity.NewDatabase(russianSubjectPath, historySubjectPath)
	database.Show()

	bot := telegram.NewBot(database)
	bot.Init()
	bot.HealthCheck()
	bot.Run()
}
