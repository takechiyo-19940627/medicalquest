package service

import (
	"context"
	"errors"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	domainErrors "github.com/takechiyo-19940627/medicalquest/domain/errors"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
	serviceErrors "github.com/takechiyo-19940627/medicalquest/service/errors"
)

type QuestionService struct {
	questionRepository repository.QuestionRepository
}

func NewQuestionService(questionRepository repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepository,
	}
}

func (s *QuestionService) FindAll(ctx context.Context) ([]dto.QuestionResult, error) {
	qs, err := s.questionRepository.FindAll(ctx)
	if err != nil {
		return nil, serviceErrors.NewInternalError("問題一覧の取得に失敗しました", err)
	}

	results := make([]dto.QuestionResult, len(qs))
	for i, q := range qs {
		results[i] = dto.QuestionResult{
			UID:           q.UID.String(),
			ReferenceCode: q.ReferenceCode,
			Title:         q.Title,
			Content:       q.Content,
		}
	}

	return results, nil
}

func (s *QuestionService) FindByID(ctx context.Context, id string) (dto.QuestionWithChoicesResult, error) {
	q, err := s.questionRepository.FindByID(ctx, entity.ToUID(id))
	if err != nil {
		if errors.Is(err, domainErrors.ErrQuestionNotFound) {
			return dto.QuestionWithChoicesResult{}, serviceErrors.NewNotFoundError("指定された問題が見つかりません", err)
		}
		return dto.QuestionWithChoicesResult{}, serviceErrors.NewInternalError("問題の取得に失敗しました", err)
	}

	results := dto.QuestionWithChoicesResult{
		QuestionResult: dto.QuestionResult{
			UID:           q.UID.String(),
			ReferenceCode: q.ReferenceCode,
			Title:         q.Title,
			Content:       q.Content,
		},
	}

	return results, nil
}

func (s *QuestionService) Create() error {
	return nil
}

func (s *QuestionService) Submit(ctx context.Context, request dto.AnswerRequest) (dto.AnswerResult, error) {
	q, err := s.questionRepository.FindByID(ctx, entity.ToUID(request.QuestionID))
	if err != nil {
		if errors.Is(err, domainErrors.ErrQuestionNotFound) {
			return dto.AnswerResult{}, serviceErrors.NewNotFoundError("指定された問題が見つかりません", err)
		}
		return dto.AnswerResult{}, serviceErrors.NewInternalError("問題の取得に失敗しました", err)
	}

	selectedUID := entity.ToUID(request.SelectedChoiceID)
	if hasChoice := q.HasChoice(selectedUID); !hasChoice {
		return dto.AnswerResult{}, serviceErrors.NewValidationError("無効な選択肢IDです", "selectedChoiceId", nil)
	}

	answer, err := entity.NewAnswer(q, selectedUID)
	if err != nil {
		return dto.AnswerResult{}, serviceErrors.NewInternalError("回答の作成に失敗しました", err)
	}

	return dto.AnswerResult{
		IsCorrect: answer.IsCorrect(),
	}, nil
}
