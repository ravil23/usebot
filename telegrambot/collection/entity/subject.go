package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

type Subject struct {
	Tasks []*Task `json:"tasks"`

	AllThemes []string `json:"-"`
}

func (s *Subject) show() {
	log.Printf("Tasks count: %d", len(s.Tasks))
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
	return &subject, nil
}

func parseSubjectFileOrPanic(path string) *Subject {
	subject, err := parseSubjectFile(path)
	if err != nil {
		log.Panic(err)
	}
	return subject
}
