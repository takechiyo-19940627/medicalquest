package persistence

import (
	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
)

type QuestionRepository struct{}

func NewQuestionRepository() repository.QuestionRepository {
	return QuestionRepository{}
}

func (q QuestionRepository) FindAll() []entity.Question {
	return []entity.Question{}
}

func (q QuestionRepository) FindByID(id string) entity.Question {
	return entity.Question{}
}
