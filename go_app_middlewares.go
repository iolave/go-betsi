package goapp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/pkg/errors"
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
			w.Header().Set("x-powered-by", "pingolabs.cl")
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// newTraceMdw returns a chi middleware that checks
// if the request has trace headers and adds them
// to the context.
//
// If the request doesn't have trace headers, it
// will populate the context with a new trace.
func newTraceMdw() mdwHandler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tr := trace.NewFromHTTPHeaders(r.Header)

			// Check/set request id
			if tr.RequestID == "" {
				tr.RequestID = uuid.New().String()
			}
			tr.SetHTTPHeaders(w.Header())

			ctx = trace.SetContext(ctx, tr)

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

func NewAcceptVersionHandler(dflt Handler, handlers map[string]Handler) Handler {
	return func(ar AppRequest) {
		var handler Handler

		acceptVersion := ar.Request.Header.Get("Accept-Version")
		if acceptVersion != "" {
			foundHandler, ok := handlers[acceptVersion]
			if !ok {
				ar.SendError(errors.NewBadRequestError("invalid accept version header", nil))
				return
			}
			handler = foundHandler
		} else {
			handler = dflt
		}

		if handler == nil {
			ar.SendError(errors.NewInternalServerError(
				"handler not implemented",
				nil,
			))
			return
		}

		handler(ar)
	}

}
