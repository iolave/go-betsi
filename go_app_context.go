package goapp

import (
	"context"
	"reflect"

	"github.com/pingolabscl/go-app/pkg/errors"
)

const (
	ctx_app_key               = "app"
	app_err_name              = "app_error"
	app_err_msg_app_not_found = "unable to get app from context because it was not found"
	app_err_msg_wrong_type    = "unable to get app from context because the type is not *goapp.App"
)

// GetFromContext returns the app from a context. If the
// context does not contain an app, or the type of the
// app is not *goapp.App, then an error is returned.
func GetFromContext(ctx context.Context) (*App, *errors.Error) {
	app := ctx.Value(ctx_app_key)
	if app == nil {
		return nil, errors.NewWithName(app_err_name, app_err_msg_app_not_found)
	}

	if reflect.TypeFor[*App]() != reflect.TypeOf(app) {
		return nil, errors.NewWithName(app_err_name, app_err_msg_wrong_type)
	}

	return app.(*App), nil
}

func setContext(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, ctx_app_key, app)
}
