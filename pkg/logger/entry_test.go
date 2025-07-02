package logger

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/pkg/errors"
	"github.com/pingolabscl/go-app/pkg/trace"
)

func Test_newEntry(t *testing.T) {
	t.Run("should return an entry with the given values", func(t *testing.T) {
		w := &Entry{
			Timestamp: 0,
			Level:     LEVEL_INFO,
			LevelStr:  LEVEL_INFO.String(),
			Name:      uuid.NewString(),
			Version:   uuid.NewString(),
			Trace:     trace.Trace{},
			Error:     nil,
			Info:      map[string]any{},
			Msg:       uuid.NewString(),
		}

		got := newEntry(
			nil,
			w.Name,
			w.Version,
			LEVEL_INFO,
			w.Msg,
			nil,
			w.Error,
		)
		got.Timestamp = 0

		if got.serialize() != w.serialize() {
			t.Errorf("newEntry() = %v, want %v", got.serialize(), w.serialize())
		}
	})
}

type nonMarshallableError struct {
	Fn func() `json:"fn"`
}

func (e nonMarshallableError) Error() string {
	return "error"
}

func Test_newEntry_serialize(t *testing.T) {
	t.Run("should return an json entry with a wrapped error cuz error is not marshalable", func(t *testing.T) {
		err := nonMarshallableError{
			Fn: func() {},
		}

		e := &Entry{
			Timestamp: 0,
			Level:     LEVEL_INFO,
			LevelStr:  LEVEL_INFO.String(),
			Name:      uuid.NewString(),
			Version:   uuid.NewString(),
			Trace:     trace.Trace{},
			Error:     err,
			Info:      map[string]any{},
			Msg:       uuid.NewString(),
		}

		w := fmt.Sprintf(
			`{"timestamp":%d,"level":"%s","name":"%s","version":"%s","trace":{},"error":%s,"info":{},"msg":"%s"}`,
			e.Timestamp,
			e.LevelStr,
			e.Name,
			e.Version,
			errors.Wrap(err).JSON(),
			e.Msg,
		)
		got := e.serialize()

		if got != w {
			t.Errorf("serialize() = %v, want %v", got, w)
		}
	})

	t.Run("should return an json entry with an empty info cuz info is not json marshalable", func(t *testing.T) {
		info := map[string]any{
			"fn": func() {},
		}

		e := &Entry{
			Timestamp: 0,
			Level:     LEVEL_INFO,
			LevelStr:  LEVEL_INFO.String(),
			Name:      uuid.NewString(),
			Version:   uuid.NewString(),
			Trace:     trace.Trace{},
			Error:     nil,
			Info:      info,
			Msg:       uuid.NewString(),
		}
		w := fmt.Sprintf(
			`{"timestamp":%d,"level":"%s","name":"%s","version":"%s","trace":{},"info":{},"msg":"%s"}`,
			e.Timestamp,
			e.LevelStr,
			e.Name,
			e.Version,
			e.Msg,
		)
		got := e.serialize()

		if got != w {
			t.Errorf("serialize() = %v, want %v", got, w)
		}
	})
}
