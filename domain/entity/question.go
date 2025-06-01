package entity

import (
	"errors"

	"github.com/samber/lo"
)

const MAX_CHOICE_COUNT = 4

var ERR_MAX_CHOICE_COUNT = errors.New("choice count is max")
var ERR_NO_CHOICE = errors.New("no choice")
var ERR_INVALID_CORRECT_CHOICE_COUNT = errors.New("invalid correct choice count")

type Question struct {
	UID           UID
	ReferenceCode string
	Title         string
	Content       string
	Choices       []Choice
}

func NewQuestion(referenceCode, title, content string, choices []Choice) (Question, error) {
	if err := validate(choices); err != nil {
		return Question{}, err
	}

	return Question{
		UID:           GenerateUID(),
		ReferenceCode: referenceCode,
		Title:         title,
		Content:       content,
		Choices:       choices,
	}, nil
}

func validate(choices []Choice) error {
	if len(choices) == 0 {
		return ERR_NO_CHOICE
	}

	if len(choices) > MAX_CHOICE_COUNT {
		return ERR_MAX_CHOICE_COUNT
	}

	corrects := lo.Filter(choices, func(c Choice, i int) bool {
		return c.IsCorrect
	})
	if len(corrects) == 0 || len(corrects) > 1 {
		return ERR_INVALID_CORRECT_CHOICE_COUNT
	}

	return nil
}
