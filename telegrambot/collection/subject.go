package collection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

const (
	SubjectNameRussian      = "Русский язык"
	SubjectNameMathAdvanced = "Математика (профильный)"
	SubjectNameMathBasic    = "Математика (базовый)"
	SubjectNamePhysics      = "Физика"
	SubjectNameChemistry    = "Химия"
	SubjectNameIT           = "Информатика и ИКТ"
	SubjectNameBiology      = "Биология"
	SubjectNameHistory      = "История"
	SubjectNameGeography    = "География"
	SubjectNameEnglish      = "Английский"
	SubjectNameGerman       = "Немецкий"
	SubjectNameFrench       = "Французский"
	SubjectNameSocial       = "Обществознание"
	SubjectNameSpanish      = "Испанский"
	SubjectNameLiterature   = "Литература"
)

var AllSubjectNames = []string{
	SubjectNameRussian,
	SubjectNameMathAdvanced,
	SubjectNameMathBasic,
	SubjectNamePhysics,
	SubjectNameIT,
	SubjectNameChemistry,
	SubjectNameBiology,
	SubjectNameGeography,
	SubjectNameHistory,
	SubjectNameSocial,
	SubjectNameEnglish,
	SubjectNameGerman,
	SubjectNameFrench,
	SubjectNameSpanish,
	SubjectNameLiterature,
}

type Subject struct {
	Tasks []*Task `json:"tasks"`

	Name             string   `json:"-"`
	LowLevelTasks    []*Task  `json:"-"`
	MediumLevelTasks []*Task  `json:"-"`
	HighLevelTasks   []*Task  `json:"-"`
	AllThemes        []string `json:"-"`
}

func (s *Subject) String() string {
	return fmt.Sprintf(
		"Subject{Tasks: %d (%d low, %d medium, %d high), Themes: %d}",
		len(s.Tasks), len(s.LowLevelTasks), len(s.MediumLevelTasks), len(s.HighLevelTasks), len(s.AllThemes),
	)
}

func (s *Subject) extractAllThemes() {
	allThemes := make(map[string]struct{})
	for _, task := range s.Tasks {
		for _, theme := range task.Themes {
			allThemes[theme] = struct{}{}
		}
	}
	s.AllThemes = make([]string, 0, len(allThemes))
	for theme := range allThemes {
		s.AllThemes = append(s.AllThemes, theme)
	}
	sort.Slice(s.AllThemes, func(i, j int) bool {
		return s.AllThemes[i] < s.AllThemes[j]
	})
}

func (s *Subject) groupTasksByLevels() {
	s.LowLevelTasks = make([]*Task, 0, len(s.Tasks))
	s.MediumLevelTasks = make([]*Task, 0, len(s.Tasks))
	s.HighLevelTasks = make([]*Task, 0, len(s.Tasks))
	for _, task := range s.Tasks {
		switch task.Level {
		case LevelLow:
			s.LowLevelTasks = append(s.LowLevelTasks, task)
		case LevelMedium:
			s.MediumLevelTasks = append(s.MediumLevelTasks, task)
		case LevelHigh:
			s.HighLevelTasks = append(s.HighLevelTasks, task)
		}
	}
}

func parseSubjectFile(name, path string) (*Subject, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var subject Subject
	err = json.Unmarshal(jsonData, &subject)
	if err != nil {
		return nil, err
	}
	subject.extractAllThemes()
	subject.groupTasksByLevels()
	for _, task := range subject.Tasks {
		task.SubjectName = name
	}
	subject.Name = name
	return &subject, nil
}

func parseSubjectFileOrPanic(name, path string) *Subject {
	subject, err := parseSubjectFile(name, path)
	if err != nil {
		log.Panic(err)
	}
	return subject
}
