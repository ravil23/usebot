package entity

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
