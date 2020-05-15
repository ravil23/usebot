package entity

import (
	"log"
)

type Database struct {
	Subjects map[string]*Subject
}

func NewDatabase(
	russianSubjectPath string,
	mathAdvancedSubjectPath string,
	mathBasicSubjectPath string,
	physicsSubjectPath string,
	chemistrySubjectPath string,
	itSubjectPath string,
	biologySubjectPath string,
	historySubjectPath string,
	geographySubjectPath string,
	englishSubjectPath string,
	germanSubjectPath string,
	frenchSubjectPath string,
	socialSubjectPath string,
	spanishSubjectPath string,
	literatureSubjectPath string,
) *Database {

	return &Database{
		Subjects: map[string]*Subject{SubjectNameRussian: parseSubjectFileOrPanic(russianSubjectPath),
			SubjectNameMathAdvanced: parseSubjectFileOrPanic(mathAdvancedSubjectPath),
			SubjectNameMathBasic:    parseSubjectFileOrPanic(mathBasicSubjectPath),
			SubjectNamePhysics:      parseSubjectFileOrPanic(physicsSubjectPath),
			SubjectNameChemistry:    parseSubjectFileOrPanic(chemistrySubjectPath),
			SubjectNameIT:           parseSubjectFileOrPanic(itSubjectPath),
			SubjectNameBiology:      parseSubjectFileOrPanic(biologySubjectPath),
			SubjectNameHistory:      parseSubjectFileOrPanic(historySubjectPath),
			SubjectNameGeography:    parseSubjectFileOrPanic(geographySubjectPath),
			SubjectNameEnglish:      parseSubjectFileOrPanic(englishSubjectPath),
			SubjectNameGerman:       parseSubjectFileOrPanic(germanSubjectPath),
			SubjectNameFrench:       parseSubjectFileOrPanic(frenchSubjectPath),
			SubjectNameSocial:       parseSubjectFileOrPanic(socialSubjectPath),
			SubjectNameSpanish:      parseSubjectFileOrPanic(spanishSubjectPath),
			SubjectNameLiterature:   parseSubjectFileOrPanic(literatureSubjectPath),
		},
	}
}

func (d *Database) Show() {
	for name, subject := range d.Subjects {
		log.Printf("Subject: %s", name)
		subject.show()
	}
}
