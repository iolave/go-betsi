package handlers

import (
	"net/http"

	"github.com/iolave/go-errors"
	"github.com/iolave/go-logger"
)

// NewNotFoundHandler returns a http.HandlerFunc that
// send a 404 error if the requested resource
// is not found. The response will be in JSON format and
// it follows the [github.com/iolave/go-errors.HTTPError]
// structure.
//
// Optionally, a logger can be passed to log the
// error. If no logger is passed, the error will
// not be logged.
func NewNotFoundHandler(logger logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.NewNotFoundError(
			"resource not found",
			nil,
		).(*errors.HTTPError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.StatusCode)
		w.Write(err.JSON())

		if logger != nil {
			logger.ErrorWithData(ctx, "resource_not_found", err, map[string]any{
				"path":   r.URL.Path,
				"method": r.Method,
			})
		}
	})
}
