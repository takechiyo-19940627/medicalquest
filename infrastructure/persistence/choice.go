package persistence

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent/question"
)

type ChoiceRepository struct {
	db *ent.Client
}

func NewChoiceRepository(db *ent.Client) repository.ChoiceRepository {
	return ChoiceRepository{
		db: db,
	}
}

func (c ChoiceRepository) Save(ctx context.Context, choice entity.Choice) error {
	question, err := c.db.
		Question.
		Query().
		Where(question.UID(choice.QuestionUID.String())).
		First(ctx)

	if err != nil {
		return err
	}

	return c.db.Choice.Create().
		SetUID(choice.UID.String()).
		SetQuestionID(question.ID).
		SetContent(choice.Content).
		SetIsCorrect(choice.IsCorrect).
		Exec(ctx)
}
