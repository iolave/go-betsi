package goapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/pingolabscl/go-app/errors"
)

type AppRequest struct {
	app     *App
	Request *http.Request
	writer  http.ResponseWriter
}

// SendError sends an error response, following the
// errors.HTTP structure.
func (req *AppRequest) SendError(err error) {
	ctx := req.Request.Context()
	if err == nil {
		err = errors.NewInternalServerError("missing error", "unknown cause")
	}

	if reflect.TypeFor[*errors.HTTP]() != reflect.TypeOf(err) {
		err = errors.NewInternalServerError("internal error", err.Error())
	}

	req.writer.Header().Set("content-type", "application/json")
	req.writer.WriteHeader(err.(*errors.HTTP).Status)
	b, _ := json.Marshal(err.(*errors.HTTP))
	req.writer.Write(b)
	req.app.Logger.ErrorWithData(ctx, "handler_failed", err, map[string]any{
		"path": req.Request.URL.Path,
	})
}

// SendJSON sends an OK json response. If the value
// param failed to be marshaled, an error will be sent.
func (req *AppRequest) SendJSON(result any) {
	ctx := req.Request.Context()
	b, err := json.Marshal(result)
	if err != nil {
		req.SendError(errors.NewInternalServerError(
			"failed to serialize result",
			"result might contain invalid json content",
		))

		return
	}

	req.writer.Header().Set("content-type", "application/json")
	req.writer.WriteHeader(http.StatusOK)
	req.writer.Write(b)
	req.app.Logger.InfoWithData(ctx, "handler_success", map[string]any{
		"path": req.Request.URL.Path,
	})
}

type Handler func(AppRequest)

func newHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app := GetFromContext(ctx)
		app.Logger.InfoWithData(ctx, "handler_started", map[string]any{
			"path": r.URL.Path,
		})
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newHealthcheckHandler() http.HandlerFunc {
	startTime := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(map[string]any{
			"uptime": time.Since(startTime).
				Truncate(time.Second).
				Seconds(),
		})
		w.Write(b)
	})

}

func newNotFoundHandler(app *App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.NewNotFoundError(
			"resource not found",
			fmt.Sprintf("path %s not found", r.URL.Path),
		)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.Status)
		b, _ := json.Marshal(err)
		w.Write(b)
		app.Logger.ErrorWithData(ctx, "resource_not_found", err, map[string]any{
			"path": r.URL.Path,
		})
	})
}
