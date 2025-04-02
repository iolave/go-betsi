package errors

import (
	"fmt"
)

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s", err.Message)
}

func New(msg string) *Error {
	return &Error{
		Name:    "error",
		Message: msg,
	}
}
