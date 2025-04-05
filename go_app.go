package goapp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingolabscl/go-app/errors"
	"github.com/pingolabscl/go-app/logger"
)

type App struct {
	ctx        context.Context
	Logger     logger.Logger
	mux        *chi.Mux
	port       int
	clients    map[string]client
	httpClient *http.Client
}

type Config struct {
	Name               string
	LogLevel           logger.Level
	Port               int
	InsecureSkipVerify bool
}

func New(cfg Config) (app *App, err error) {
	if cfg.Name == "" {
		return nil, errors.New("app name is empty")
	}

	port := 3000
	if cfg.Port != 0 {
		port = cfg.Port
	}

	app = &App{
		port: port,
		ctx:  context.Background(),
		Logger: logger.New(logger.Config{
			Name:  cfg.Name,
			Level: cfg.LogLevel,
		}),
		mux:        chi.NewRouter(),
		clients:    make(map[string]client),
		httpClient: newHTTPClient(cfg.InsecureSkipVerify),
	}

	app.mux.Use(newAppContextMdw(app))
	app.mux.Use(newRequestIdMdw())
	app.mux.Use(newXPoweredByMdw())
	app.mux.Get("/healthcheck", newHealthcheckHandler())
	app.mux.Get("/healthcheck/", newHealthcheckHandler())
	app.mux.NotFound(newNotFoundHandler(app))

	return app, nil
}

func (app *App) Start() {
	app.Logger.InfoWithData(app.ctx, "app_starting", map[string]any{
		"port":       app.port,
		"tlsEnabled": false,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		app.Logger.FatalWithData(app.ctx, "app_crashed", err, map[string]any{
			"port":       app.port,
			"tlsEnabled": false,
		})
	}

	app.Logger.InfoWithData(app.ctx, "app_started", map[string]any{
		"port": app.port,
	})
	err = http.Serve(listener, app.mux)
	if err != nil {
		app.Logger.FatalWithData(app.ctx, "app_crashed", err, map[string]any{
			"port":       app.port,
			"tlsEnabled": false,
		})
	}
}

func (app *App) StartTLS(certFile, keyFile string) {
	app.Logger.InfoWithData(app.ctx, "app_starting", map[string]any{
		"port":       app.port,
		"tlsEnabled": true,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		app.Logger.FatalWithData(app.ctx, "app_crashed", err, map[string]any{
			"port":       app.port,
			"tlsEnabled": true,
		})
	}

	app.Logger.InfoWithData(app.ctx, "app_started", map[string]any{
		"port": app.port,
	})
	err = http.ServeTLS(listener, app.mux, certFile, keyFile)
	if err != nil {
		app.Logger.FatalWithData(app.ctx, "app_crashed", err, map[string]any{
			"port":       app.port,
			"tlsEnabled": true,
		})
	}
}

func (app *App) Get(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Get(path, wrappedHandler)
	app.mux.Get(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Post(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Post(path, wrappedHandler)
	app.mux.Post(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Put(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Put(path, wrappedHandler)
	app.mux.Put(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Delete(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Delete(path, wrappedHandler)
	app.mux.Delete(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Patch(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Patch(path, wrappedHandler)
	app.mux.Patch(fmt.Sprintf("%s/", path), wrappedHandler)
}
