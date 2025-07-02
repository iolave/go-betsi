package goapp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/pingolabscl/go-app/errors"
	"github.com/pingolabscl/go-app/trace"
)

func newHTTPClient(insecureSkipVerify bool) *http.Client {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	return &http.Client{
		Transport: &transport,
	}
}

type client struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	TLSEnabled bool   `json:"tlsEnabled"`
}

func (app *App) RegisterClient(name, host string, port int, tlsEnabled bool) {
	client := client{
		Name:       name,
		Host:       host,
		Port:       port,
		TLSEnabled: tlsEnabled,
	}
	app.clients[name] = client
	app.Logger.DebugWithData(app.ctx, "register_client", map[string]any{
		"client": client,
	})
}

type JSONRequest struct {
	ClientName string            `json:"clientName"`
	Path       string            `json:"path"`
	Method     string            `json:"method"`
	Body       any               `json:"body"`
	Headers    map[string]string `json:"headers"`
	// TODO: Add query params support
	// Query      map[string]string `json:"query"`
	Result any `json:"-"`
}

// RequestJSON sends a request to the specified client and returns the response, it only supports JSON requests at the moment.
// TODO: Add support for streaming requests.
func (app *App) RequestJSON(ctx context.Context, r JSONRequest) error {
	client, ok := app.clients[r.ClientName]
	if !ok {
		return errors.NewInternalServerError("unable to send request", "client not found")

	}

	var url string
	if client.TLSEnabled {
		url = fmt.Sprintf("https://%s:%d%s", client.Host, client.Port, r.Path)
	} else {
		url = fmt.Sprintf("http://%s:%d%s", client.Host, client.Port, r.Path)
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return errors.NewBadRequestError("unable to send request", "failed to marshal body")
	}

	tr := trace.GetFromContext(ctx)
	if tr.RequestID == "" {
		tr.RequestID = uuid.NewString()
	}
	ctx = trace.SetContext(ctx, tr)

	req, err := http.NewRequestWithContext(ctx, r.Method, url, bytes.NewReader(b))
	if err != nil {
		return errors.NewInternalServerError("unable to send request", "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", tr.RequestID)

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	res, err := app.httpClient.Do(req)
	if err != nil {
		return errors.NewInternalServerError("failed to send request", err.Error())
	}

	if res.Header.Get("X-Powered-By") != "pingolabs.cl" {
		return errors.NewInternalServerError("failed to read response", "client is not a pingolabs service")
	}

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return errors.NewInternalServerError("failed to read response body", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		httpErr := new(errors.HTTPError)
		if err := json.Unmarshal(b, httpErr); err != nil {
			return errors.NewInternalServerError("failed to unmarshal error response", err.Error())
		}
		return httpErr
	}

	return json.Unmarshal(b, r.Result)
}
