package logger

import (
	"context"
	"encoding/json"
	"time"

	anyutils "github.com/pingolabscl/go-app/internal/any_utils"
	"github.com/pingolabscl/go-app/pkg/errors"
	"github.com/pingolabscl/go-app/pkg/trace"
)

// Entry is the json representation of a log entry.
type Entry struct {
	Timestamp int64          `json:"timestamp"`
	Level     Level          `json:"-"`
	LevelStr  string         `json:"level"`
	Name      string         `json:"name"`
	Version   string         `json:"version"`
	Trace     trace.Trace    `json:"trace"`
	Error     error          `json:"error,omitempty"`
	Info      map[string]any `json:"info"`
	Msg       string         `json:"msg"`
}

// newEntry returns a new entry.
//
// It computes the trace from the context.
func newEntry(
	ctx context.Context,
	name string,
	version string,
	level Level,
	msg string,
	info map[string]any,
	err error,
) *Entry {
	if info == nil {
		info = map[string]any{}
	}

	return &Entry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		LevelStr:  level.String(),
		Name:      name,
		Version:   version,
		Trace:     trace.GetFromContext(ctx),
		Error:     err,
		Info:      info,
		Msg:       msg,
	}
}

// serialize serializes the entry to a json string.
//
// If the info are not serializable it will
// be set to it's zero value.
//
// If the error is not json marshalable, then it is
// wrapped in an error.
func (e *Entry) serialize() string {
	// Remove any circular references
	e.Info = anyutils.DestroyCircular(e.Info).(map[string]any)

	// If the entry has an error, check if it is
	// json marshalable and if not, wrap it in an error
	// so it can be serialized to something else other
	// than {}.
	if e.Error != nil {
		b, jsonErr := json.Marshal(e.Error)
		if jsonErr != nil || string(b) == "{}" {
			e.Error = errors.Wrap(e.Error)
		}
	}

	// Serialize the entry to json. If it succeeds,
	// return the string. If it fails, it means the
	// info property is not json marshalable, so
	// set it to an empty map.
	b, err := json.Marshal(e)
	if err == nil {
		return string(b)
	}
	e.Info = map[string]any{}
	b, _ = json.Marshal(e)
	return string(b)
}
