package entity

import "errors"

const MAX_CHOICE_COUNT = 4

var ERR_MAX_CHOICE_COUNT = errors.New("choice count is max")

type Question struct {
	UID           UID
	ReferenceCode string
	Title         string
	Content       string
	Choices       []Choice
}

func NewQuestion(referenceCode, title, content string) Question {
	return Question{
		UID:           GenerateUID(),
		ReferenceCode: referenceCode,
		Title:         title,
		Content:       content,
	}
}

func (q *Question) AddChoice(choice Choice) error {
	if !q.canAddChoice() {
		return ERR_MAX_CHOICE_COUNT
	}

	q.Choices = append(q.Choices, choice)
	return nil
}

func (q Question) canAddChoice() bool {
	return len(q.Choices) < MAX_CHOICE_COUNT
}
