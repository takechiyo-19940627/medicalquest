package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	domainErrors "github.com/takechiyo-19940627/medicalquest/domain/errors"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
	serviceErrors "github.com/takechiyo-19940627/medicalquest/service/errors"
	"github.com/takechiyo-19940627/medicalquest/service/mocks"
)

func TestNewQuestionService(t *testing.T) {
	mockRepo := new(mocks.MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.questionRepository)
}

func TestQuestionService_FindAll(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockQuestionRepository)
		expectedResult []dto.QuestionResult
		expectedError  error
		errorContains  string
	}{
		{
			name: "正常系: 問題一覧を正常に取得できる",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				questions := []entity.Question{
					{
						UID:           entity.ToUID("uid1"),
						ReferenceCode: "REF001",
						Title:         "問題1",
						Content:       "問題1の内容",
					},
					{
						UID:           entity.ToUID("uid2"),
						ReferenceCode: "REF002",
						Title:         "問題2",
						Content:       "問題2の内容",
					},
				}
				m.On("FindAll", mock.Anything).Return(questions, nil)
			},
			expectedResult: []dto.QuestionResult{
				{
					UID:           "uid1",
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
				},
				{
					UID:           "uid2",
					ReferenceCode: "REF002",
					Title:         "問題2",
					Content:       "問題2の内容",
				},
			},
			expectedError: nil,
		},
		{
			name: "正常系: 問題が0件の場合",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindAll", mock.Anything).Return([]entity.Question{}, nil)
			},
			expectedResult: []dto.QuestionResult{},
			expectedError:  nil,
		},
		{
			name: "異常系: リポジトリでエラーが発生",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindAll", mock.Anything).Return([]entity.Question{}, errors.New("DB connection error"))
			},
			expectedResult: nil,
			errorContains:  "問題一覧の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockQuestionRepository)
			tt.mockSetup(mockRepo)
			
			service := NewQuestionService(mockRepo)
			result, err := service.FindAll(context.Background())
			
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else if tt.errorContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				serviceErr, ok := err.(*serviceErrors.ServiceError)
				assert.True(t, ok)
				assert.Equal(t, serviceErrors.TypeInternal, serviceErr.Type)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestQuestionService_FindByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*mocks.MockQuestionRepository)
		expectedResult dto.QuestionWithChoicesResult
		expectedError  error
		errorContains  string
		errorType      serviceErrors.ErrorType
	}{
		{
			name: "正常系: IDで問題を正常に取得できる",
			id:   "uid1",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				question := entity.Question{
					UID:           entity.ToUID("uid1"),
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
					Choices: []entity.Choice{
						{
							UID:         entity.ToUID("choice1"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "選択肢1",
							IsCorrect:   true,
						},
						{
							UID:         entity.ToUID("choice2"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "選択肢2",
							IsCorrect:   false,
						},
					},
				}
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(question, nil)
			},
			expectedResult: dto.QuestionWithChoicesResult{
				QuestionResult: dto.QuestionResult{
					UID:           "uid1",
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
				},
			},
			expectedError: nil,
		},
		{
			name: "異常系: 問題が見つからない",
			id:   "not_found",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindByID", mock.Anything, entity.ToUID("not_found")).Return(entity.Question{}, domainErrors.ErrQuestionNotFound)
			},
			errorContains: "指定された問題が見つかりません",
			errorType:     serviceErrors.TypeNotFound,
		},
		{
			name: "異常系: リポジトリでその他のエラーが発生",
			id:   "uid1",
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(entity.Question{}, errors.New("DB error"))
			},
			errorContains: "問題の取得に失敗しました",
			errorType:     serviceErrors.TypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockQuestionRepository)
			tt.mockSetup(mockRepo)
			
			service := NewQuestionService(mockRepo)
			result, err := service.FindByID(context.Background(), tt.id)
			
			if tt.errorContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				serviceErr, ok := err.(*serviceErrors.ServiceError)
				assert.True(t, ok)
				assert.Equal(t, tt.errorType, serviceErr.Type)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestQuestionService_Create(t *testing.T) {
	mockRepo := new(mocks.MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	
	err := service.Create()
	assert.NoError(t, err)
}

func TestQuestionService_Submit(t *testing.T) {
	tests := []struct {
		name           string
		request        dto.AnswerRequest
		mockSetup      func(*mocks.MockQuestionRepository)
		expectedResult dto.AnswerResult
		expectedError  error
		errorContains  string
		errorType      serviceErrors.ErrorType
	}{
		{
			name: "正常系: 正解の選択肢を選んだ場合",
			request: dto.AnswerRequest{
				QuestionID:       "uid1",
				SelectedChoiceID: "choice1",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				question := entity.Question{
					UID:           entity.ToUID("uid1"),
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
					Choices: []entity.Choice{
						{
							UID:         entity.ToUID("choice1"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "正解",
							IsCorrect:   true,
						},
						{
							UID:         entity.ToUID("choice2"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "不正解",
							IsCorrect:   false,
						},
					},
				}
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(question, nil)
			},
			expectedResult: dto.AnswerResult{
				IsCorrect: true,
			},
		},
		{
			name: "正常系: 不正解の選択肢を選んだ場合",
			request: dto.AnswerRequest{
				QuestionID:       "uid1",
				SelectedChoiceID: "choice2",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				question := entity.Question{
					UID:           entity.ToUID("uid1"),
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
					Choices: []entity.Choice{
						{
							UID:         entity.ToUID("choice1"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "正解",
							IsCorrect:   true,
						},
						{
							UID:         entity.ToUID("choice2"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "不正解",
							IsCorrect:   false,
						},
					},
				}
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(question, nil)
			},
			expectedResult: dto.AnswerResult{
				IsCorrect: false,
			},
		},
		{
			name: "異常系: 問題が見つからない",
			request: dto.AnswerRequest{
				QuestionID:       "not_found",
				SelectedChoiceID: "choice1",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindByID", mock.Anything, entity.ToUID("not_found")).Return(entity.Question{}, domainErrors.ErrQuestionNotFound)
			},
			errorContains: "指定された問題が見つかりません",
			errorType:     serviceErrors.TypeNotFound,
		},
		{
			name: "異常系: 無効な選択肢ID",
			request: dto.AnswerRequest{
				QuestionID:       "uid1",
				SelectedChoiceID: "invalid_choice",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				question := entity.Question{
					UID:           entity.ToUID("uid1"),
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
					Choices: []entity.Choice{
						{
							UID:         entity.ToUID("choice1"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "正解",
							IsCorrect:   true,
						},
						{
							UID:         entity.ToUID("choice2"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "不正解",
							IsCorrect:   false,
						},
					},
				}
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(question, nil)
			},
			errorContains: "無効な選択肢IDです",
			errorType:     serviceErrors.TypeValidation,
		},
		{
			name: "異常系: リポジトリでエラー発生",
			request: dto.AnswerRequest{
				QuestionID:       "uid1",
				SelectedChoiceID: "choice1",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(entity.Question{}, errors.New("DB error"))
			},
			errorContains: "問題の取得に失敗しました",
			errorType:     serviceErrors.TypeInternal,
		},
		{
			name: "異常系: 空の選択肢IDの場合",
			request: dto.AnswerRequest{
				QuestionID:       "uid1",
				SelectedChoiceID: "",
			},
			mockSetup: func(m *mocks.MockQuestionRepository) {
				question := entity.Question{
					UID:           entity.ToUID("uid1"),
					ReferenceCode: "REF001",
					Title:         "問題1",
					Content:       "問題1の内容",
					Choices: []entity.Choice{
						{
							UID:         entity.ToUID("choice1"),
							QuestionUID: entity.ToUID("uid1"),
							Content:     "正解",
							IsCorrect:   true,
						},
					},
				}
				m.On("FindByID", mock.Anything, entity.ToUID("uid1")).Return(question, nil)
			},
			errorContains: "無効な選択肢IDです",
			errorType:     serviceErrors.TypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockQuestionRepository)
			tt.mockSetup(mockRepo)
			
			service := NewQuestionService(mockRepo)
			result, err := service.Submit(context.Background(), tt.request)
			
			if tt.errorContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				serviceErr, ok := err.(*serviceErrors.ServiceError)
				assert.True(t, ok)
				assert.Equal(t, tt.errorType, serviceErr.Type)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}