package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ravil23/usebot/telegrambot/collection/entity"
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

func init() {
	flag.StringVar(&russianSubjectPath, "russian", "", "Database of Russian subject tasks")
	flag.StringVar(&mathAdvancedSubjectPath, "math-advanced", "", "Database of Advanced Math subject tasks")
	flag.StringVar(&mathBasicSubjectPath, "math-basic", "", "Database of Basic Math subject tasks")
	flag.StringVar(&physicsSubjectPath, "physics", "", "Database of Physics subject tasks")
	flag.StringVar(&chemistrySubjectPath, "chemistry", "", "Database of Chemistry subject tasks")
	flag.StringVar(&itSubjectPath, "it", "", "Database of IT subject tasks")
	flag.StringVar(&biologySubjectPath, "biology", "", "Database of Biology subject tasks")
	flag.StringVar(&historySubjectPath, "history", "", "Database of History subject tasks")
	flag.StringVar(&geographySubjectPath, "geography", "", "Database of Geography subject tasks")
	flag.StringVar(&englishSubjectPath, "english", "", "Database of English subject tasks")
	flag.StringVar(&germanSubjectPath, "german", "", "Database of German subject tasks")
	flag.StringVar(&frenchSubjectPath, "french", "", "Database of French subject tasks")
	flag.StringVar(&socialSubjectPath, "social", "", "Database of Social subject tasks")
	flag.StringVar(&spanishSubjectPath, "spanish", "", "Database of Spanish subject tasks")
	flag.StringVar(&literatureSubjectPath, "literature", "", "Database of Literature subject tasks")
}

func parseArguments() {
	log.Printf("Arguments: [%s]", strings.Join(flag.Args(), ", "))
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
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	parseArguments()

	database := entity.NewDatabase(
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
	bot.HealthCheck()
	bot.Run()
}
