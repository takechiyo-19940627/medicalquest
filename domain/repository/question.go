package repository

import "github.com/takechiyo-19940627/medicalquest/domain/entity"

type QuestionRepository interface {
	FindAll() []entity.Question
	FindByID(id string) entity.Question
}
