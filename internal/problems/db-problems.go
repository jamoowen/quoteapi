package problems

import "fmt"

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s", e.Message)
	}
	return fmt.Sprintf("Not Found Error")
}
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}
