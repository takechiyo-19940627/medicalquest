package service

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	result "github.com/takechiyo-19940627/medicalquest/service/dto"
)

type QuestionService struct {
	questionRepository repository.QuestionRepository
}

func NewQuestionService(questionRepository repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepository,
	}
}

func (s *QuestionService) FindAll(ctx context.Context) ([]result.QuestionResult, error) {
	qs, err := s.questionRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]result.QuestionResult, len(qs))
	for i, q := range qs {
		results[i] = result.QuestionResult{
			UID:           q.UID.String(),
			ReferenceCode: q.ReferenceCode,
			Title:         q.Title,
			Content:       q.Content,
		}
	}

	return results, nil
}

func (s *QuestionService) FindByID(ctx context.Context, id string) (entity.Question, error) {
	return s.questionRepository.FindByID(ctx, entity.ToUID(id))
}

func (s *QuestionService) Create() error {
	return nil
}
