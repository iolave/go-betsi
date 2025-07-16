package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/pingolabscl/go-app/pkg/errors"
	"golang.org/x/mod/semver"
)

const (
	error_name_logger       = "logger_error"
	error_msg_invalid_level = "invalid level"
)

type Config struct {
	// Name is the name of the app.
	Name string
	// Level is the level of the logger.
	Level Level
}

// Logger is a json based logger that log messages
// appending the current timestamp, level, name, version,
// trace, error and info.
type Logger struct {
	// Config is the configuration of the logger.
	Config
	// version is the version of the app and is
	// retrieved from the environment variable VERSION.
	version string
}

// New returns a new logger with the given configuration
// and sets the version of the app from the environment
// variable VERSION. If the version is not a valid semver
// version, then it is set to the first 7 characters of the
// version. If version is an empty string, then it is set
// to develop.
//
// If the level is not valid, then an error is returned.
//
// The error returned is of type *errors.Error.
func New(cfg Config) (*Logger, error) {
	if !cfg.Level.IsValid() {
		return nil, errors.NewWithName(error_name_logger, error_msg_invalid_level)
	}

	version := os.Getenv("VERSION")
	if !semver.IsValid(version) {
		version = fmt.Sprintf("%.7s", version)
	}

	if version == "" {
		version = "develop"
	}

	return &Logger{
		Config:  cfg,
		version: version,
	}, nil
}

// log logs an entry and applies necessary filters to sanitize sensitive data.
func (l *Logger) log(entry *Entry) {
	// checks if the entry should be printed
	if entry.Level < l.Level {
		return
	}

	// serializes the entry
	msg := entry.serialize()

	// prints the entry
	fmt.Println(msg)
}

// Trace logs a message with the LEVEL_TRACE level.
func (l Logger) Trace(ctx context.Context, msg string) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_TRACE,
		msg,
		nil,
		nil,
	)
	l.log(e)
}

// TraceWithData logs a message with the LEVEL_TRACE level
// and the given data.
func (l *Logger) TraceWithData(ctx context.Context, msg string, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_TRACE,
		msg,
		data,
		nil,
	)
	l.log(e)
}

// Debug logs a message with the LEVEL_DEBUG level.
func (l Logger) Debug(ctx context.Context, msg string) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_DEBUG,
		msg,
		nil,
		nil,
	)
	l.log(e)
}

// DebugWithData logs a message with the LEVEL_DEBUG level
// and the given data.
func (l *Logger) DebugWithData(ctx context.Context, msg string, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_DEBUG,
		msg,
		data,
		nil,
	)
	l.log(e)
}

// Info logs a message with the LEVEL_INFO level.
func (l Logger) Info(ctx context.Context, msg string) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_INFO,
		msg,
		nil,
		nil,
	)
	l.log(e)
}

// InfoWithData logs a message with the LEVEL_INFO level
// and the given data.
func (l *Logger) InfoWithData(ctx context.Context, msg string, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_INFO,
		msg,
		data,
		nil,
	)
	l.log(e)
}

// Warn logs a message with the LEVEL_WARN level.
func (l *Logger) Warn(ctx context.Context, msg string, err error) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_WARN,
		msg,
		nil,
		err,
	)
	l.log(e)
}

// WarnWithData logs a message with the LEVEL_WARN level
// and the given data.
func (l *Logger) WarnWithData(ctx context.Context, msg string, err error, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_WARN,
		msg,
		data,
		err,
	)
	l.log(e)
}

// Error logs a message with the LEVEL_ERROR level.
func (l *Logger) Error(ctx context.Context, msg string, err error) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_ERROR,
		msg,
		nil,
		err,
	)
	l.log(e)
}

// ErrorWithData logs a message with the LEVEL_ERROR level
// and the given data.
func (l *Logger) ErrorWithData(ctx context.Context, msg string, err error, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_ERROR,
		msg,
		data,
		err,
	)
	l.log(e)
}

// Fatal logs a message with the LEVEL_FATAL level
// and exits with code 1.
func (l *Logger) Fatal(ctx context.Context, msg string, err error) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_FATAL,
		msg,
		nil,
		err,
	)
	l.log(e)
	os.Exit(1)
}

// FatalWithData logs a message with the LEVEL_FATAL level
// with the given data and exits with code 1.
func (l *Logger) FatalWithData(ctx context.Context, msg string, err error, data map[string]any) {
	e := newEntry(
		ctx,
		l.Name,
		l.version,
		LEVEL_FATAL,
		msg,
		data,
		err,
	)
	l.log(e)
	os.Exit(1)
}
