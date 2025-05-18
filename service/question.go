package service

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
)

type QuestionService struct {
	questionRepository repository.QuestionRepository
}

func NewQuestionService(questionRepository repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepository,
	}
}

func (s *QuestionService) FindAll(ctx context.Context) ([]entity.Question, error) {
	return s.questionRepository.FindAll(ctx)
}

func (s *QuestionService) Create() error {
	return nil
}
