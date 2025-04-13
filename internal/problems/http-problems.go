package problems

import (
	"fmt"
)

type HTTPError struct {
	Code    int    // HTTP status code
	Message string // Error message
	Err     error  // Wrapped error (optional)
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func (e *HTTPError) Is(target error) bool {
	t, ok := target.(*HTTPError)
	if !ok {
		return false
	}
	return t.Code == e.Code
}

func NewHTTPError(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// // Common HTTP errors
// var (
// 	ErrBadRequest          = &HTTPError{Code: http.StatusBadRequest, Message: "bad request"}
// 	ErrUnauthorized        = &HTTPError{Code: http.StatusUnauthorized, Message: "unauthorized"}
// 	ErrForbidden           = &HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
// 	ErrNotFound            = &HTTPError{Code: http.StatusNotFound, Message: "not found"}
// 	ErrInternalServerError = &HTTPError{Code: http.StatusInternalServerError, Message: "internal server error"}
// )
