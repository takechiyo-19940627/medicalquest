package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/takechiyo-19940627/medicalquest/domain/entity"
)

type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Question), args.Error(1)
}

func (m *MockQuestionRepository) FindByID(ctx context.Context, id entity.UID) (entity.Question, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Question), args.Error(1)
}

func (m *MockQuestionRepository) Save(ctx context.Context, id entity.UID, referenceCode, title, content string) error {
	args := m.Called(ctx, id, referenceCode, title, content)
	return args.Error(0)
}