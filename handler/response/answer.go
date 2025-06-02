package response

import "github.com/takechiyo-19940627/medicalquest/service/dto"

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
