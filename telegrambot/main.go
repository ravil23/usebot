package main

import (
	"flag"
	"log"
	"os"

	"github.com/ravil23/usebot/telegrambot/collection/entity"
	"github.com/ravil23/usebot/telegrambot/telegram"
)

var russianSubjectPath string

func init() {
	flag.StringVar(&russianSubjectPath, "rus", "", "Database of russian subject tasks")
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

	database := entity.NewDatabase(russianSubjectPath)
	log.Printf("Found Russian tasks count: %d", len(database.Russian.Tasks))

	bot := telegram.NewBot(database)
	bot.Init()
	bot.HealthCheck()
	bot.Run()
}
