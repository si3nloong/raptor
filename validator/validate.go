package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validate : validate fields
func Validate(i interface{}) (*map[string]interface{}, error) {
	var validate *validator.Validate
	validate = validator.New()
	registerValidation(validate)
	registerAlias(validate)

	err := validate.Struct(i)

	if err != nil {

		errs := make(map[string]interface{})

		r := reflect.Indirect(reflect.ValueOf(i)).Type()

		structName := r.Name()
		for _, eachError := range err.(validator.ValidationErrors) {
			ns := eachError.Namespace()
			if structName != "" {
				ns = strings.Replace(ns, fmt.Sprintf("%s.", structName), "", 1)
			}
			fields := strings.Split(ns, ".")[:]

			diveIn(errs, r, fields, eachError)
		}
		return &errs, errors.New("Validation errors")
	}
	return nil, nil
}
