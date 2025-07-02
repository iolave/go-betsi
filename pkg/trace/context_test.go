package trace

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetFromContext(t *testing.T) {
	t.Run("should return an empty trace if the context is nil", func(t *testing.T) {
		want := Trace{}
		got := GetFromContext(nil)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetFromContext() = %v, want %v", got, want)
		}
	})

	t.Run("should return an empty trace if the context does not contain a trace", func(t *testing.T) {
		ctx := context.Background()
		want := Trace{}
		got := GetFromContext(ctx)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetFromContext() = %v, want %v", got, want)
		}
	})

	t.Run("should return an empty trace if the context contains a trace of the wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctx_trace_key, 1)
		want := Trace{}
		got := GetFromContext(ctx)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetFromContext() = %v, want %v", got, want)
		}
	})

	t.Run("should return the trace from the context", func(t *testing.T) {
		want := Trace{RequestID: uuid.NewString()}
		ctx := SetContext(t.Context(), want)
		got := GetFromContext(ctx)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetFromContext() = %v, want %v", got, want)
		}
	})
}
