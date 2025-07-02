package goapp

import (
	"context"
	"reflect"
)

const appKey = "app"

func GetFromContext(ctx context.Context) *App {
	ctxApp := ctx.Value(appKey)
	if ctxApp == nil {
		return nil
	}

	if reflect.TypeFor[*App]() != reflect.TypeOf(ctxApp) {
		return nil
	}

	return ctxApp.(*App)
}

func setContext(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appKey, app)
}
