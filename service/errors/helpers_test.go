package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		cause       error
		expectedMsg string
	}{
		{
			name:        "原因エラーあり",
			message:     "リソースが見つかりません",
			cause:       errors.New("original error"),
			expectedMsg: "リソースが見つかりません",
		},
		{
			name:        "原因エラーなし",
			message:     "データが見つかりません",
			cause:       nil,
			expectedMsg: "データが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNotFoundError(tt.message, tt.cause)
			
			assert.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
			
			serviceErr := err
			assert.Equal(t, TypeNotFound, serviceErr.Type)
			assert.Equal(t, tt.message, serviceErr.Message)
			assert.Equal(t, tt.cause, serviceErr.Unwrap())
			assert.Empty(t, serviceErr.Field)
		})
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		field       string
		cause       error
		expectedMsg string
	}{
		{
			name:        "フィールドと原因エラーあり",
			message:     "無効な値です",
			field:       "email",
			cause:       errors.New("invalid format"),
			expectedMsg: "無効な値です",
		},
		{
			name:        "フィールドのみ",
			message:     "必須項目です",
			field:       "name",
			cause:       nil,
			expectedMsg: "必須項目です",
		},
		{
			name:        "メッセージのみ",
			message:     "バリデーションエラー",
			field:       "",
			cause:       nil,
			expectedMsg: "バリデーションエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.field, tt.cause)
			
			assert.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
			
			serviceErr := err
			assert.Equal(t, TypeValidation, serviceErr.Type)
			assert.Equal(t, tt.message, serviceErr.Message)
			assert.Equal(t, tt.field, serviceErr.Field)
			assert.Equal(t, tt.cause, serviceErr.Unwrap())
		})
	}
}

func TestNewInternalError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		cause       error
		expectedMsg string
	}{
		{
			name:        "原因エラーあり",
			message:     "内部エラーが発生しました",
			cause:       errors.New("database connection failed"),
			expectedMsg: "内部エラーが発生しました",
		},
		{
			name:        "原因エラーなし",
			message:     "処理に失敗しました",
			cause:       nil,
			expectedMsg: "処理に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewInternalError(tt.message, tt.cause)
			
			assert.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
			
			serviceErr := err
			assert.Equal(t, TypeInternal, serviceErr.Type)
			assert.Equal(t, tt.message, serviceErr.Message)
			assert.Equal(t, tt.cause, serviceErr.Unwrap())
			assert.Empty(t, serviceErr.Field)
		})
	}
}

func TestNewConflictError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		cause       error
		expectedMsg string
	}{
		{
			name:        "原因エラーあり",
			message:     "既に存在します",
			cause:       errors.New("duplicate key"),
			expectedMsg: "既に存在します",
		},
		{
			name:        "原因エラーなし",
			message:     "競合が発生しました",
			cause:       nil,
			expectedMsg: "競合が発生しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewConflictError(tt.message, tt.cause)
			
			assert.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
			
			serviceErr := err
			assert.Equal(t, TypeConflict, serviceErr.Type)
			assert.Equal(t, tt.message, serviceErr.Message)
			assert.Equal(t, tt.cause, serviceErr.Unwrap())
			assert.Empty(t, serviceErr.Field)
		})
	}
}