package goapp

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingolabscl/go-app/errors"
	"github.com/pingolabscl/go-app/logger"
)

type App struct {
	ctx    context.Context
	Logger logger.Logger
	mux    *chi.Mux
}

type Config struct {
	Name     string
	LogLevel logger.Level
}

func New(cfg Config) (app *App, err error) {
	if cfg.Name == "" {
		return nil, errors.New("app name is empty")
	}

	app = &App{
		ctx: context.Background(),
		Logger: logger.New(logger.Config{
			Name:  cfg.Name,
			Level: cfg.LogLevel,
		}),
		mux: chi.NewRouter(),
	}

	app.mux.Use(newAppContextMdw(app))
	app.mux.Use(newRequestIdMdw())
	app.mux.Use(newXPoweredByMdw())
	app.mux.Get("/healthcheck", newHealthcheckHandler())

	return app, nil
}

func (app *App) Start() {
	app.Logger.Info(app.ctx, "app_starting")

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		app.Logger.Fatal(app.ctx, "app_crashed", err)
	}

	app.Logger.Info(app.ctx, "app_started")
	err = http.Serve(listener, app.mux)
	if err != nil {
		app.Logger.Fatal(app.ctx, "app_crashed", err)
	}
}

func (app *App) StartTLS(certFile, keyFile string) {
	app.Logger.Info(app.ctx, "app_starting")

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		app.Logger.Fatal(app.ctx, "app_crashed", err)
	}

	app.Logger.Error(app.ctx, "app_started", err)
	err = http.ServeTLS(listener, app.mux, certFile, keyFile)
	if err != nil {
		app.Logger.Fatal(app.ctx, "app_crashed", err)
	}
}

func (app *App) Get(path string, handler Handler) {
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Get(path, wrappedHandler)
}

func (app *App) Post(path string, handler Handler) {
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Post(path, wrappedHandler)
}

func (app *App) Put(path string, handler Handler) {
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Put(path, wrappedHandler)
}

func (app *App) Delete(path string, handler Handler) {
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Delete(path, wrappedHandler)
}

func (app *App) Patch(path string, handler Handler) {
	wrappedHandler := newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			app:     app,
			Request: r,
			writer:  w,
		})
	})
	app.mux.Patch(path, wrappedHandler)
}
