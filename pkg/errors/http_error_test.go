package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHTTPError(t *testing.T) {
	t.Run("should print error", func(t *testing.T) {
		v := &HTTPError{
			StatusCode: 404,
			Name:       "not_found_error",
			Message:    "not found",
			Handled:    true,
			Err:        nil,
		}
		want := "not_found_error: not found"
		if got := v.Error(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should serialize error", func(t *testing.T) {
		v := &HTTPError{
			StatusCode: 404,
			Name:       "not_found_error",
			Message:    "not found",
			Handled:    true,
			Err: map[string]any{
				"key": map[string]any{
					"nested": map[string]any{
						"nestedChild": map[string]any{
							"hello": "world",
						},
					},
				},
			},
		}
		want := `{"statusCode":404,"name":"not_found_error","message":"not found","handled":true,"error":{"key":{"nested":{"nestedChild":{"hello":"world"}}}}}`

		if got := v.JSON(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should serialize even when json.Marshal fails", func(t *testing.T) {
		v := &HTTPError{
			StatusCode: 404,
			Name:       "not_found_error",
			Message:    "not found",
			Handled:    true,
			Err: map[string]any{
				"fn":  func() {},
				"key": "value",
			},
		}
		want := `{"statusCode":404,"name":"not_found_error","message":"not found","handled":true,"error":null}`
		if got := v.JSON(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should wrap built-in error", func(t *testing.T) {
		builtinErr := fmt.Errorf("error")
		err := NewHTTPError(404, "not_found_error", "not found", builtinErr)
		want := `{"statusCode":404,"name":"not_found_error","message":"not found","handled":true,"error":{"name":"error","message":"error","original":{}}}`
		if got := err.(*HTTPError).JSON(); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestNewNotFoundError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 404,
			Name:       "not_found_error",
			Message:    "not found",
			Handled:    true,
			Err:        nil,
		}
		v := NewNotFoundError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewNotFoundError() = %v, want %v", v, want)
		}
	})
}

func TestNewBadRequestError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 400,
			Name:       "bad_request_error",
			Message:    "bad request",
			Handled:    true,
			Err:        nil,
		}
		v := NewBadRequestError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewBadRequestError() = %v, want %v", v, want)
		}
	})
}

func TestNewInternalServerError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 500,
			Name:       "internal_server_error",
			Message:    "internal server error",
			Handled:    true,
			Err:        nil,
		}
		v := NewInternalServerError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewInternalServerError() = %v, want %v", v, want)
		}
	})
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 401,
			Name:       "unauthorized_error",
			Message:    "unauthorized",
			Handled:    true,
			Err:        nil,
		}
		v := NewUnauthorizedError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewUnauthorizedError() = %v, want %v", v, want)
		}
	})
}

func TestNewForbiddenError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 403,
			Name:       "forbidden_error",
			Message:    "forbidden",
			Handled:    true,
			Err:        nil,
		}
		v := NewForbiddenError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewForbiddenError() = %v, want %v", v, want)
		}
	})
}

func TestNewConflictError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 409,
			Name:       "conflict_error",
			Message:    "conflict",
			Handled:    true,
			Err:        nil,
		}
		v := NewConflictError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewConflictError() = %v, want %v", v, want)
		}
	})
}

func TestNewTooManyRequestsError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 429,
			Name:       "too_many_requests_error",
			Message:    "too many requests",
			Handled:    true,
			Err:        nil,
		}
		v := NewTooManyRequestsError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewTooManyRequestsError() = %v, want %v", v, want)
		}
	})
}

func TestNewBadGatewayError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 502,
			Name:       "bad_gateway_error",
			Message:    "bad gateway",
			Handled:    true,
			Err:        nil,
		}
		v := NewBadGatewayError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewBadGatewayError() = %v, want %v", v, want)
		}
	})
}

func TestNewServiceUnavailableError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 503,
			Name:       "service_unavailable_error",
			Message:    "service unavailable",
			Handled:    true,
			Err:        nil,
		}
		v := NewServiceUnavailableError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewServiceUnavailableError() = %v, want %v", v, want)
		}
	})
}

func TestNewGatewayTimeoutError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 504,
			Name:       "gateway_timeout_error",
			Message:    "gateway timeout",
			Handled:    true,
			Err:        nil,
		}
		v := NewGatewayTimeoutError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewGatewayTimeoutError() = %v, want %v", v, want)
		}
	})
}

func TestNewMethodNotAllowedError(t *testing.T) {
	t.Run("should return a new HTTPError", func(t *testing.T) {
		want := &HTTPError{
			StatusCode: 405,
			Name:       "method_not_allowed_error",
			Message:    "method not allowed",
			Handled:    true,
			Err:        nil,
		}
		v := NewMethodNotAllowedError(want.Message, nil)
		if !reflect.DeepEqual(want, v) {
			t.Errorf("NewMethodNotAllowedError() = %v, want %v", v, want)
		}
	})
}
