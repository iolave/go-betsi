package middlewares

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/iolave/go-errors"
	"github.com/iolave/go-logger"
)

// RateLimit is the rate limit for an IP address.
type RateLimit struct {
	// IP is the IP address of the request.
	IP string `json:"ip"`

	// LastRequestedAt is the time the request was made.
	LastRequestedAt time.Time `json:"lastRequestedAt"`

	// Metric is the time metric for the rate limit.
	// For example, if Metric is set to 1 minute, then a request
	// will be rate limited for 1 minute.
	Metric time.Duration `json:"timeMetric"`

	// Limit is the maximum number of requests allowed within the
	// time metric.
	Limit int `json:"limit"`

	// Remaining is the number of requests remaining within the
	// time metric.
	Remaining int `json:"remaining"`
}

// Decrement decrements the remaining requests by 1.
func (r *RateLimit) Decrement() {
	r.Remaining--
}

// CanReset returns true if the rate limit can be reset.
func (r *RateLimit) CanReset() bool {
	if r.Remaining > 0 {
		return false
	}

	return time.Now().After(r.LastRequestedAt.Add(r.Metric))
}

// Reset resets the rate limit.
func (r *RateLimit) Reset() {
	r.LastRequestedAt = time.Now()
	r.Remaining = r.Limit
}

// IsLimited returns true if the rate limit is limited.
func (r RateLimit) IsLimited() bool {
	return r.Remaining <= 0
}

// SetResponseHeaders sets the response headers for the rate limit.
// It sets the X-Rate-Limit-Limit, X-Rate-Limit-Remaining, and
// X-Rate-Limit-Reset headers.
func (r RateLimit) SetResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Rate-Limit-Limit", strconv.Itoa(r.Limit))
	w.Header().Set("X-Rate-Limit-Remaining", strconv.Itoa(r.Remaining))
	w.Header().Set("X-Rate-Limit-Reset", strconv.FormatInt(r.LastRequestedAt.Add(r.Metric).Unix(), 10))
}

// RateLimitStore is the interface for a store to use for rate limiting.
type RateLimitStore interface {
	GetLimit(ip string) (*RateLimit, error)
	UpsertLimit(limit RateLimit) error
}

// RateLimitConfig is the configuration for the rate limiting middleware.
type RateLimitConfig struct {
	// Store is the store to use for rate limiting.
	Store RateLimitStore

	// Metric is the time metric to use for rate limiting.
	// For example, if Metric is set to 1 minute, then a request
	// will be rate limited for 1 minute.
	Metric time.Duration

	// Limit is the maximum number of requests allowed within the
	// time metric.
	Limit int

	// Logger is an optional logger to use for logging errors and
	// rate limit exceeded errors.
	Logger logger.Logger
}

func NewRateLimitMdwWithJSONError(config RateLimitConfig) func(next http.Handler) http.Handler {
	if config.Store == nil {
		panic("store is required")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := getIPFromRequest(r)
			if err != nil {
				httpErr := errors.NewInternalServerError(
					"failed to get determine rate limit for incoming request",
					err,
				).(*errors.HTTPError)
				w.WriteHeader(httpErr.StatusCode)
				w.Write([]byte(httpErr.JSON()))
				if config.Logger != nil {
					config.Logger.Error(
						r.Context(),
						"middleware_rate_limiting_failed",
						err,
					)
				}
				return
			}

			rl, err := config.Store.GetLimit(ip)
			if err != nil {
				httpErr := errors.NewInternalServerError(
					"failed to get rate limit",
					err,
				).(*errors.HTTPError)

				w.WriteHeader(httpErr.StatusCode)
				w.Write([]byte(httpErr.JSON()))
				if config.Logger != nil {
					config.Logger.Error(
						r.Context(),
						"middleware_rate_limiting_failed",
						err,
					)
				}
				return
			}
			if rl == nil {
				rl = &RateLimit{
					IP:              ip,
					LastRequestedAt: time.Now(),
					Metric:          config.Metric,
					Limit:           config.Limit,
					Remaining:       config.Limit,
				}
				err := config.Store.UpsertLimit(*rl)
				if err != nil {
					httpErr := errors.NewInternalServerError(
						"failed to update rate limit",
						err,
					).(*errors.HTTPError)

					w.WriteHeader(httpErr.StatusCode)
					w.Write([]byte(httpErr.JSON()))
					if config.Logger != nil {
						config.Logger.Error(
							r.Context(),
							"middleware_rate_limiting_failed",
							err,
						)
					}
					return
				}
			}
			if rl.IsLimited() {
				if !rl.CanReset() {
					err := errors.NewTooManyRequestsError(
						"rate limit exceeded",
						nil,
					).(*errors.HTTPError)
					rl.SetResponseHeaders(w)
					w.WriteHeader(err.StatusCode)
					w.Write(err.JSON())
					if config.Logger != nil {
						config.Logger.ErrorWithData(r.Context(), "rate_limit_exceeded_error", err, map[string]any{"ip": ip})
					}
					return
				}

				rl.Reset()
			}

			rl.Decrement()
			config.Store.UpsertLimit(*rl)
			rl.SetResponseHeaders(w)
			next.ServeHTTP(w, r)
		})
	}
}

// getIPFromRequest returns the IP address of the request.
func getIPFromRequest(r *http.Request) (string, error) {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip, nil
	}

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip, nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return ip, nil
}
