package goapp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/pkg/trace"
)

type mdwHandler func(http.Handler) http.Handler

// newXPoweredByMdw returns a chi middleware that sets
// the X-Powered-By header to a value that identifies
// this app follows the pingolabscl standard and thus,
// errors can be parsed and handled.
func newXPoweredByMdw() mdwHandler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			w.Header().Set("x-powered-by", "pingolabs.cl")
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func newRequestIdMdw() mdwHandler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestId := r.Header.Get("x-request-id")
			if requestId == "" {
				requestId = uuid.New().String()
			}

			ctx = trace.SetContext(ctx, trace.Trace{
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
