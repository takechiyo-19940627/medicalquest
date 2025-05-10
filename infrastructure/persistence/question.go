package persistence

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent"
)

type QuestionRepository struct {
	db *ent.Client
}

func NewQuestionRepository(db *ent.Client) repository.QuestionRepository {
	return QuestionRepository{
		db,
	}
}

func (q QuestionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
	qs, err := q.db.Question.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var questions []entity.Question
	for _, q := range qs {
		var rc string
		if q.ReferenceCode != nil {
			rc = *q.ReferenceCode
		}

		m := entity.NewQuestionFromPersistence(
			q.UID,
			rc,
			q.Title,
			q.Content,
		)
		questions = append(questions, m)
	}

	return questions, nil
}

func (q QuestionRepository) FindByID(ctx context.Context, id string) (entity.Question, error) {
	return entity.Question{}, nil
}
