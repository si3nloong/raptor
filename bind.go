package raptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	errReflectIsNotPointer  = errors.New("Struct is not a pointer")
	errUnSupportedMediaType = errors.New("Unsupported media type to bind")
)

var (
	typeOfByte = reflect.TypeOf([]byte(nil))
)

func getString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func loadValue(v reflect.Value, b []string) error {
	t := v.Type()
	switch t.Kind() {
	case reflect.String:
		var str string
		if len(b) > 0 {
			str = b[0]
		}
		v.SetString(str)
	case reflect.Bool:
		x, err := strconv.ParseBool(getString(b))
		if err != nil {
			return err
		}
		v.SetBool(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(getString(b), 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(x) {
			return fmt.Errorf("overflow integer value, %v", x)
		}
		v.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(getString(b), 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(x) {
			return fmt.Errorf("overflow unsigned integer value, %v", x)
		}
		v.SetUint(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(getString(b), 64)
		if err != nil {
			return err
		}
		if v.OverflowFloat(x) {
			return fmt.Errorf("overflow float value , %v", x)
		}
		v.SetFloat(x)
	case reflect.Interface:
		var it interface{}
		if err := json.Unmarshal([]byte(getString(b)), &it); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(it))
	case reflect.Slice, reflect.Array:
		if t == typeOfByte {
			v.SetBytes([]byte(getString(b)))
			return nil
		}

		arr := make([]json.RawMessage, 0)
		if err := json.Unmarshal([]byte(getString(b)), &arr); err != nil {
			return err
		}

		vv := reflect.MakeSlice(t, len(arr), len(arr))
		vi := reflect.New(vv.Type())
		vi.Elem().Set(vv)
		for i := 0; i < len(arr); i++ {
			if err := loadValue(vi.Elem().Index(i), []string{string(arr[i])}); err != nil {
				return err
			}
		}
		v.Set(vi.Elem())
	case reflect.Ptr:
		if getString(b) == "" {
			v.Set(reflect.New(t))
			return nil
		}
		if err := loadValue(v.Elem(), b); err != nil {
			return err
		}
	case reflect.Struct:
		vv := reflect.New(v.Type())
		if err := json.Unmarshal(json.RawMessage(getString(b)), vv.Interface()); err != nil {
			return err
		}
		v.Set(vv.Elem())
	}

	return nil
}

func bindQuery(tag string, v reflect.Value, l map[string][]string) error {
	vi := reflect.Indirect(v)
	vv := reflect.New(vi.Type())

	for i := 0; i < vv.Elem().NumField(); i++ {
		fv := vv.Elem().Field(i)
		ft := vv.Type().Elem().Field(i)

		tag := strings.Split(ft.Tag.Get(tag), ",")
		name := strings.TrimSpace(tag[0])
		if name == "-" {
			continue
		}

		if name == "" {
			name = ft.Name
		}

		if strings.TrimSpace(name) == "" {
			name = ft.Name
		}

		val, isOk := l[name]
		if !isOk {
			continue
		}

		if err := loadValue(fv, val); err != nil {
			return err
		}
	}

	vi.Set(vv.Elem())
	return nil
}
