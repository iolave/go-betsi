package betsi

import (
	"context"
	"net/http"
	"reflect"

	"github.com/iolave/go-errors"
	"github.com/iolave/go-trace"
)

// NewRequest returns a new http request. It takes a url param that can contain
// path params (i.e "/users/{id}") and a v param that has to be pointer to
// a struct or a struct with "ar" tags:
//
//   - path={VALUE}: path param key within the url.
//   - body={oneof json}: property that contains the body type (json or xml).
//
// v example:
//
//	type PatchUserRequest struct {
//		Name string `ar:"path=name"`
//		Body struct{
//			Age int `json:"age"`
//			// ...
//		} `ar:"body=json"`
//	}
//
// Any returned error is of type [github.com/iolave/go-errors.HTTPError].
func NewRequest(ctx context.Context, method string, url string, v any) (*http.Request, error) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
			return nil, errors.NewWithName(
				"app_error",
				"v has to be a pointer to a struct or a struct",
			)
		}
	case reflect.Struct:
		v = reflect.ValueOf(v).Interface()
	default:
		return nil, errors.NewWithName(
			"app_error",
			"v has to be a pointer to a struct or a struct",
		)

	}

	url, reader, err := encodeAppRequest(url, v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	t := trace.GetFromContext(ctx)
	t.SetHTTPHeaders(req.Header)
	return req, nil
}
