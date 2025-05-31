package repository

import (
	"context"

	"github.com/takechiyo-19940627/medicalquest/domain/entity"
)

type ChoiceRepository interface {
	Save(ctx context.Context, choice entity.Choice) error
}
