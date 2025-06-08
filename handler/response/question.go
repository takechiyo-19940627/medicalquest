package response

import (
	"github.com/takechiyo-19940627/medicalquest/service/dto"
)

type QuestionsResponse struct {
	Data []QuestionItem `json:"data"`
}

type QuestionItem struct {
	UID           string `json:"uid"`
	ReferenceCode string `json:"reference_code"`
	Title         string `json:"title"`
	Content       string `json:"content"`
}

func NewQuestionResponse(questions []dto.QuestionResult) QuestionsResponse {
	if len(questions) == 0 {
		return QuestionsResponse{}
	}

	items := make([]QuestionItem, len(questions))

	for _, q := range questions {
		item := QuestionItem{
			UID:           q.UID,
			ReferenceCode: q.ReferenceCode,
			Title:         q.Title,
			Content:       q.Content,
		}
		items = append(items, item)
	}

	return QuestionsResponse{
		Data: items,
	}
}

type AnswerResponse struct {
	Data AnswerItem `json:"data"`
}

type AnswerItem struct {
	IsCorrect bool `json:"is_correct"`
}

func NewAnswerResponse(result dto.AnswerResult) AnswerResponse {
	return AnswerResponse{
		Data: AnswerItem{
			IsCorrect: result.IsCorrect,
		},
	}
}
