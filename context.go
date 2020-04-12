package raptor

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	json "encoding/json"

	"github.com/si3nloong/raptor/validator"
	"github.com/valyala/fasthttp"
)

// Context :
type Context struct {
	*fasthttp.RequestCtx
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
func (c *Context) Param(key string) (str string) {
	switch vi := c.UserValue(key).(type) {
	case []byte:
		str = string(vi)
	case string:
		str = vi
	case bool:
		str = fmt.Sprintf("%t", vi)
	case int:
		str = strconv.FormatInt(int64(vi), 10)
	case int8:
		str = strconv.FormatInt(int64(vi), 10)
	case int16:
		str = strconv.FormatInt(int64(vi), 10)
	case int32:
		str = strconv.FormatInt(int64(vi), 10)
	case int64:
		str = strconv.FormatInt(int64(vi), 10)
	case uint:
		str = strconv.FormatUint(uint64(vi), 10)
	case uint8:
		str = strconv.FormatUint(uint64(vi), 10)
	case uint16:
		str = strconv.FormatUint(uint64(vi), 10)
	case uint32:
		str = strconv.FormatUint(uint64(vi), 10)
	case uint64:
		str = strconv.FormatUint(uint64(vi), 10)
	case float32:
		str = strconv.FormatFloat(float64(vi), 'f', 10, 32)
	case float64:
		str = strconv.FormatFloat(vi, 'f', 10, 64)
	case time.Time:
		str = vi.Format(time.RFC3339)
	case fmt.Stringer:
		str = vi.String()
	case nil:
		str = ""
	default:
		str = fmt.Sprintf("%v", vi)
	}
	return
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

func (c *Context) Blob(contentType string, b []byte, statusCode ...int) error {
	httpStatusCode := fasthttp.StatusOK
	if len(statusCode) > 0 {
		httpStatusCode = statusCode[0]
	}
	c.RequestCtx.Response.Header.Set(HeaderContentType, contentType)
	c.RequestCtx.Response.Header.SetStatusCode(httpStatusCode)
	if _, err := c.Write(b); err != nil {
		return err
	}
	return nil
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

// HTML :
func (c *Context) HTML(html string, statusCode ...int) error {
	return c.HTMLBlob([]byte(html), statusCode...)
}

// Render :
func (c *Context) Render(b []byte) error {
	return c.HTMLBlob(b, fasthttp.StatusOK)
}

// NoContent :
func (c *Context) NoContent(statusCode ...int) error {
	c.ResetBody()
	if len(statusCode) > 0 {
		c.RequestCtx.Response.Header.SetStatusCode(statusCode[0])
		return nil
	}
	c.RequestCtx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}

// Response :
func (c *Context) Response() *Response {
	return &Response{c.RequestCtx}
}

// JSON :
func (c *Context) JSON(value interface{}, statusCode ...int) error {
	code := fasthttp.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	c.RequestCtx.Response.Header.Set(HeaderContentType, MIMEApplicationJSON)
	c.RequestCtx.Response.Header.SetStatusCode(code)
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
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
		e.Detail = params[2]
	}
	return e
}
