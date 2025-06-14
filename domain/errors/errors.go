package errors

import "errors"

// ドメイン層のセンチネルエラー
var (
	ErrNoChoiceUID               = errors.New("choice UID is required")
	ErrMaxChoiceCount            = errors.New("choice count exceeds maximum limit")
	ErrNoChoice                  = errors.New("at least one choice is required")
	ErrInvalidCorrectChoiceCount = errors.New("exactly one correct choice is required")
	ErrQuestionNotFound          = errors.New("question not found")
	ErrChoiceNotFound            = errors.New("choice not found")
)
