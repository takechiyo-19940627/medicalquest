package persistence

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	domainErrors "github.com/takechiyo-19940627/medicalquest/domain/errors"
	"github.com/takechiyo-19940627/medicalquest/domain/repository"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent/question"
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

		m := entity.Question{
			UID:           entity.ToUID(q.UID),
			ReferenceCode: rc,
			Title:         q.Title,
			Content:       q.Content,
			Choices:       []entity.Choice{},
		}
		questions = append(questions, m)
	}

	return questions, nil
}

func (q QuestionRepository) FindByID(ctx context.Context, id entity.UID) (entity.Question, error) {
	qs, err := q.db.Question.
		Query().
		Where(question.UID(id.String())).
		WithChoices().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// NotFoundの場合はドメインエラーを返す
			return entity.Question{}, domainErrors.ErrQuestionNotFound
		}
		// その他のデータベースエラーはラップして返す
		return entity.Question{}, fmt.Errorf("failed to find question: %w", err)
	}

	choices := make([]entity.Choice, len(qs.Edges.Choices))
	for i, c := range qs.Edges.Choices {
		choices[i] = entity.Choice{
			UID:         entity.ToUID(c.UID),
			QuestionUID: entity.ToUID(qs.UID),
			Content:     c.Content,
			IsCorrect:   c.IsCorrect,
		}
	}

	return entity.Question{
		UID:           entity.ToUID(qs.UID),
		ReferenceCode: lo.FromPtr(qs.ReferenceCode),
		Title:         qs.Title,
		Content:       qs.Content,
		Choices:       choices,
	}, nil
}

func (q QuestionRepository) Save(ctx context.Context, uid entity.UID, referenceCode, title, content string) error {
	_, err := q.db.Question.
		Create().
		SetUID(uid.String()).
		SetReferenceCode(referenceCode).
		SetTitle(title).
		SetContent(content).
		Save(ctx)

	return err
}
