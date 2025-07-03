package logger

import "github.com/pingolabscl/go-app/pkg/errors"

// Level is the level of the logger.
type Level int

const (
	// LEVEL_TRACE is the trace level.
	LEVEL_TRACE Level = (iota + 1) * 10
	// LEVEL_DEBUG is the debug level.
	LEVEL_DEBUG
	// LEVEL_INFO is the info level.
	LEVEL_INFO
	// LEVEL_WARN is the warn level.
	LEVEL_WARN
	// LEVEL_ERROR is the error level.
	LEVEL_ERROR
	// LEVEL_FATAL is the fatal level.
	LEVEL_FATAL
)

// String returns the string representation of the level
// in lower case.
func (l Level) String() string {
	switch l {
	case LEVEL_TRACE:
		return "trace"
	case LEVEL_DEBUG:
		return "debug"
	case LEVEL_INFO:
		return "info"
	case LEVEL_WARN:
		return "warn"
	case LEVEL_ERROR:
		return "error"
	case LEVEL_FATAL:
		return "fatal"
	default:
		return "unknown"
	}
}

// IsValid checks if the value is a valid level.
func (l Level) IsValid() bool {
	// Check if the value is a multiple of 10
	if l%10 != 0 {
		return false
	}

	// Check if the value is one on the defined levels.
	if l < LEVEL_TRACE || l > LEVEL_FATAL {
		return false
	}

	return true
}

// NewLevelFromString returns a level from a string.
// If the string is not a valid level, an error is returned.
func NewLevelFromString(s string) (Level, *errors.Error) {
	switch s {
	case "trace":
		return LEVEL_TRACE, nil
	case "debug":
		return LEVEL_DEBUG, nil
	case "info":
		return LEVEL_INFO, nil
	case "warn":
		return LEVEL_WARN, nil
	case "error":
		return LEVEL_ERROR, nil
	case "fatal":
		return LEVEL_FATAL, nil
	default:
		return LEVEL_TRACE, errors.NewWithName(
			"level_error",
			"invalid level",
		)
	}
}
