package logger

import (
	"testing"
)

func TestLevel_String(t *testing.T) {
	t.Run("should return the string representation of the level", func(t *testing.T) {
		levels := []Level{
			LEVEL_TRACE,
			LEVEL_DEBUG,
			LEVEL_INFO,
			LEVEL_WARN,
			LEVEL_ERROR,
			LEVEL_FATAL,
			-1,
		}

		for _, level := range levels {
			got := level.String()

			if got != level.String() {
				t.Errorf("String() = %v, want %v", got, level.String())
			}
		}
	})
}
