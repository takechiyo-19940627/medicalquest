package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ServiceError
		expected string
	}{
		{
			name: "メッセージのみ",
			err: &ServiceError{
				Type:    TypeNotFound,
				Message: "データが見つかりません",
			},
			expected: "データが見つかりません",
		},
		{
			name: "原因エラー付き",
			err: &ServiceError{
				Type:    TypeInternal,
				Message: "処理に失敗しました",
				cause:   errors.New("database error"),
			},
			expected: "処理に失敗しました",
		},
		{
			name: "フィールド付き",
			err: &ServiceError{
				Type:    TypeValidation,
				Message: "無効な値です",
				Field:   "email",
			},
			expected: "無効な値です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestServiceError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		err      *ServiceError
		expected error
	}{
		{
			name: "原因エラーあり",
			err: &ServiceError{
				Type:    TypeInternal,
				Message: "エラー",
				cause:   errors.New("original error"),
			},
			expected: errors.New("original error"),
		},
		{
			name: "原因エラーなし",
			err: &ServiceError{
				Type:    TypeNotFound,
				Message: "エラー",
				cause:   nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unwrapped := tt.err.Unwrap()
			if tt.expected == nil {
				assert.Nil(t, unwrapped)
			} else {
				assert.Equal(t, tt.expected.Error(), unwrapped.Error())
			}
		})
	}
}

func TestErrorTypes(t *testing.T) {
	assert.Equal(t, ErrorType("not_found"), TypeNotFound)
	assert.Equal(t, ErrorType("validation"), TypeValidation)
	assert.Equal(t, ErrorType("conflict"), TypeConflict)
	assert.Equal(t, ErrorType("internal"), TypeInternal)
}