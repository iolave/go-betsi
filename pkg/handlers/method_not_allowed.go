package handlers

import (
	"net/http"

	"github.com/iolave/go-errors"
	"github.com/iolave/go-logger"
)

// NewMethodNotAllowedHandler returns a http.HandlerFunc that
// send a 405 error if the requested method is not allowed.
// The response will be in JSON format and it follows the
// [github.com/iolave/go-errors.HTTPError] structure.
//
// Optionally, a logger can be passed to log the
// error. If no logger is passed, the error will
// not be logged.
func NewMethodNotAllowedHandler(logger logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.NewHTTPError(
			http.StatusMethodNotAllowed,
			"method_not_allowed_error",
			"method not allowed",
			nil,
		).(*errors.HTTPError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.StatusCode)
		w.Write(err.JSON())

		if logger != nil {
			logger.ErrorWithData(ctx, "method_not_allowed", err, map[string]any{
				"path":   r.URL.Path,
				"method": r.Method,
			})
		}
	})
}
