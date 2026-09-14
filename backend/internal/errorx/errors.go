package errorx

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CodeError) Error() string {
	return e.Message
}

func NewCodeError(code int, message string) *CodeError {
	return &CodeError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrNotFound         = NewCodeError(http.StatusNotFound, "resource not found")
	ErrBadRequest       = NewCodeError(http.StatusBadRequest, "bad request")
	ErrVersionConflict  = NewCodeError(http.StatusConflict, "version conflict: resource has been modified")
	ErrInvalidTransition = NewCodeError(http.StatusBadRequest, "invalid status transition")
	ErrDeleteNonDraft   = NewCodeError(http.StatusBadRequest, "can only delete draft workflows")
	ErrInternalServer   = NewCodeError(http.StatusInternalServerError, "internal server error")
)

func NewNotFoundError(msg string) *CodeError {
	return NewCodeError(http.StatusNotFound, msg)
}

func NewBadRequestError(msg string) *CodeError {
	return NewCodeError(http.StatusBadRequest, msg)
}

func NewInternalError(msg string) *CodeError {
	return NewCodeError(http.StatusInternalServerError, msg)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error) {
	if codeErr, ok := err.(*CodeError); ok {
		httpx.WriteJson(w, codeErr.Code, &ErrorResponse{
			Code:    codeErr.Code,
			Message: codeErr.Message,
		})
		return
	}

	httpx.WriteJson(w, http.StatusInternalServerError, &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}
