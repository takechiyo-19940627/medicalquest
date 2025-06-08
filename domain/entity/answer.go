package entity

import (
	"errors"

	"github.com/samber/lo"
)

type Answer struct {
	question          Question
	selectedChoiceUID UID
}

var ERR_NO_CHOICE_UID = errors.New("no choice uid")

func NewAnswer(question Question, selectedChoiceUID UID) (Answer, error) {
	if selectedChoiceUID.String() == "" {
		return Answer{}, ERR_NO_CHOICE_UID
	}

	return Answer{
		question:          question,
		selectedChoiceUID: selectedChoiceUID,
	}, nil
}

func (a Answer) IsCorrect() bool {
	correct := lo.Filter(a.question.Choices, func(c Choice, _ int) bool {
		return c.IsCorrect
	})[0]

	return a.selectedChoiceUID == correct.UID
}
