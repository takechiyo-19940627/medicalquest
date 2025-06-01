package repository

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
)

type QuestionRepository interface {
	FindAll(ctx context.Context) ([]entity.Question, error)
	FindByID(ctx context.Context, id entity.UID) (entity.Question, error)
	Save(ctx context.Context, id entity.UID, referenceCode, title, content string) error
}
