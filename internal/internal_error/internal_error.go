package internal_error

type InternalError struct {
	Message string
	Err     string
}

func NewNotFoundError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "not found",
	}
}

func (ie *InternalError) Error() string {
	return ie.Message
}

func NewInternalServerError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "internal_server_error",
	}
}

func NewBadRequestError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "internal_bad_request",
	}
}
