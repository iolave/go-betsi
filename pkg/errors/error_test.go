package errors

import (
	"errors"
	"reflect"
	"testing"
)

type AnyError struct {
	Fn func() `json:"fn"`
}

func (e *AnyError) Error() string {
	return "error"
}

func TestError(t *testing.T) {
	t.Run("should return new error", func(t *testing.T) {
		want := Error{
			Name:    "error",
			Message: "error",
		}
		got := New("error")

		if reflect.DeepEqual(got, want) {
			t.Errorf("New() = %v, want %v", got, want)
		}
	})

	t.Run("should return new error with custom name", func(t *testing.T) {
		customName := "custom_error"
		msg := "error"
		want := Error{
			Name:    customName,
			Message: msg,
		}
		got := NewWithName(customName, msg)

		if reflect.DeepEqual(got, want) {
			t.Errorf("NewWithName() = %v, want %v", got, want)
		}
	})

	t.Run("should wrap an error", func(t *testing.T) {
		toWrap := errors.New("normal error")
		want := Error{
			Name:    "error",
			Message: toWrap.Error(),
		}
		got := Wrap(toWrap)
		if reflect.DeepEqual(got, want) {
			t.Errorf("Wrap() = %v, want %v", got, want)
		}
	})

	t.Run("should return error message", func(t *testing.T) {
		want := "error: error"
		got := New("error").Error()
		if got != want {
			t.Errorf("Error() = %v, want %v", got, want)
		}
	})

	t.Run("should return original error message when Original is not nil", func(t *testing.T) {
		original := errors.New("original error")
		want := original.Error()
		got := Wrap(original).Error()
		if got != want {
			t.Errorf("Error() = %v, want %v", got, want)
		}
	})

	t.Run("should return json error", func(t *testing.T) {
		want := `{"name":"error","message":"error","original":null}`
		got := New("error").JSON()
		if got != want {
			t.Errorf("JSON() = %v, want %v", got, want)
		}
	})

	t.Run("should return json error when original is not marshalable", func(t *testing.T) {
		want := `{"name":"error","message":"error","original":null}`
		someErr := &AnyError{
			Fn: func() {},
		}
		got := Wrap(someErr).JSON()
		if got != want {
			t.Errorf("JSON() = %v, want %v", got, want)
		}
	})

	t.Run("should return a nil error when nil is passed to Wrap", func(t *testing.T) {
		got := Wrap(nil)
		if got != nil {
			t.Errorf("Wrap() = %v, want %v", got, nil)
		}
	})

	t.Run("should return original error when Wrap is called with an error that is already an Error", func(t *testing.T) {
		original := New("error")
		got := Wrap(original)
		if got != original {
			t.Errorf("Wrap() = %v, want %v", got, original)
		}
	})
}

func TestNewWithNameAndErr(t *testing.T) {
	t.Run("should return new error with the given name and error", func(t *testing.T) {
		want := Error{
			Name:     "error",
			Message:  "error",
			Original: errors.New("error"),
		}
		got := NewWithNameAndErr("error", "error", errors.New("error"))
		if reflect.DeepEqual(got, want) {
			t.Errorf("NewWithNameAndErr() = %v, want %v", got, want)
		}
	})
}
