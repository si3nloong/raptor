package raptor

import (
	"bytes"
	"encoding"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"errors"

	"github.com/ajg/form"
	"github.com/gorilla/schema"
)

var (
	errReflectIsNotPointer  = errors.New("Struct is not a pointer")
	errUnsupportedMediaType = errors.New("Unsupported media type to bind")
)

var (
	typeOfByte = reflect.TypeOf([]byte(nil))
)

// Bind :
func (c *Context) Bind(dst interface{}) error {
	v := reflect.ValueOf(dst)
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("layout is not addressable")
	}

	c.QueryArgs().QueryString()
	query := b2s(bytes.TrimSpace(c.QueryArgs().QueryString()))
	if query != "" {
		values, err := url.ParseQuery(query)
		if err != nil {
			return err
		}

		dec := schema.NewDecoder()
		dec.SetAliasTag("query")
		if err := dec.Decode(dst, values); err != nil {
			return err
		}
	}

	if c.IsMethod(GET) {
		return nil
	}

	paths := bytes.Split(c.Request.Header.Peek(HeaderContentType), []byte{59})
	switch b2s(bytes.TrimSpace(paths[0])) {
	case MIMEApplicationForm, MIMEMultipartForm:
		if err := form.DecodeString(&dst, string(c.Request.Body())); err != nil {
			return err
		}
	case MIMEApplicationXML:
		if err := xml.Unmarshal(c.Request.Body(), dst); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(c.Request.Body(), dst); err != nil {
			return err
		}
	}

	return nil
}

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
		if x, isOk := vv.Interface().(encoding.TextUnmarshaler); isOk {
			x.UnmarshalText([]byte(getString(b)))
			v.Set(vv.Elem())
		}
	}

	return nil
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	return v
}

func bindQuery(tag string, v reflect.Value, l map[string][]string) error {
	vi := indirect(v)
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
