package goapp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/trace"
)

type mdwHandler func(http.Handler) http.Handler

func newRequestIdMdw() mdwHandler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestId := r.Header.Get("x-request-id")
			if requestId == "" {
				requestId = uuid.New().String()
			}

			ctx = trace.Set(ctx, trace.Trace{
				RequestID: requestId,
			})

			w.Header().Set("x-request-id", requestId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func newAppContextMdw(app *App) mdwHandler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = setContext(ctx, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
