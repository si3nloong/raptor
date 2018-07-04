package validator

import (
	"log"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

func diveIn(errs map[string]interface{}, r reflect.Type, fields []string, err validator.FieldError) error {
	first := fields[0]
	if strings.IndexRune(first, '[') > 0 {
		first = first[:strings.IndexRune(first, '[')]
	}

	var v reflect.StructField
	var name string

	if r.Kind() == reflect.Slice {
		r = r.Elem()
	}

	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	v, _ = r.FieldByName(first)
	name = v.Tag.Get("json")
	name = strings.Replace(name, ",omitempty", "", -1)
	if name == "" {
		name = first
	}

	if len(fields) > 1 {
		_, isNil := errs[name]
		if !isNil {
			errs[name] = make(map[string]interface{})
		}

		if v.Type.Kind() == reflect.Ptr {
			return diveIn(errs[name].(map[string]interface{}), v.Type.Elem(), fields[1:], err)
		}

		return diveIn(errs[name].(map[string]interface{}), v.Type, fields[1:], err)
	}

	message, isExist := ValidationErrorMessages[err.Tag()]
	if isExist && reflect.TypeOf(message).Kind() == reflect.Map {
		m, isOK := message.(map[string]string)[err.Kind().String()]
		if !isOK {
			m = ValidationErrorMessages["default"].(string)
		}
		n := strings.Replace(err.Field(), first, name, -1)
		m = strings.Replace(m, ":field", n, -1)
		m = strings.Replace(m, ":value", err.Param(), -1)
		log.Println("err", n)

		errs[n] = m
	} else {
		if !isExist {
			message = ValidationErrorMessages["default"]
		}
		n := strings.Replace(err.Field(), first, name, -1)
		n = strings.Replace(n, ",omitempty", "", -1)
		message = strings.Replace(message.(string), ":field", n, -1)
		errs[n] = message
	}

	return nil
}
