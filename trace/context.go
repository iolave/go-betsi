package trace

import (
	"context"
	"reflect"
)

// ctx_trace_key is the key used to store
// and retrieve the trace from a context.
const ctx_trace_key = "trace"

// GetFromContext returns the trace from a context.
//
// If the context does not contain a trace or its
// type is not Trace, then a new Trace is zero-initialized.
func GetFromContext(ctx context.Context) Trace {
	trace := ctx.Value(ctx_trace_key)
	if trace == nil {
		return Trace{}
	}

	if reflect.TypeFor[Trace]() != reflect.TypeOf(trace) {
		return Trace{}
	}

	return trace.(Trace)
}

// SetContext returns a new context with the given trace.
func SetContext(ctx context.Context, trace Trace) context.Context {
	return context.WithValue(ctx, ctx_trace_key, trace)
}
