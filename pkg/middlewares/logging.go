package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/iolave/go-betsi/internal/utils"
	"github.com/iolave/go-logger"
)

// RequestLoggingMdwConfig holds the configuration for the request logging middleware.
type RequestLoggingMdwConfig struct {
	// Logger is the logger to use for logging requests.
	Logger logger.Logger

	// LogPath determines whether to log the request path.
	LogPath bool

	// LogPathParams determines whether to log the request path parameters.
	LogPathParams bool

	// LogQueryParams determines whether to log the request query parameters.
	LogQueryParams bool

	// LogJSONBody determines whether to log the request body for POST and PUT requests
	// with "application/json" content type.
	LogJSONBody bool
}

// NewRequestLoggingMdw creates a new request logging middleware.
//
// This middleware logs the start, success, or failure of each request.
// It can be configured to log the request path, path parameters, query parameters, and JSON body.
// It uses a custom response writer to capture the status code and any errors that occur during the request.
// The log messages are formatted as "<method>_<path>_<status>".
//
// Example:
//
//	r.Use(middlewares.NewRequestLoggingMdw(middlewares.RequestLoggingMdwConfig{
//		Logger:         l,
//		LogPath:        true,
//		LogPathParams:  true,
//		LogQueryParams: true,
//		LogJSONBody:    true,
//	}))
func NewRequestLoggingMdw(cfg RequestLoggingMdwConfig) func(next http.Handler) http.Handler {
	if cfg.Logger == nil {
		panic("logger cannot be nil")
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Create a new custom response writer
			// that stores the sent status code (and error if any)
			// when writing the header and body.
			w = &utils.ResponseWriter{
				Original: w,
			}

			l := cfg.Logger

			// Get the context from the request
			ctx := r.Context()
			rc := chi.RouteContext(ctx)

			// Determine if the request path belongs to a
			// pattern that has been registered with the router.
			// If not, we can skip the logging middleware.
			path := rc.Routes.Find(rc, r.Method, r.URL.Path)
			if path == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Convert the path to lower case and replaces
			// "/", "-", "{", and "}" with "_".
			if path[0] == '/' {
				path = path[1:]
			}
			path = strings.ToLower(path)
			path = strings.ReplaceAll(path, "/", "_")
			path = strings.ReplaceAll(path, "-", "_")
			path = strings.ReplaceAll(path, "{", "")
			path = strings.ReplaceAll(path, "}", "")

			// Build the log message for starting the request
			msg := fmt.Sprintf(
				"%s_%s_started",
				strings.ToLower(r.Method),
				path,
			)

			// data is a map of key/value pairs that will be
			// logged with the message. Here we'll store the
			// request path, path params, query params, and
			// request body if configured.
			data := map[string]any{}

			if cfg.LogPath {
				data["path"] = r.URL.Path
			}

			if cfg.LogPathParams {
				pathParams := map[string]string{}
				for _, k := range rc.URLParams.Keys {
					pathParams[k] = rc.URLParam(k)
				}
				data["pathParams"] = pathParams
			}

			if cfg.LogQueryParams {
				queryParams := map[string]string{}
				for k, v := range r.URL.Query() {
					queryParams[k] = v[0]
				}
				data["queryParams"] = queryParams
			}

			if cfg.LogJSONBody {
				if r.Method == "POST" || r.Method == "PUT" {
					if r.Header.Get("Content-Type") == "application/json" {
						buf, _ := io.ReadAll(r.Body)
						reader := io.NopCloser(bytes.NewBuffer(buf))
						r.Body = reader

						var body any
						err := json.Unmarshal(buf, &body)
						if err == nil {
							data["body"] = body
						} else {
							data["body"] = nil

						}
					}
				}
			}

			// Defers a function call to be executed after the function returns.
			// This function uses the custom response writer to determine
			// if the request was successful or not and then logs the message
			// accordingly.
			defer func() {
				w := w.(*utils.ResponseWriter)
				if w.SentErr == nil {
					// Logs the request succeeded message
					msg := fmt.Sprintf(
						"%s_%s_succeeded",
						strings.ToLower(r.Method),
						path,
					)
					l.InfoWithData(ctx, msg, data)
					return
				} else {
					// Logs the request failed message
					msg := fmt.Sprintf(
						"%s_%s_failed",
						strings.ToLower(r.Method),
						path,
					)
					l.ErrorWithData(ctx, msg, w.SentErr, data)
					return
				}
			}()

			// Logs the request started message
			l.InfoWithData(ctx, msg, data)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
