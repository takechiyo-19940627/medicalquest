package service

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
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
		return nil, err
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
		return dto.QuestionWithChoicesResult{}, err
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
