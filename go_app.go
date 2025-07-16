package goapp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/vault-client-go"
	"github.com/pingolabscl/go-app/pkg/errors"
	"github.com/pingolabscl/go-app/pkg/logger"
)

type App struct {
	ctx         context.Context
	Logger      *logger.Logger
	mux         *chi.Mux
	port        int
	clients     map[string]client
	httpClient  *http.Client
	validator   *validator.Validate
	vault       *vault.Client
	vaultConfig VaultConfig
}

type Config struct {
	Name               string
	LogLevel           logger.Level
	Port               int
	InsecureSkipVerify bool
	Vault              VaultConfig
}

// New creates a new App. It returns an error
// of type *errors.Error if the app cannot be created.
func New(cfg Config) (*App, error) {
	ctx := context.Background()

	vaultConfig := determineVaultConfig(cfg.Vault)
	var vaultClient *vault.Client
	if vaultConfig.Addr != "" {
		vault, err := vault.New(
			vault.WithAddress(vaultConfig.Addr),
			vault.WithRequestTimeout(30*time.Second),
		)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		customHeaders := http.Header{}
		customHeaders.Set("CF-Access-Client-Id", os.Getenv("CF_ACCESS_CLIENT_ID"))
		customHeaders.Set("CF-Access-Client-Secret", os.Getenv("CF_ACCESS_CLIENT_SECRET"))
		if err := vault.SetCustomHeaders(customHeaders); err != nil {
			return nil, errors.Wrap(err)
		}
		vaultClient = vault
	}

	if cfg.Name == "" {
		return nil, errors.New("app name is empty")
	}

	port := 3000
	if cfg.Port != 0 {
		port = cfg.Port
	}

	logger, err := logger.New(logger.Config{
		Name:  cfg.Name,
		Level: cfg.LogLevel,
	})
	if err != nil {
		return nil, err
	}

	app := &App{
		ctx:         ctx,
		port:        port,
		Logger:      logger,
		mux:         chi.NewRouter(),
		clients:     make(map[string]client),
		httpClient:  newHTTPClient(cfg.InsecureSkipVerify),
		validator:   validator.New(validator.WithRequiredStructEnabled()),
		vault:       vaultClient,
		vaultConfig: vaultConfig,
	}

	app.mux.Use(newAppContextMdw(app))
	app.mux.Use(newTraceMdw())
	app.mux.Use(newXPoweredByMdw())
	app.mux.Get("/healthcheck", newHealthcheckHandler())
	app.mux.Get("/healthcheck/", newHealthcheckHandler())
	app.mux.MethodNotAllowed(newMethodNotAllowedHandler(app))
	app.mux.NotFound(newNotFoundHandler(app))

	return app, nil
}

func (app *App) Start() {
	app.Logger.InfoWithData(app.ctx, "app_starting", map[string]any{
		"port":       app.port,
		"tlsEnabled": false,
	})

	if app.vault != nil {
		if err := app.vault.SetToken(app.vaultConfig.Token); err != nil {
			app.Logger.FatalWithData(app.ctx, "app_crashed", err, map[string]any{
				"port":       app.port,
				"tlsEnabled": false,
			})
		}
		go app.renewVaultToken()
	}

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
	wrappedHandler := app.newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			App:     app,
			Request: r,
			writer:  w,
		})
	})
	if path == "" {
		app.mux.Get("/", wrappedHandler)
		return
	}
	app.mux.Get(path, wrappedHandler)
	app.mux.Get(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Post(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := app.newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			App:     app,
			Request: r,
			writer:  w,
		})
	})
	if path == "" {
		app.mux.Post("/", wrappedHandler)
		return
	}
	app.mux.Post(path, wrappedHandler)
	app.mux.Post(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Put(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := app.newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			App:     app,
			Request: r,
			writer:  w,
		})
	})
	if path == "" {
		app.mux.Put("/", wrappedHandler)
		return
	}
	app.mux.Put(path, wrappedHandler)
	app.mux.Put(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Delete(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := app.newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			App:     app,
			Request: r,
			writer:  w,
		})
	})
	if path == "" {
		app.mux.Delete("/", wrappedHandler)
		return
	}
	app.mux.Delete(path, wrappedHandler)
	app.mux.Delete(fmt.Sprintf("%s/", path), wrappedHandler)
}

func (app *App) Patch(path string, handler Handler) {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	wrappedHandler := app.newHandler(func(w http.ResponseWriter, r *http.Request) {
		handler(AppRequest{
			App:     app,
			Request: r,
			writer:  w,
		})
	})
	if path == "" {
		app.mux.Patch("/", wrappedHandler)
		return
	}
	app.mux.Patch(path, wrappedHandler)
	app.mux.Patch(fmt.Sprintf("%s/", path), wrappedHandler)
}
