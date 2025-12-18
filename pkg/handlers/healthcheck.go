package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// NewHealthcheckHandler returns a http.HandlerFunc that
// sends a 200 OK response with a json body containing
// the uptime.
func NewHealthcheckHandler() http.HandlerFunc {
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
