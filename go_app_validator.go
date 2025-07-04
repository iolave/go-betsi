package goapp

import (
	"reflect"

	"github.com/pingolabscl/go-app/pkg/errors"
)

// recursiveValidation validates any type recursively
// using the go-playground/validator rules. It returns
// an error if the validation fails.
func (app *App) recursiveValidation(v any) *errors.Error {
	switch kind := reflect.TypeOf(v).Kind(); kind {
	case reflect.Slice:
		valueOf := reflect.ValueOf(v)
		length := valueOf.Len()
		for i := range length {
			v := valueOf.Index(i).Interface()
			err := app.recursiveValidation(v)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		valueOf := reflect.ValueOf(v)
		iter := valueOf.MapRange()
		for iter.Next() {
			v := iter.Value()
			if err := app.recursiveValidation(v.Interface()); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if err := app.validator.Struct(v); err != nil {
			return errors.NewWithNameAndErr(
				"validation_error",
				"failed to validate struct",
				err,
			)
		}
	}

	return nil
}
