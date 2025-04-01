package trace

import (
	"context"
	"reflect"
)

const traceKey = "trace"

type Trace struct {
	RequestID string `json:"requestId"`
}

func Get(ctx context.Context) Trace {
	ctxTrace := ctx.Value(traceKey)
	if ctxTrace == nil {
		return Trace{}
	}

	if reflect.TypeFor[Trace]() != reflect.TypeOf(ctxTrace) {
		return Trace{}
	}

	return ctxTrace.(Trace)
}

func Set(ctx context.Context, trace Trace) context.Context {
	return context.WithValue(ctx, traceKey, trace)
}
