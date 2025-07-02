package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/pingolabscl/go-app/errors"
	"github.com/pingolabscl/go-app/trace"
	"golang.org/x/mod/semver"
)

type Config struct {
	Name  string
	Level Level
}

type Logger struct {
	Config
	version string
}

func New(cfg Config) Logger {
	if cfg.Level == "" {
		cfg.Level = LEVEL_INFO
	}
	version := os.Getenv("VERSION")
	if !semver.IsValid(version) {
		version = fmt.Sprintf("%.7s", version)
	}
	return Logger{
		Config:  cfg,
		version: version,
	}
}

func (l Logger) Debug(ctx context.Context, msg string) {
	level := LEVEL_DEBUG
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, map[string]any{}, nil))
}

func (l Logger) DebugWithData(ctx context.Context, msg string, data map[string]any) {
	level := LEVEL_DEBUG
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, data, nil))
}

func (l Logger) Info(ctx context.Context, msg string) {
	level := LEVEL_INFO
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, map[string]any{}, nil))
}

func (l Logger) InfoWithData(ctx context.Context, msg string, data map[string]any) {
	level := LEVEL_INFO
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, data, nil))
}

func (l Logger) Warn(ctx context.Context, msg string) {
	level := LEVEL_WARN
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, map[string]any{}, nil))
}

func (l Logger) WarnWithData(ctx context.Context, msg string, data map[string]any) {
	level := LEVEL_WARN
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, data, nil))
}

func (l Logger) Error(ctx context.Context, msg string, err error) {
	level := LEVEL_ERROR
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, map[string]any{}, wrapError(err)))
}

func (l Logger) ErrorWithData(ctx context.Context, msg string, err error, data map[string]any) {
	level := LEVEL_ERROR
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, data, wrapError(err)))
}

// Fatal exits with code 1
func (l Logger) Fatal(ctx context.Context, msg string, err error) {
	level := LEVEL_FATAL
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, map[string]any{}, wrapError(err)))
	os.Exit(1)
}

// FatalWithData exits with code 1
func (l Logger) FatalWithData(ctx context.Context, msg string, err error, data map[string]any) {
	level := LEVEL_FATAL
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, data, wrapError(err)))
	os.Exit(1)
}

func (l Logger) shouldPrint(level Level) bool {
	funcLevel := level.toInt()
	settedLevel := l.Level.toInt()

	if funcLevel >= settedLevel {
		return true
	}

	return false
}

type entry struct {
	Timestamp int64          `json:"timestamp"`
	Level     Level          `json:"level"`
	Name      string         `json:"name"`
	Version   string         `json:"version"`
	Trace     trace.Trace    `json:"trace"`
	Error     error          `json:"error,omitempty"`
	Info      map[string]any `json:"info"`
	Msg       string         `json:"msg"`
}

func (l Logger) buildLogEntry(
	ctx context.Context,
	level Level,
	msg string,
	info map[string]any,
	err error,
) string {
	entry := entry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		Name:      l.Name,
		Version:   l.version,
		Trace:     trace.GetFromContext(ctx),
		Error:     err,
		Info:      info,
		Msg:       msg,
	}

	b, _ := json.Marshal(entry)
	return string(b)
}

func wrapError(err error) error {
	if err == nil {
		return errors.New("empty error")
	}

	if "*errors.errorString" != reflect.TypeOf(err).String() {
		return err
	}

	return errors.New(err.Error())
}
