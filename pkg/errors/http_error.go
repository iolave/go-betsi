package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	anyutils "github.com/pingolabscl/go-app/internal/any_utils"
)

type HTTPError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	Handled    bool   `json:"handled"`
	Err        any    `json:"error"`
}

var _ JSONError = &HTTPError{}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

// NewHTTPError creates a new HTTPError.
//
// If the error is a built-in error, it will be wrapped
// using this package's Wrap function.
func NewHTTPError(statusCode int, name, message string, err any) error {
	if err != nil {
		_, ok := err.(error)
		if ok {
			b, marshalErr := json.Marshal(err)
			if marshalErr == nil && string(b) == "{}" {
				err = Wrap(err.(error))
			}
		}
	}

	return &HTTPError{
		StatusCode: statusCode,
		Name:       name,
		Message:    message,
		Handled:    true,
		Err:        err,
	}
}

// JSON returns the JSON representation of the error.
//
// If the Err property is not marshalable, it will
// be omitted from the output.
func (e HTTPError) JSON() string {
	e = anyutils.DestroyCircular(e).(HTTPError)
	b, err := json.Marshal(e)
	if err != nil {
		e.Err = nil
	}

	b, _ = json.Marshal(e)
	return string(b)
}

// NewBadRequestError creates a new HTTPError with a 400 status code.
func NewBadRequestError(message string, err any) error {
	return NewHTTPError(http.StatusBadRequest, "bad_request_error", message, err)
}

// NewNotFoundError creates a new HTTPError with a 404 status code.
func NewNotFoundError(message string, err any) error {
	return NewHTTPError(http.StatusNotFound, "not_found_error", message, err)
}

// NewInternalServerError creates a new HTTPError with a 500 status code.
func NewInternalServerError(message string, err any) error {
	return NewHTTPError(http.StatusInternalServerError, "internal_server_error", message, err)
}

// NewUnauthorizedError creates a new HTTPError with a 401 status code.
func NewUnauthorizedError(message string, err any) error {
	return NewHTTPError(http.StatusUnauthorized, "unauthorized_error", message, err)
}

// NewForbiddenError creates a new HTTPError with a 403 status code.
func NewForbiddenError(message string, err any) error {
	return NewHTTPError(http.StatusForbidden, "forbidden_error", message, err)
}

// NewConflictError creates a new HTTPError with a 409 status code.
func NewConflictError(message string, err any) error {
	return NewHTTPError(http.StatusConflict, "conflict_error", message, err)
}

// NewTooManyRequestsError creates a new HTTPError with a 429 status code.
func NewTooManyRequestsError(message string, err any) error {
	return NewHTTPError(http.StatusTooManyRequests, "too_many_requests_error", message, err)
}

// NewBadGatewayError creates a new HTTPError with a 502 status code.
func NewBadGatewayError(message string, err any) error {
	return NewHTTPError(http.StatusBadGateway, "bad_gateway_error", message, err)
}

// NewServiceUnavailableError creates a new HTTPError with a 503 status code.
func NewServiceUnavailableError(message string, err any) error {
	return NewHTTPError(http.StatusServiceUnavailable, "service_unavailable_error", message, err)
}

// NewGatewayTimeoutError creates a new HTTPError with a 504 status code.
func NewGatewayTimeoutError(message string, err any) error {
	return NewHTTPError(http.StatusGatewayTimeout, "gateway_timeout_error", message, err)
}

// NewHTTPError creates a new HTTPError.
func NewMethodNotAllowedError(message string, err any) error {
	return NewHTTPError(http.StatusMethodNotAllowed, "method_not_allowed_error", message, err)
}
