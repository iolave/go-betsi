package goapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// Context returns the context within the http request.
// Is the same as ar.Request.Context().
func (ar *AppRequest) Context() context.Context {
	if ar.Request == nil {
		return nil
	}

	return ar.Request.Context()
}

// ReadJSONBody unmarshals a request json body into v.
// It also does struct validation using go-playground/validator
// rules.
func (ar *AppRequest) ReadJSONBody(v any) error {
	b, err := io.ReadAll(ar.Request.Body)
	if err != nil {
		return errors.NewInternalServerError("failed to read request body", err.Error())
	}
	if err := json.Unmarshal(b, v); err != nil {
		return errors.NewBadRequestError("invalid json content", err.Error())
	}

	valueOf := reflect.ValueOf(v)
	var valueToValidate any
	if valueOf.Kind() == reflect.Ptr {
		valueOf := reflect.ValueOf(valueOf.Elem().Interface())
		valueToValidate = valueOf.Interface()
	} else {
		valueToValidate = v
	}

	if reflect.TypeOf(valueToValidate).Kind() == reflect.Slice {
		valueOf := reflect.ValueOf(valueToValidate)
		length := valueOf.Len()
		for i := 0; i < length; i++ {
			valueToValidate := valueOf.Index(i).Interface()
			if err := ar.app.validator.Struct(valueToValidate); err != nil {
				return errors.NewInternalServerError(
					"unable to send response, response didn't passed validation",
					err.Error(),
				)
			}
		}
	} else {
		if err := ar.app.validator.Struct(valueToValidate); err != nil {
			return errors.NewInternalServerError(
				"unable to send response, response didn't passed validation",
				err.Error(),
			)
		}
	}

	return nil
}

// SendError sends an error response, following the
// errors.HTTP structure.
//   - If the error is not of type *errors.HTTP, an internal server
//     error will be sent instead.
func (ar *AppRequest) SendError(err error) {
	ctx := ar.Context()
	if err == nil {
		err = errors.NewInternalServerError("missing error", "unknown cause")
	}

	if reflect.TypeFor[*errors.HTTP]() != reflect.TypeOf(err) {
		err = errors.NewInternalServerError("internal error", err.Error())
	}

	ar.writer.Header().Set("content-type", "application/json")
	ar.writer.WriteHeader(err.(*errors.HTTP).Status)
	b, _ := json.Marshal(err.(*errors.HTTP))
	ar.writer.Write(b)
	ar.app.Logger.ErrorWithData(ctx, "handler_failed", err, map[string]any{
		"path": ar.Request.URL.Path,
	})
}

// SendJSON sends an OK http response with json content. Make
// sure the result param is a struct.
//
//   - If go-playground/validator validation requirements are
//     not met, an internal server error will be sent instead.
//
//   - if the result param failed to be marshaled, an internal
//     server error will be sent instead.
//
//   - If the result param is a map, go-playground/validator
//     validation requirements are not met, and therefore an
//     internal server error will be sent instead.
func (ar *AppRequest) SendJSON(result any) {
	ctx := ar.Context()

	if reflect.TypeOf(result).Kind() == reflect.Slice {
		valueOf := reflect.ValueOf(result)
		length := valueOf.Len()
		for i := 0; i < length; i++ {
			v := valueOf.Index(i).Interface()

			if err := ar.app.validator.Struct(v); err != nil {
				ar.SendError(errors.NewInternalServerError(
					"unable to send response, response didn't passed validation",
					err.Error(),
				))

				return
			}
		}

	} else {
		if err := ar.app.validator.Struct(result); err != nil {
			ar.SendError(errors.NewInternalServerError(
				"unable to send response, response didn't passed validation",
				err.Error(),
			))

			return
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		ar.SendError(errors.NewInternalServerError(
			"failed to serialize result",
			"result might contain invalid json content",
		))

		return
	}

	ar.writer.Header().Set("content-type", "application/json")
	ar.writer.WriteHeader(http.StatusOK)
	ar.writer.Write(b)
	ar.app.Logger.InfoWithData(ctx, "handler_success", map[string]any{
		"path": ar.Request.URL.Path,
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

func newMethodNotAllowedHandler(app *App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.NewMethodNotAllowedError(
			"method not allowed",
			fmt.Sprintf("method %s not allowed", r.Method),
		)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(err.Status)
		b, _ := json.Marshal(err)
		w.Write(b)
		app.Logger.ErrorWithData(ctx, "method_not_allowed", err, map[string]any{
			"path":   r.URL.Path,
			"method": r.Method,
		})
	})
}
