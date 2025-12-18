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
func (r Router) Route(pattern string, router Router) {
	r.mux.Mount(pattern, router.mux)
}

// Use adds middleware to the router.
func (r Router) Use(mdw func(next http.Handler) http.Handler) {
	r.mux.Use(mdw)
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

// NewHandler returns a new generic handler with it's
// input and output types set to any and any respectively
// from a type-safe handler.
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

// AppRequest is the type of the request that is
// passed to the handler function and is used to
// retrieve request data, such as the request body
// and headers. It is also used to send a response
// back to the client.
//
// It contains various useful methods to be used
// inside a handler function.
type AppRequest[In, Out any] struct {
	Req *http.Request
	w   http.ResponseWriter
}

// Context returns the context of the request.
// It is a shortcut for ar.Req.Context().
//
// If the request is nil, a new context will be created.
func (ar AppRequest[_, _]) Context() context.Context {
	if ar.Req == nil {
		return context.Background()
	}

	return ar.Req.Context()
}

// SendError sends a json error response to the client. The error
// will be of type *errors.HTTPError. If the error is nil or it
// is not of type *errors.HTTPError an internal server error will
// be sent to the client.
func (ar AppRequest[_, _]) SendError(ctx context.Context, err error) {
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

// SendJSON sends a 200 json response to the client with type Out.
//
//   - If v contains valid go-playground/validator tags, v will be
//     validated before sending it to the client to ensure the client
//     receives a proper response. If v failed to be validated, an
//     internal server error will be sent.
//   - If v can not be json marshaled, an internal server error will
//     be sent to the client.
func (ar AppRequest[_, Out]) SendJSON(ctx context.Context, v Out) {
	if err := utils.ValidateRecursively(v); err != nil {
		ar.SendError(ctx, errors.NewInternalServerError(
			ERR_SRV_AR_SEND_JSON_VALIDATION_ERR,
			err,
		))
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ar.SendError(ctx, errors.NewInternalServerError(
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

// ParseRequest parses the [http.Request] and returns a new instance of In
// with it's properties set to the values of the request (path params and body).
//
// Any error returned by this method will be of type [github.com/iolave/go-errors.GenericError].
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
