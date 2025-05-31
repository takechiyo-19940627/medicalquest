package response

import "github.com/takechiyo-19940627/medicalquest/service/dto"

type QuestionWithChoicesResponse struct {
	Data QuestionWithChoicesItem `json:"data"`
}

type QuestionWithChoicesItem struct {
	UID           string       `json:"uid"`
	ReferenceCode string       `json:"reference_code"`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	Choices       []ChoiceItem `json:"choices"`
}

type ChoiceItem struct {
	UID       string `json:"uid"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

func NewQuestionWithChoicesResponse(question dto.QuestionWithChoicesResult) QuestionWithChoicesResponse {
	choices := make([]ChoiceItem, len(question.Choices))
	for i, c := range question.Choices {
		choices[i] = ChoiceItem{
			UID:       c.UID,
			Content:   c.Content,
			IsCorrect: c.IsCorrect,
		}
	}

	return QuestionWithChoicesResponse{
		Data: QuestionWithChoicesItem{
			UID:           question.UID,
			ReferenceCode: question.ReferenceCode,
			Title:         question.Title,
			Content:       question.Content,
			Choices:       choices,
		},
	}
}
