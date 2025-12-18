package utils

import (
	"encoding/json"
	"net/http"

	"github.com/iolave/go-errors"
)

var _ http.ResponseWriter = &ResponseWriter{}

// customResponseWriter is a wrapper around http.ResponseWriter that
// allows us to store the status code and error that was sent to the
// client, it implements http.ResponseWriter interface.
type ResponseWriter struct {
	SentStatus int
	SentErr    error
	Original   http.ResponseWriter
}

func (w ResponseWriter) Header() http.Header {
	return w.Original.Header()
}
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if w.SentStatus < 200 || w.SentStatus > 299 {
		err := errors.HTTPError{}
		json.Unmarshal(b, &err)
		w.SentErr = &err
	}
	return w.Original.Write(b)
}
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.SentStatus = statusCode
	w.Original.WriteHeader(statusCode)
}
