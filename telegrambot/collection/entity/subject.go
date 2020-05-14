package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

type Subject struct {
	Tasks []*Task `json:"tasks"`

	LowLevelTasks    []*Task  `json:"-"`
	MediumLevelTasks []*Task  `json:"-"`
	HighLevelTasks   []*Task  `json:"-"`
	AllThemes        []string `json:"-"`
}

func (s *Subject) show() {
	log.Printf("Tasks count: %d", len(s.Tasks))
	log.Printf("Low level tasks count: %d", len(s.LowLevelTasks))
	log.Printf("Medium level tasks count: %d", len(s.MediumLevelTasks))
	log.Printf("High level tasks count: %d", len(s.HighLevelTasks))
	log.Printf("Themes count: %d", len(s.AllThemes))
	for _, theme := range s.AllThemes {
		log.Printf("- %s", theme)
	}
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

func parseSubjectFile(path string) (*Subject, error) {
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
	return &subject, nil
}

func parseSubjectFileOrPanic(path string) *Subject {
	subject, err := parseSubjectFile(path)
	if err != nil {
		log.Panic(err)
	}
	return subject
}
