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
	"github.com/pingolabscl/go-app/pkg/errors"
	"github.com/pingolabscl/go-app/pkg/trace"
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
}

// RequestJSON sends an http request to a goapp client
// and returns a json response.
//
// It currently supports sending json requests.
// TODO: Add support for streaming requests.
//
// It retrieves the goapp from the passed context and
// uses the client name to retrieve the client from the
// goapp clients map. If the client is not found, it
// returns an error.
//
// If the context does not contain a goapp, it returns
// an internal server error and might indicate the context
// is not set correctly.
func RequestJSON[Response any](ctx context.Context, r JSONRequest) (Response, *errors.HTTPError) {
	// gets the trace from the context and sets the request id
	// if it's empty to trace following requests.
	tr := trace.GetFromContext(ctx)
	if tr.RequestID == "" {
		tr.RequestID = uuid.NewString()
	}
	ctx = trace.SetContext(ctx, tr)

	result := new(Response)

	app, err := GetFromContext(ctx)
	if err != nil {
		return *result, errors.NewInternalServerError("unable to send request", err)
	}

	client, ok := app.clients[r.ClientName]
	if !ok {
		return *result, errors.NewInternalServerError(
			"unable to send request, client not found",
			nil,
		)
	}

	var url string
	if client.TLSEnabled {
		url = fmt.Sprintf("https://%s:%d%s", client.Host, client.Port, r.Path)
	} else {
		url = fmt.Sprintf("http://%s:%d%s", client.Host, client.Port, r.Path)
	}

	b, error := json.Marshal(r.Body)
	if error != nil {
		return *result, errors.NewInternalServerError(
			"unable to send request, failed to marshal body",
			errors.Wrap(error),
		)
	}

	req, error := http.NewRequestWithContext(ctx, r.Method, url, bytes.NewReader(b))
	if error != nil {
		return *result, errors.NewInternalServerError(
			"unable to send request, failed to create request",
			errors.Wrap(error),
		)
	}

	req.Header.Set("Content-Type", "application/json")
	tr.SetHTTPHeaders(req.Header)

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	res, error := app.httpClient.Do(req)
	if error != nil {
		return *result, errors.NewInternalServerError(
			"unable to send request",
			errors.Wrap(error),
		)
	}

	if res.Header.Get("X-Powered-By") != "pingolabs.cl" {
		return *result, errors.NewInternalServerError(
			"unable to process response, client is not a pingolabs service",
			nil,
		)
	}

	b, error = io.ReadAll(res.Body)
	if error != nil {
		return *result, errors.NewInternalServerError(
			"failed to read response body",
			errors.Wrap(error),
		)
	}

	if res.StatusCode != http.StatusOK {
		httpErr := new(errors.HTTPError)
		if err := json.Unmarshal(b, httpErr); err != nil {
			return *result, errors.NewInternalServerError(
				"failed to unmarshal error response",
				errors.Wrap(err),
			)
		}

		return *result, httpErr
	}

	if err := json.Unmarshal(b, result); err != nil {
		return *result, errors.NewInternalServerError(
			"failed to unmarshal response",
			errors.Wrap(err),
		)
	}

	return *result, nil
}
