package middlewares

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/iolave/go-trace"
)

// newTraceMdw returns a chi middleware that checks
// if the request has trace headers and adds them
// to the context.
//
// If the request doesn't have trace headers, it
// will populate the context with a new trace.
func NewTraceMdw(mapHeaderToKey map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			t, _ := trace.GetFromHTTPRequest(r)

			for h, k := range mapHeaderToKey {
				v := r.Header.Get(h)
				t.Set(k, v)
			}

			reqId := t.Get("request_id")
			if reqId == "" {
				reqId = uuid.New().String()
			}

			t.SetHTTPHeaders(w.Header())
			ctx = t.SetInContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
