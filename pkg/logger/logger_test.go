package logger

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/pkg/errors"
)

// Run a fork test that may crash using os.exit.
func RunForkTest(t *testing.T, testName string, id string) (string, string, error) {
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%v", testName))
	cmd.Env = append(os.Environ(), fmt.Sprintf("FORK=%s", id))

	var stdoutB, stderrB bytes.Buffer
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB

	err := cmd.Run()

	return stdoutB.String(), stderrB.String(), err
}

func TestNewLogger(t *testing.T) {
	t.Run("should return a new logger with the given configuration", func(t *testing.T) {
		cfg := Config{
			Name:  "test",
			Level: LEVEL_INFO,
		}
		want := &Logger{
			Config:  cfg,
			version: "develop",
		}

		got, err := New(cfg)

		if err != nil {
			t.Errorf("New() = (%v, %v), want (%v, <nil>)", got, err, want)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("New() = (%v, %v), want (%v, <nil>)", got.Config, err, want.Config)
		}
	})

	t.Run("should return an error if the level is not valid", func(t *testing.T) {
		levels := []Level{
			-1,
			-1000,
			1000,
		}

		for _, level := range levels {
			w := errors.NewWithName(error_name_logger, error_msg_invalid_level)
			cfg := Config{
				Name:  "test",
				Level: level,
			}

			l, got := New(cfg)

			if !reflect.DeepEqual(got, w) {
				t.Errorf("New() = (%v, %v), want (<nil>, %v)", l, got, w)
			}
		}
	})
}

func TestLogger(t *testing.T) {
	logger := Logger{
		Config: Config{
			Name:  uuid.NewString(),
			Level: LEVEL_DEBUG,
		},
		version: uuid.NewString(),
	}

	t.Run("should not log a trace message", func(t *testing.T) {
		ctx := t.Context()
		logger.Trace(ctx, "trace")
	})
	t.Run("should log a trace message with data", func(t *testing.T) {
		ctx := t.Context()
		logger.TraceWithData(ctx, "trace", map[string]any{})
	})
	t.Run("should log a debug message", func(t *testing.T) {
		ctx := t.Context()
		logger.Debug(ctx, "debug")
	})
	t.Run("should log a debug message with data", func(t *testing.T) {
		ctx := t.Context()
		logger.DebugWithData(ctx, "debug", map[string]any{})
	})
	t.Run("should log an info message", func(t *testing.T) {
		ctx := t.Context()
		logger.Info(ctx, "info")
	})
	t.Run("should log an info message with data", func(t *testing.T) {
		ctx := t.Context()
		logger.InfoWithData(ctx, "info", map[string]any{})
	})
	t.Run("should log a warn message", func(t *testing.T) {
		ctx := t.Context()
		logger.Warn(ctx, "warn", nil)
	})
	t.Run("should log a warn message with data", func(t *testing.T) {
		ctx := t.Context()
		logger.WarnWithData(ctx, "warn", nil, map[string]any{})
	})
	t.Run("should log an error message", func(t *testing.T) {
		ctx := t.Context()
		logger.Error(ctx, "error", nil)
	})
	t.Run("should log an error message with data", func(t *testing.T) {
		ctx := t.Context()
		logger.ErrorWithData(ctx, "error", nil, map[string]any{})
	})
	t.Run("should log a fatal message", func(t *testing.T) {
		forkId := "fork-fatal"
		switch os.Getenv("FORK") {
		case forkId:
			logger.Fatal(t.Context(), "fatal", nil)
		case "":
			_, _, err := RunForkTest(t, "TestLogger", forkId)

			if err == nil {
				t.Errorf("expected Fatal() to exit")
				return
			}

			if err.Error() != "exit status 1" {
				t.Errorf("Fatal() = %v, want %v", err.Error(), "exit status 1")
			}
		}
	})
	t.Run("should log a fatal message with data", func(t *testing.T) {
		forkId := "fork-fatal-with-data"
		if os.Getenv("FORK") == forkId {
			logger.FatalWithData(t.Context(), "fatal", nil, nil)
		}
		if os.Getenv("FORK") == "" {
			_, _, err := RunForkTest(t, "TestLogger", forkId)

			if err == nil {
				t.Errorf("FatalWithData() = <nil>, want an error")
				return
			}

			if err.Error() != "exit status 1" {
				t.Errorf("FatalWithData() = %v, want %v", err.Error(), "exit status 1")
			}
		}
	})
}
