package trace

import "net/http"

const (
	header_trace_request_id = "X-Trace-Request-ID"
)

// SetHTTPHeaders sets trace attributes
// to http headers.
func (t Trace) SetHTTPHeaders(h http.Header) {
	h.Set(header_trace_request_id, t.RequestID)
}

// NewFromHTTPHeaders returns a new Trace from the given http headers.
func NewFromHTTPHeaders(h http.Header) Trace {
	if h == nil {
		return Trace{}
	}

	return Trace{
		RequestID: h.Get(header_trace_request_id),
	}
}
