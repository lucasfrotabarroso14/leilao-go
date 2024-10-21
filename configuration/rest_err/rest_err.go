package rest_err

import (
	"leilao-go/internal/internal_error"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Err     string `json:"error"`
	Code    int    `json:"code"`
	Causes  []Causes
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func ConvertError(internalError *internal_error.InternalError) *RestErr {
	switch internalError.Err {
	case "bad_request":
		return NewBadRequestError(internalError.Error())
	case "not_found":
		return NewNotFoundError(internalError.Error())
	default:
		return NewInternalServerError(internalError.Error())
	}
}

func NewBadRequestError(message string, causes ...Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
		Causes:  nil,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_Server",
		Code:    http.StatusInternalServerError,
		Causes:  nil,
	}
}
