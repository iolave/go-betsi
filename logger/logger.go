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
)

type Config struct {
	Name  string
	Level Level
}

type Logger struct {
	Config
}

func New(cfg Config) Logger {
	if cfg.Level == "" {
		cfg.Level = LEVEL_INFO
	}

	return Logger{
		Config: cfg,
	}
}

func (l Logger) Debug(ctx context.Context, msg string) {
	level := LEVEL_DEBUG
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, nil))
}

func (l Logger) Info(ctx context.Context, msg string) {
	level := LEVEL_INFO
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, nil))
}

func (l Logger) Warn(ctx context.Context, msg string) {
	level := LEVEL_WARN
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, nil))
}

func (l Logger) Error(ctx context.Context, msg string, err error) {
	level := LEVEL_ERROR
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, wrapError(err)))
}

// Fatal exits with code 1
func (l Logger) Fatal(ctx context.Context, msg string, err error) {
	level := LEVEL_FATAL
	if !l.shouldPrint(level) {
		return
	}

	fmt.Println(l.buildLogEntry(ctx, level, msg, wrapError(err)))
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
	Timestamp int64       `json:"timestamp"`
	Level     Level       `json:"level"`
	Name      string      `json:"name"`
	Trace     trace.Trace `json:"trace"`
	Error     error       `json:"error,omitempty"`
	Msg       string      `json:"msg"`
}

func (l Logger) buildLogEntry(ctx context.Context, level Level, msg string, err error) string {
	entry := entry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		Name:      l.Name,
		Error:     err,
		Trace:     trace.Get(ctx),
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
