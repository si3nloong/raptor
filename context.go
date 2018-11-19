package raptor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/erikdubbelboer/fasthttp"
	json "github.com/pquerna/ffjson/ffjson"
	"github.com/si3nloong/raptor/validator"
)

// Context :
type Context struct {
	*fasthttp.RequestCtx
	isDebug bool
}

// ErrUnSupportedMediaType :
var ErrUnSupportedMediaType = errors.New("Unsupported media type to bind")

// Response :
func (c *Context) Response() *Response {
	return &Response{c.RequestCtx}
}

// Bind :
func (c *Context) Bind(dst interface{}) error {
	v := reflect.ValueOf(dst)
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("layout is not addressable")
	}

	query := b2s(bytes.TrimSpace(c.QueryArgs().QueryString()))
	if query != "" {
		values, err := url.ParseQuery(query)
		if err != nil {
			return err
		}
		if err := bindQuery("query", v, values); err != nil {
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
	case MIMEApplicationJSON:
		if err := json.Unmarshal(c.Request.Body(), dst); err != nil {
			return err
		}
	default:
		return ErrUnSupportedMediaType
	}

	return nil
}

// QueryString :
func (c *Context) QueryString() string {
	buf := new(bytes.Buffer)
	c.QueryArgs().VisitAll(func(k, v []byte) {
		buf.WriteString(fmt.Sprintf("%s=%s&", k, v))
	})
	return strings.TrimRight(buf.String(), "&")
}

// IsMethod :
func (c *Context) IsMethod(method string) bool {
	method = strings.TrimSpace(strings.ToLower(method))
	m := strings.TrimSpace(strings.ToLower(string(c.Method())))
	return method == m
}

// Param :
func (c *Context) Param(key string) string {
	var str string
	switch vi := c.UserValue(key).(type) {
	case []byte:
		str = string(vi)
	case string:
		str = vi
	case bool:
		str = fmt.Sprintf("%t", vi)
	case int, int8, int16, int32, int64:
		str = fmt.Sprintf("%d", vi)
	case uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("%d", vi)
	case float32:
		str = strconv.FormatFloat(float64(vi), 'f', 10, 64)
	case float64:
		str = strconv.FormatFloat(vi, 'f', 10, 64)
	case time.Time:
		str = vi.Format(time.RFC3339)
	case fmt.Stringer:
		str = vi.String()
	default:
		str = fmt.Sprintf("%v", vi)
	}

	return str
}

// Redirect :
func (c *Context) Redirect(uri string, statusCode ...int) error {
	code := fasthttp.StatusMovedPermanently
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	c.RequestCtx.Redirect(uri, code)
	return nil
}

// Validate :
func (c *Context) Validate(i interface{}) error {
	if c.IsMethod("GET") {
		return validator.Validate("query", i)
	}
	return validator.Validate("json", i)
}

// SetCookie :
func (c *Context) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		c.RequestCtx.Response.Header.Set("Set-Cookie", v)
	}
}

// HTML :
func (c *Context) HTML(html string, statusCode ...int) error {
	return c.HTMLBlob([]byte(html), statusCode...)
}

// HTMLBlob :
func (c *Context) HTMLBlob(b []byte, statusCode ...int) error {
	httpStatusCode := fasthttp.StatusOK
	if len(statusCode) > 0 {
		httpStatusCode = statusCode[0]
	}
	c.RequestCtx.Response.Header.Set(HeaderContentType, "text/html; charset=utf-8")
	c.RequestCtx.Response.Header.SetStatusCode(httpStatusCode)
	c.Write(b)
	return nil
}

// Render :
func (c *Context) Render(b []byte) error {
	c.RequestCtx.Response.Header.Set(HeaderContentType, "text/html; charset=utf-8")
	c.RequestCtx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	c.Write(b)
	return nil
}

// NewAPIError :
func (c *Context) NewAPIError(err error, params ...interface{}) error {
	e := new(APIError)
	e.Inner = err
	if len(params) > 0 {
		x, _ := params[0].(string)
		e.Code = x
	}
	if len(params) > 1 {
		x, _ := params[1].(string)
		e.Message = x
	}
	if len(params) > 2 {
		e.Description = params[2]
	}
	e.isDebug = c.isDebug
	return e
}
