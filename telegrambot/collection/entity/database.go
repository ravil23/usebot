package entity

type Database struct {
	Russian *Subject
}

func NewDatabase(russianSubjectPath string) *Database {
	return &Database{
		Russian: parseSubjectFileOrPanic(russianSubjectPath),
	}
}
