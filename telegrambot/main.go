package main

import (
	"flag"
	"log"
	"os"

	"github.com/ravil23/usebot/telegrambot/collection"
	"github.com/ravil23/usebot/telegrambot/telegram"
)

var russianSubjectPath string
var mathAdvancedSubjectPath string
var mathBasicSubjectPath string
var physicsSubjectPath string
var chemistrySubjectPath string
var itSubjectPath string
var biologySubjectPath string
var historySubjectPath string
var geographySubjectPath string
var englishSubjectPath string
var germanSubjectPath string
var frenchSubjectPath string
var socialSubjectPath string
var spanishSubjectPath string
var literatureSubjectPath string
var test bool

func init() {
	russianSubjectPath = os.Getenv("SUBJECT_RUSSIAN")
	mathAdvancedSubjectPath = os.Getenv("SUBJECT_MATH_ADVANCED")
	mathBasicSubjectPath = os.Getenv("SUBJECT_MATH_BASIC")
	physicsSubjectPath = os.Getenv("SUBJECT_PHYSICS")
	chemistrySubjectPath = os.Getenv("SUBJECT_CHEMISTRY")
	itSubjectPath = os.Getenv("SUBJECT_IT")
	biologySubjectPath = os.Getenv("SUBJECT_BIOLOGY")
	historySubjectPath = os.Getenv("SUBJECT_HISTORY")
	geographySubjectPath = os.Getenv("SUBJECT_GEOGRAPHY")
	englishSubjectPath = os.Getenv("SUBJECT_ENGLISH")
	germanSubjectPath = os.Getenv("SUBJECT_GERMAN")
	frenchSubjectPath = os.Getenv("SUBJECT_FRENCH")
	socialSubjectPath = os.Getenv("SUBJECT_SOCIAL")
	spanishSubjectPath = os.Getenv("SUBJECT_SPANISH")
	literatureSubjectPath = os.Getenv("SUBJECT_LITERATURE")
	flag.BoolVar(&test, "test", false, "Send all tasks to alerts channel for testing")
}

func parseArguments() {
	flag.Parse()
	if russianSubjectPath == "" ||
		mathAdvancedSubjectPath == "" ||
		mathBasicSubjectPath == "" ||
		physicsSubjectPath == "" ||
		chemistrySubjectPath == "" ||
		itSubjectPath == "" ||
		biologySubjectPath == "" ||
		historySubjectPath == "" ||
		geographySubjectPath == "" ||
		englishSubjectPath == "" ||
		germanSubjectPath == "" ||
		frenchSubjectPath == "" ||
		socialSubjectPath == "" ||
		spanishSubjectPath == "" ||
		literatureSubjectPath == "" {
		log.Printf("Some of subject paths are empty")
		os.Exit(2)
	}
}

func main() {
	parseArguments()

	database := collection.NewDatabase(
		russianSubjectPath,
		mathAdvancedSubjectPath,
		mathBasicSubjectPath,
		physicsSubjectPath,
		chemistrySubjectPath,
		itSubjectPath,
		biologySubjectPath,
		historySubjectPath,
		geographySubjectPath,
		englishSubjectPath,
		germanSubjectPath,
		frenchSubjectPath,
		socialSubjectPath,
		spanishSubjectPath,
		literatureSubjectPath,
	)
	database.Show()

	bot := telegram.NewBot(database)
	bot.Init()
	if test {
		bot.TestAllTasks(telegram.AlertsChatID)
		return
	}
	bot.HealthCheck()
	bot.Run()
}
