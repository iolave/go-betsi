package goapp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/pingolabscl/go-app/pkg/errors"
)

type AppRequest struct {
	App     *App
	Request *http.Request
	writer  http.ResponseWriter
}

// Context returns the context within the http request.
// Is the same as ar.Request.Context().
func (ar *AppRequest) Context() context.Context {
	if ar.Request == nil {
		return nil
	}

	return ar.Request.Context()
}

// ReadJSONBody unmarshals a request json body into v.
// It also validates v using go-playground/validator
// rules.
//
// error is of type *errors.HTTPError.
func (ar *AppRequest) ReadJSONBody(v any) error {
	b, err := io.ReadAll(ar.Request.Body)
	if err != nil {
		return errors.NewInternalServerError(
			"failed to read request body",
			errors.Wrap(err),
		)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.NewBadRequestError("request body is not valid json", errors.Wrap(err))
	}

	valueOf := reflect.ValueOf(v)
	var valueToValidate any
	if valueOf.Kind() == reflect.Ptr {
		valueOf := reflect.ValueOf(valueOf.Elem().Interface())
		valueToValidate = valueOf.Interface()
	} else {
		valueToValidate = v
	}

	if err := ar.App.recursiveValidation(valueToValidate); err != nil {
		return errors.NewInternalServerError(
			"client response doesn't meet validation rules",
			err,
		)
	}

	return nil
}

// SendError sends an error response, following the
// errors.HTTPError structure.
//   - If the error is not of type *errors.HTTPError, an internal server
//     error will be sent instead.
func (ar *AppRequest) SendError(err error) {
	ctx := ar.Context()
	if err == nil {
		err = errors.NewInternalServerError("nil error sent", nil)
	}

	switch err.(type) {
	case *errors.Error:
		err = errors.NewInternalServerError("internal error", err)
	default:
		err = errors.NewInternalServerError("internal error", errors.Wrap(err))
	}

	ar.writer.Header().Set("content-type", "application/json")
	ar.writer.WriteHeader(err.(*errors.HTTPError).StatusCode)
	b, _ := json.Marshal(err)
	ar.writer.Write(b)
	ar.App.Logger.ErrorWithData(ctx, "handler_failed", err, map[string]any{
		"path": ar.Request.URL.Path,
	})
}

// SendJSON sends an OK http response with json content. Make
// sure the result param either a struct, slice, or map. Otherwise,
// an internal server error will be sent.
//
//   - If go-playground/validator validation requirements are
//     not met, an internal server error will be sent instead.
//
//   - if the result param failed to be marshaled, an internal
//     server error will be sent instead.
func (ar *AppRequest) SendJSON(v any) {
	ctx := ar.Context()

	allowedKinds := []reflect.Kind{
		reflect.Slice,
		reflect.Map,
		reflect.Struct,
	}
	if kind := reflect.TypeOf(v).Kind(); !slices.Contains(allowedKinds, kind) {
		ar.SendError(errors.NewInternalServerError(
			"unable to send response, response is not a struct, slice or map",
			nil,
		))

		return
	}

	if err := ar.App.recursiveValidation(v); err != nil {
		ar.SendError(errors.NewInternalServerError(
			"unable to send response, response doesn't meet validation rules",
			err,
		))

		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ar.SendError(errors.NewInternalServerError(
			"failed to serialize response",
			errors.Wrap(err),
		))

		return
	}

	ar.writer.Header().Set("content-type", "application/json")
	ar.writer.WriteHeader(http.StatusOK)
	ar.writer.Write(b)
	ar.App.Logger.InfoWithData(ctx, "handler_success", map[string]any{
		"path": ar.Request.URL.Path,
	})
}

type Handler func(ar AppRequest)

func (app *App) newHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		app.Logger.InfoWithData(ctx, "handler_started", map[string]any{
			"path": r.URL.Path,
		})

		handler.ServeHTTP(w, r)
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
		err := errors.NewNotFoundError("resource not found", nil)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.(*errors.HTTPError).StatusCode)
		b, _ := json.Marshal(err)
		w.Write(b)
		app.Logger.ErrorWithData(ctx, "resource_not_found", err, map[string]any{
			"path": r.URL.Path,
		})
	})
}

func newMethodNotAllowedHandler(app *App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.NewMethodNotAllowedError("method not allowed", nil)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.(*errors.HTTPError).StatusCode)
		b, _ := json.Marshal(err)
		w.Write(b)
		app.Logger.ErrorWithData(ctx, "method_not_allowed", err, map[string]any{
			"path":   r.URL.Path,
			"method": r.Method,
		})
	})
}
