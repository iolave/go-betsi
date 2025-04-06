package errors

import (
	"fmt"
	"net/http"
)

type HTTP struct {
	Status  int    `json:"status"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

func (err *HTTP) Error() string {
	return fmt.Sprintf("[%d] %s - %s: %s", err.Status, err.Name, err.Message, err.Cause)
}

func NewBadRequestError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusBadRequest,
		Name:    "bad_request_error",
		Message: msg,
		Cause:   cause,
	}
}

func NewUnauthorizedError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusUnauthorized,
		Name:    "unauthorized_error",
		Message: msg,
		Cause:   cause,
	}
}

func NewNotFoundError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusNotFound,
		Name:    "not_found_error",
		Message: msg,
		Cause:   cause,
	}
}

func NewMethodNotAllowedError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusMethodNotAllowed,
		Name:    "method_not_allowed_error",
		Message: msg,
		Cause:   cause,
	}
}

func NewTimeoutError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusGatewayTimeout,
		Name:    "gateway_timeout_error",
		Message: msg,
		Cause:   cause,
	}
}

func NewInternalServerError(msg, cause string) *HTTP {
	return &HTTP{
		Status:  http.StatusInternalServerError,
		Name:    "internal_server_error",
		Message: msg,
		Cause:   cause,
	}
}
