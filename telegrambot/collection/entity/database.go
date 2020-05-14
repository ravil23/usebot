package entity

import (
	"log"
)

type Database struct {
	Russian *Subject
	History *Subject
}

func NewDatabase(russianSubjectPath, historySubjectPath string) *Database {
	return &Database{
		Russian: parseSubjectFileOrPanic(russianSubjectPath),
		History: parseSubjectFileOrPanic(historySubjectPath),
	}
}

func (d *Database) Show() {
	log.Print("Russian subject")
	d.Russian.show()

	log.Print("History subject")
	d.History.show()
}
