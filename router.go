package betsi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/iolave/go-betsi/internal/utils"
	"github.com/iolave/go-errors"
	"github.com/iolave/go-trace"
)

// Verify that Router implements the http.Handler interface.
var _ http.Handler = &Router{}

// NewRouter returns a new initialized Router.
func NewRouter() *Router {
	r := &Router{
		mux: chi.NewRouter(),
	}

	return r
}

// Router is a chi.Router wrapper that implements
// the http.Handler interface and is adapted to
// work with the appfactory.Handler.
type Router struct {
	mux *chi.Mux
}

// ServeHTTP is the single method of the http.Handler
// interface that makes Router interoperable with the
// standard library. It uses a sync.Pool to get
// and reuse routing contexts for each request.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Route attaches another appfactory.Router as a
// subrouter along a routing path. It's very useful
// to split up a large API as many independent
// routers and compose them as a single service
// using Route.
func (r Router) Route(pattern string, router *Router) {
	r.mux.Mount(pattern, router.mux)
}

// Use adds middleware to the router.
func (r Router) Use(mdw func(next http.Handler) http.Handler) {
	r.mux.Use(mdw)
}

// With adds inline middleware for an endpoint handler. It uses a
// standard [http.Hanlder] mdw.
func (r Router) With(mdw func(next http.Handler) http.Handler) Router {
	return Router{
		mux: r.mux.With(mdw).(*chi.Mux),
	}
}

// NotFoundHandler sets the handler to be called when a route is not found.
func (r Router) NotFoundHandler(h http.Handler) {
	r.mux.NotFound(h.(http.HandlerFunc))
}

// MethodNotAllowedHandler sets the handler to be called when a request
// comes with an invalid method.
func (r Router) MethodNotAllowedHandler(h http.Handler) {
	r.mux.MethodNotAllowed(h.(http.HandlerFunc))
}

// buildPatterns builds the patterns for the router by generating
// a `/` sufixed and non-sufixed pattern.
//
//   - If the pattern is empty it panics.
//   - If the pattern is `/` it returns a slice with a single `/`.
func (r Router) buildPatterns(pattern string) []string {
	patterns := []string{}

	if pattern == "" {
		panic("pattern cannot be empty")
	}
	pat := strings.TrimSuffix(pattern, "/")
	if pat == "" {
		patterns = append(patterns, "/")
	} else {
		patterns = append(patterns, fmt.Sprintf("%s/", pat))
		patterns = append(patterns, pat)
	}
	return patterns
}

// convertToAcceptedHandler converts the appfactory.Handler
// to http.HandlerFunc.
func (r Router) convertToAcceptedHandlerWithLog(
	handler Handler[any, any],
) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest[any, any]{
			Req: r,
			w:   w,
		})
	})

	wrappedFn := func(w http.ResponseWriter, r *http.Request) {
		// app, _ := GetFromContext(r.Context())
		// if app == nil {
		fn(w, r)
		// } else {
		// 	//app.mdwAppRequestLog(fn).ServeHTTP(w, r)
		// }
	}

	return wrappedFn
}

// Get adds the route `pattern` that matches a GET http
// method to execute the `handler` appfactory.Handler.
func (r Router) Get(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Get(pat, h)
	}
}

// Post adds the route `pattern` that matches a POST http
// method to execute the `handler` appfactory.Handler.
func (r Router) Post(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Post(pat, h)
	}
}

// Put adds the route `pattern` that matches a PUT http
// method to execute the `handler` appfactory.Handler.
func (r Router) Put(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Put(pat, h)
	}
}

// Delete adds the route `pattern` that matches a DELETE http
// method to execute the `handler` appfactory.Handler.
func (r Router) Delete(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Delete(pat, h)
	}
}

// Path adds the route `pattern` that matches a PATCH http
// method to execute the `handler` appfactory.Handler.
func (r Router) Patch(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Patch(pat, h)
	}
}

// Head adds the route `pattern` that matches a HEAD http
// method to execute the `handler` appfactory.Handler.
func (r Router) Head(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Head(pat, h)
	}
}

// Options adds the route `pattern` that matches an OPTIONS
// http method to execute the `handler` appfactory.Handler.
func (r Router) Options(pattern string, handler Handler[any, any]) {
	h := r.convertToAcceptedHandlerWithLog(handler)
	for _, pat := range r.buildPatterns(pattern) {
		r.mux.Options(pat, h)
	}
}

// Handler is the type of the handler function that
// is used to handle requests within an betsi app.
type Handler[In, Out any] func(ar AppRequest[In, Out])

// ServeHTTP is the single method of the http.Handler
// interface that makes Router interoperable with the
// standard library. It uses a sync.Pool to get
// and reuse routing contexts for each request.
func (h Handler[In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(AppRequest[In, Out]{
		Req: r,
		w:   w,
	})
}

// NewHandler wraps a type-safe Handler[In, Out] into a generic Handler[any, any].
//
// This function is useful for adapting specific, strongly-typed handlers to contexts
// where the router expects a more general handler signature (e.g., when registering
// routes where the exact input/output types are not yet known or are managed
// internally by the AppRequest's mechanisms).
//
// The returned handler internally converts the generic AppRequest[any, any] back
// into the original AppRequest[In, Out] before invoking the provided type-safe handler `h`.
func NewHandler[In, Out any](h Handler[In, Out]) Handler[any, any] {
	nh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ar := AppRequest[In, Out]{
			Req: r,
			w:   w,
		}
		h(ar)
	})

	return func(ar AppRequest[any, any]) {
		nh.ServeHTTP(ar.w, ar.Req)
	}
}

