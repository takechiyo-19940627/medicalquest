package service

import (
	"context"
	"errors"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
)

type AnswerService struct {
	questionRepository repository.QuestionRepository
}

var ERR_INVALID_CHOICE_ID = errors.New("invalid choice id")

func NewAnswerService(questionRepository repository.QuestionRepository) *AnswerService {
	return &AnswerService{
		questionRepository,
	}
}

func (a *AnswerService) Submit(ctx context.Context, request dto.AnswerRequest) (dto.AnswerResult, error) {
	q, err := a.questionRepository.FindByID(ctx, entity.ToUID(request.QuestionID))
	if err != nil {
		return dto.AnswerResult{}, err
	}

	selectedUID := entity.ToUID(request.QuestionID)
	if hasChoice := q.HasChoice(selectedUID); !hasChoice {
		return dto.AnswerResult{}, ERR_INVALID_CHOICE_ID
	}

	answer, err := entity.NewAnswer(q, selectedUID)
	if err != nil {
		return dto.AnswerResult{}, err
	}

	return dto.AnswerResult{
		IsCorrect: answer.IsCorrect(),
	}, nil
}
