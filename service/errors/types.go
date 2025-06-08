package errors

type ErrorType string

const (
	TypeNotFound   ErrorType = "not_found"
	TypeValidation ErrorType = "validation"
	TypeConflict   ErrorType = "conflict"
	TypeInternal   ErrorType = "internal"
)

type ServiceError struct {
	Type    ErrorType
	Message string
	Field   string // バリデーションエラーの場合のフィールド名
	cause   error  // 元のエラー（外部に公開しない）
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) Unwrap() error {
	return e.cause
}