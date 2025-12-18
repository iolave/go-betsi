package betsi

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iolave/go-errors"
	"github.com/iolave/go-logger"
)

type App struct {
	cfg       Config
	validator *validator.Validate
}

// ServerConfig is the configuration for the http
// server.
type ServerConfig struct {
	// Port is the port for the http server.
	Port int
	// Router is the router used to handle the requests.
	Router *Router
}

type Config struct {
	// Logger is the logger used to log messages.
	Logger logger.Logger

	// Server is the configuration for the http server.
	// If nil, the http server will not be started.
	Server *ServerConfig
}

// New returns a new app. It errors if the logger is nil.
//
// Any returned error is an implementation of the [github.com/iolave/go-errors.Error]
// interface, in specific, of type [github.com/iolave/go-errors.GenericError]
//
// These errors can be casted to the original error by doing:
//
//	gerr := err.(*errors.GenericError)
//	// or
//	gerr := err.(errors.Error)
func New(cfg Config) (*App, error) {
	app := &App{
		cfg:       cfg,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	if cfg.Logger == nil {
		return nil, errors.NewWithName(ERR_NAME, ERR_NIL_LOGGER)
	}

	return app, nil
}

// Start starts the app. It takes a list of functions that will be executed
// before starting the app. The functions will be executed in the order they
// are passed to the function.
func (app *App) Start(preExecution ...func()) {
	for _, f := range preExecution {
		f()
	}

	ctx := context.Background()
	tlsEnabled := false

	if app.cfg.Server == nil {
		err := errors.NewWithName(ERR_NAME, ERR_START_W_NIL_SERVER).(errors.Error)
		panic(err.JSON())
	}

	if app.cfg.Server.Router == nil {
		err := errors.NewWithName(ERR_NAME, ERR_START_W_NIL_ROUTER).(errors.Error)
		panic(err.JSON())
	}

	logger := app.cfg.Logger
	if logger == nil {
		err := errors.NewWithName(ERR_NAME, ERR_NIL_LOGGER).(errors.Error)
		panic(err.JSON())
	}

	logger.InfoWithData(ctx, "app_starting", map[string]any{
		"port":       app.cfg.Server.Port,
		"tlsEnabled": tlsEnabled,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.cfg.Server.Port))
	if err != nil {
		logger.FatalWithData(ctx, "app_crashed", err, map[string]any{
			"port":       app.cfg.Server.Port,
			"tlsEnabled": tlsEnabled,
		})
	}

	logger.InfoWithData(ctx, "app_started", map[string]any{
		"port":       app.cfg.Server.Port,
		"tlsEnabled": tlsEnabled,
	})
	err = http.Serve(listener, app.cfg.Server.Router)
	if err != nil {
		logger.FatalWithData(ctx, "app_crashed", err, map[string]any{
			"port":       app.cfg.Server.Port,
			"tlsEnabled": tlsEnabled,
		})
	}
}

// StartTLS starts the app. It takes a list of functions that will be executed
// before starting the app. The functions will be executed in the order they
// are passed to the function.
func (app *App) StartTLS(certFile, keyFile string, preExecution ...func()) {
	for _, f := range preExecution {
		f()
	}

	ctx := context.Background()
	tlsEnabled := true

	if app.cfg.Server == nil {
		err := errors.NewWithName(ERR_NAME, ERR_START_W_NIL_SERVER).(errors.Error)
		panic(err.JSON())
	}

	if app.cfg.Server.Router == nil {
		err := errors.NewWithName(ERR_NAME, ERR_START_W_NIL_ROUTER).(errors.Error)
		panic(err.JSON())
	}

	logger := app.cfg.Logger
	if logger == nil {
		err := errors.NewWithName(ERR_NAME, ERR_NIL_LOGGER).(errors.Error)
		panic(err.JSON())
	}

	logger.InfoWithData(ctx, "app_starting", map[string]any{
		"port":       app.cfg.Server.Port,
		"tlsEnabled": tlsEnabled,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.cfg.Server.Port))
	if err != nil {
		logger.FatalWithData(ctx, "app_crashed", err, map[string]any{
			"port":       app.cfg.Server.Port,
			"tlsEnabled": tlsEnabled,
		})
	}

	logger.InfoWithData(ctx, "app_started", map[string]any{
		"port": app.cfg.Server.Port,
	})
	err = http.ServeTLS(listener, app.cfg.Server.Router, certFile, keyFile)
	if err != nil {
		logger.FatalWithData(ctx, "app_crashed", err, map[string]any{
			"port":       app.cfg.Server.Port,
			"tlsEnabled": tlsEnabled,
		})
	}
}
