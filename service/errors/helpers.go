package errors

// NewNotFoundError creates a new ServiceError with NotFound type
func NewNotFoundError(message string, cause error) *ServiceError {
	return &ServiceError{
		Type:    TypeNotFound,
		Message: message,
		cause:   cause,
	}
}

// NewValidationError creates a new ServiceError with Validation type
func NewValidationError(message string, field string, cause error) *ServiceError {
	return &ServiceError{
		Type:    TypeValidation,
		Message: message,
		Field:   field,
		cause:   cause,
	}
}

// NewInternalError creates a new ServiceError with Internal type
func NewInternalError(message string, cause error) *ServiceError {
	return &ServiceError{
		Type:    TypeInternal,
		Message: message,
		cause:   cause,
	}
}

// NewConflictError creates a new ServiceError with Conflict type
func NewConflictError(message string, cause error) *ServiceError {
	return &ServiceError{
		Type:    TypeConflict,
		Message: message,
		cause:   cause,
	}
}