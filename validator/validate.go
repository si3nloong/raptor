package validator

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// ValidationError :
type ValidationError struct {
	errs validator.ValidationErrors
}

// Error :
func (ve ValidationError) Error() string {
	return ve.errs.Error()
}

// MarshalJSON :
func (ve ValidationError) MarshalJSON() ([]byte, error) {
	if len(ve.errs) == 0 {
		return []byte(`null`), nil
	}

	buf := new(bytes.Buffer)
	buf.WriteRune('{')

	for i, err := range ve.errs {
		if i > 0 {
			buf.WriteRune(',')
		}

		msg, ok := ValidationErrorMessages[err.Tag()]
		if !ok {
			msg = ValidationErrorMessages["default"]
		}

		ns := err.Namespace()
		msg = strings.Replace(msg, ":field", ns, 1)
		msg = strings.Replace(msg, ":value", err.Param(), 1)
		buf.WriteString(strconv.Quote(ns))
		buf.WriteRune(':')
		buf.WriteString(strconv.Quote(msg))
	}

	buf.WriteRune('}')

	return buf.Bytes(), nil
}

// Validate : validate fields
func Validate(tag string, i interface{}) error {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		tag = "json"
	}
	vldr := validator.New()
	vldr.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err := vldr.RegisterValidation("required_if", validateRequiredIf); err != nil {
		panic(err)
	}
	if err := vldr.Struct(i); err != nil {
		return ValidationError{errs: err.(validator.ValidationErrors)}
	}
	return nil
}
