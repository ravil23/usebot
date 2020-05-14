package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Subject struct {
	Tasks []*Task `json:"tasks"`
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
	return &subject, nil
}

func parseSubjectFileOrPanic(path string) *Subject {
	subject, err := parseSubjectFile(path)
	if err != nil {
		log.Panic(err)
	}
	return subject
}