// AppRequest is a generic wrapper around the standard http.Request and
// http.ResponseWriter. It provides type-safe methods for parsing incoming
// requests and sending outgoing responses.
//
// The type parameters `In` and `Out` define the expected request and response
// body types, respectively.
//
// # Request Parsing (In)
//
// The `In` type parameter specifies the shape of the data to be parsed from an
// incoming request, typically by the ParseRequest method. It must be a struct
// that uses `ar` tags to map request data:
//
//   - `ar:"path={name}"`: Maps a URL path parameter to a string field.
//     For example, in `/users/{id}`, a field with `ar:"path=id"` would be populated
//     with the value of {id}.
//   - `ar:"body=json"`: Designates a struct field as the target for the
//     deserialized JSON request body.
//
// Example of an `In` type:
//
//	type CreateUserRequest struct {
//		ID    string `ar:"path=id"`
//		// The user's details, deserialized from the JSON body.
//		Body struct {
//			Name  string `json:"name"`
//			Email string `json:"email"`
//		} `ar:"body=json"`
//	}
//
// # Response Sending (Out)
//
// The `Out` type parameter defines the data structure that will be sent as a
// response body, typically using the SendJSON method. This can be any
// serializable type. For JSON responses, standard `json` tags are used.
//
// Example of an `Out` type:
//
//	type UserResponse struct {
//		ID        string    `json:"id"`
//		Name      string    `json:"name"`
//		CreatedAt time.Time `json:"createdAt"`
//	}
type AppRequest[In, Out any] struct {
	Req *http.Request
	w   http.ResponseWriter
}

// Context returns the request's context, serving as a shortcut for `ar.Req.Context()`.
//
// As a safeguard, if the underlying request (`ar.Req`) is nil, it returns a
// non-nil, empty context via `context.Background()` to prevent panics.
func (ar AppRequest[_, _]) Context() context.Context {
	if ar.Req == nil {
		return context.Background()
	}

	return ar.Req.Context()
}

// SendJSONError sends a structured JSON error response to the client.
//
// It intelligently handles the provided error:
//   - If `err` is of type [github.com/iolave/go-errors.HTTPError], it
//     uses its status code and
//     JSON body directly for the response.
//   - If `err` is nil, it generates a new internal server error.
//   - For any other error type, it wraps the original error in a new
//     internal server error.
//
// The context `ctx` is used to retrieve trace information, which is injected
// into the response headers for observability. The Content-Type is always
// set to "application/json".
func (ar AppRequest[_, _]) SendJSONError(ctx context.Context, err error) {
	t := trace.GetFromContext(ctx)

	switch err.(type) {
	case *errors.HTTPError:
		break
	case nil:
		err = errors.NewInternalServerError(
			ERR_SRV_AR_NIL_ERR,
			err,
		)
	default:
		err = errors.NewInternalServerError(
			ERR_SRV_AR_GENERIC_ERR,
			err,
		)
	}

	httperr := err.(*errors.HTTPError)
	t.SetHTTPHeaders(ar.w.Header())
	ar.w.Header().Set("Content-Type", "application/json")
	ar.w.WriteHeader(httperr.StatusCode)
	ar.w.Write(httperr.JSON())
}

// SendJSON marshals the provided data `v` into a JSON response, sets the
// Content-Type header to "application/json", and writes an http.StatusOK (200)
// status code.
//
// The context `ctx` is used to retrieve trace information which is then injected
// into the response headers for observability.
//
// Before sending, it recursively validates the payload `v` using any associated
// go-playground/validator tags.
//
// It handles two primary error scenarios:
//   - If validation fails, it calls SendJSONError with an internal server error.
//   - If JSON marshaling fails, it also calls SendJSONError with an internal
//     server error.
func (ar AppRequest[_, Out]) SendJSON(ctx context.Context, v Out) {
	if err := utils.ValidateRecursively(v); err != nil {
		ar.SendJSONError(ctx, errors.NewInternalServerError(
			ERR_SRV_AR_SEND_JSON_VALIDATION_ERR,
			err,
		))
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ar.SendJSONError(ctx, errors.NewInternalServerError(
			ERR_SRV_AR_SEND_JSON_MARSHALL_ERR,
			err,
		))
		return
	}

	t := trace.GetFromContext(ctx)
	t.SetHTTPHeaders(ar.w.Header())

	ar.w.Header().Set("Content-Type", "application/json")
	ar.w.WriteHeader(http.StatusOK)
	ar.w.Write(b)
}

// ParseRequest populates and returns a new instance of the generic type `In` by
// decoding data from the HTTP request. It uses the `ar` struct tags on the `In`
// type to map URL path parameters and the request body to the struct's fields.
//
// This method returns an error under the following conditions:
//   - The generic type `In` is not a struct.
//   - The underlying http.Request or its body is nil.
//   - The internal decoding of the request fails (e.g., malformed JSON).
//
// On success, it returns a pointer to the populated struct. Any returned error
// will be of type [github.com/iolave/go-errors.GenericError].
func (ar AppRequest[In, _]) ParseRequest() (*In, error) {
	if reflect.TypeFor[In]().Kind() != reflect.Struct {
		return nil, errors.NewWithName(
			ERR_NAME_PARSE,
			fmt.Sprintf(ERR_INVALID_TYPE, "struct", reflect.TypeFor[In]().Kind().String()),
		)
	}

	if ar.Req == nil {
		return nil, errors.NewWithName(
			ERR_NAME_PARSE,
			ERR_AR_NIL_REQ,
		)
	}

	if ar.Req.Body == nil {
		return nil, errors.NewWithName(
			ERR_NAME_PARSE,
			ERR_AR_NIL_REQ_BODY,
		)
	}

	v := new(In)
	err := decodeAppRequest(ar.Req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
