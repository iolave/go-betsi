package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/iolave/go-errors"
)

// validate is the validator throughout the app.
// It is initialized here cuz it caches structs,
// which is useful for writing appfactory handlers.
var validate = validator.New(
	validator.WithRequiredStructEnabled(),
)

// ValidateRecursively validates any type recursively
// using the go-playground/validator rules. It returns
// an error if the validation fails.
//
// error is of type *errors.Error.
func ValidateRecursively(v any) error {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}

	switch kind := t.Kind(); kind {
	case reflect.Ptr:
		v := reflect.ValueOf(v)
		return ValidateRecursively(v.Elem().Interface())
	case reflect.Slice:
		valueOf := reflect.ValueOf(v)
		length := valueOf.Len()
		for i := range length {
			v := valueOf.Index(i).Interface()
			err := ValidateRecursively(v)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		valueOf := reflect.ValueOf(v)
		iter := valueOf.MapRange()
		for iter.Next() {
			v := iter.Value()
			if err := ValidateRecursively(v.Interface()); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if err := validate.Struct(v); err != nil {
			return errors.NewWithNameAndErr(
				"validation_error",
				"failed to validate struct",
				err,
			)
		}
	}

	return nil
}
