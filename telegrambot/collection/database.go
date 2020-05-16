package collection

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
		Subjects: map[string]*Subject{
			SubjectNameRussian:      parseSubjectFileOrPanic(SubjectNameRussian, russianSubjectPath),
			SubjectNameMathAdvanced: parseSubjectFileOrPanic(SubjectNameMathAdvanced, mathAdvancedSubjectPath),
			SubjectNameMathBasic:    parseSubjectFileOrPanic(SubjectNameMathBasic, mathBasicSubjectPath),
			SubjectNamePhysics:      parseSubjectFileOrPanic(SubjectNamePhysics, physicsSubjectPath),
			SubjectNameChemistry:    parseSubjectFileOrPanic(SubjectNameChemistry, chemistrySubjectPath),
			SubjectNameIT:           parseSubjectFileOrPanic(SubjectNameIT, itSubjectPath),
			SubjectNameBiology:      parseSubjectFileOrPanic(SubjectNameBiology, biologySubjectPath),
			SubjectNameHistory:      parseSubjectFileOrPanic(SubjectNameHistory, historySubjectPath),
			SubjectNameGeography:    parseSubjectFileOrPanic(SubjectNameGeography, geographySubjectPath),
			SubjectNameEnglish:      parseSubjectFileOrPanic(SubjectNameEnglish, englishSubjectPath),
			SubjectNameGerman:       parseSubjectFileOrPanic(SubjectNameGerman, germanSubjectPath),
			SubjectNameFrench:       parseSubjectFileOrPanic(SubjectNameFrench, frenchSubjectPath),
			SubjectNameSocial:       parseSubjectFileOrPanic(SubjectNameSocial, socialSubjectPath),
			SubjectNameSpanish:      parseSubjectFileOrPanic(SubjectNameSpanish, spanishSubjectPath),
			SubjectNameLiterature:   parseSubjectFileOrPanic(SubjectNameLiterature, literatureSubjectPath),
		},
	}
}

func (d *Database) Show() {
	for name, subject := range d.Subjects {
		log.Printf("%s: %s", name, subject)
	}
}
