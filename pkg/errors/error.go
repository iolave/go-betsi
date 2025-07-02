package errors

import (
	"encoding/json"
	"fmt"
)

// Error is an error struct that can be serialized to json.
// It contains an OriginalError property that can be used to
// wrap an error.
type Error struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Original error  `json:"original"`
}

// Error returns Original.Error() if Original is not nil and
// the concatenation of the Name and Message properties otherwise.
func (e *Error) Error() string {
	if e.Original != nil {
		return e.Original.Error()
	}

	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

// MarshalJSON returns the JSON representation of the error.
func (e Error) MarshalJSON() ([]byte, error) {
	if _, err := json.Marshal(e.Original); err != nil {
		e.Original = nil
	}

	return json.Marshal(struct {
		Name     string `json:"name"`
		Message  string `json:"message"`
		Original error  `json:"original"`
	}{
		Name:     e.Name,
		Message:  e.Message,
		Original: e.Original,
	})
}

// JSON returns the JSON representation of the error.
//
// If the Original property is not marshalable, it will
// be omitted from the output.
func (e Error) JSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New creates a new Error, it sets by default the name to "error".
func New(message string) *Error {
	return &Error{
		Name:    "error",
		Message: message,
	}
}

// NewWithName creates a new Error.
func NewWithName(name string, message string) *Error {
	return &Error{
		Name:    name,
		Message: message,
	}
}

// NewWithNameAndErr creates a new Error with the given name and error.
func NewWithNameAndErr(name string, message string, err error) *Error {
	return &Error{
		Name:     name,
		Message:  message,
		Original: err,
	}
}

// Wrap wraps an error with a new Error. If
// the error is of type *Error, it returns the
// original error.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return &Error{
		Name:     "error",
		Message:  err.Error(),
		Original: err,
	}
}
