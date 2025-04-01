package goapp

import (
	"encoding/json"
	"net/http"
	"reflect"

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
	req.app.Logger.Error(ctx, "handler_failed", err)
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
	req.app.Logger.Info(ctx, "handler_success")
}

type Handler func(AppRequest)

func newHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app := GetFromContext(ctx)
		app.Logger.Info(ctx, "handler_started")
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
