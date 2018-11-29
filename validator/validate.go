package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// ValidationError :
type ValidationError struct {
	errs validator.ValidationErrors
}

// Error :
func (ve ValidationError) Error() string {
	return ""
}

// MarshalJSON :
func (ve ValidationError) MarshalJSON() ([]byte, error) {
	if len(ve.errs) == 0 {
		return []byte(`null`), nil
	}

	var buf bytes.Buffer
	buf.WriteByte(123)

	for _, err := range ve.errs {
		if buf.Len() > 1 {
			buf.WriteByte(44)
		}

		message, isExist := ValidationErrorMessages[err.Tag()]
		if isExist && reflect.TypeOf(message).Kind() == reflect.Map {
			// m, isOK := message.(map[string]string)[err.Kind().String()]
			// if !isOK {
			// 	m = ValidationErrorMessages["default"].(string)
			// }
			// n := strings.Replace(err.Field(), first, name, -1)
			// m = strings.Replace(m, ":field", n, 1)
			// m = strings.Replace(m, ":value", err.Param(), 1)

			// errs[n] = m
		} else {
			if !isExist {
				message = ValidationErrorMessages["default"]
			}
			ns := err.Namespace()
			message = strings.Replace(message.(string), ":field", ns, 1)
			buf.WriteByte(34)
			buf.WriteString(ns)
			buf.WriteByte(34)
			buf.WriteByte(58)
			buf.WriteString(fmt.Sprintf("%q", message))
		}
	}

	buf.WriteByte(125)

	return buf.Bytes(), nil
}

// Validate : validate fields
func Validate(tag string, i interface{}) error {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		tag = "json"
	}
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if errs := v.Struct(i); errs != nil {
		return &ValidationError{errs: errs.(validator.ValidationErrors)}
	}
	return nil
}
