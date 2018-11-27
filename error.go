package raptor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/si3nloong/raptor/validator"
)

type (
	errorClaim struct {
		Error struct {
			Code        string      `json:"code" xml:"code"`
			Message     string      `json:"message" xml:"message"`
			Debug       string      `json:"debug,omitempty" xml:"debug"`
			Description interface{} `json:"description,omitempty" xml:"description"`
		} `json:"error" xml:"error"`
	}
)

// APIError :
type APIError struct {
	Inner       error
	Code        string
	Message     string
	Description interface{}
	Detail      interface{}
	isDebug     bool
}

// MarshalJSON :
func (e *APIError) MarshalJSON() (b []byte, err error) {
	r := new(errorClaim)
	r.Error.Code = e.Code
	r.Error.Message = e.Message
	if e.isDebug && e.Inner != nil {
		r.Error.Debug = e.Inner.Error()
	}
	r.Error.Description = e.Description
	b, err = ffjson.Marshal(r)
	return
}

func (e *APIError) Error() string {
	buff := new(strings.Builder)
	if e.Inner != nil {
		buff.WriteString(fmt.Sprintf("debug=%s, ", e.Inner.Error()))
	}
	buff.WriteString(fmt.Sprintf("code=%s, message=%s", e.Code, e.Message))
	if e.Description != nil {
		buff.WriteString(fmt.Sprintf(", description=%v", e.Description))
	}
	return buff.String()
}

// HTTPError :
type HTTPError struct {
	StatusCode int
	Message    interface{}
	Inner      error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", e.StatusCode, e.Message)
}

// DefaultErrorHandler :
func (r *Raptor) DefaultErrorHandler(ctx *Context, err error) {
	statusCode := ctx.RequestCtx.Response.StatusCode()
	if statusCode <= 0 {
		statusCode = fasthttp.StatusInternalServerError
	}

	switch ve := err.(type) {
	case *APIError:
		if ve.Code == "" {
			ve.Code = strconv.Itoa(statusCode)
		}
		if ve.Message == "" {
			ve.Message = http.StatusText(statusCode)
		}
		err = ctx.Response().compileResponse(ve, statusCode)
	case *validator.ValidationError:

		err = ctx.Response().compileResponse(ve, statusCode)
	default:
		err = ctx.Response().compileResponse(ve, statusCode)
	}
}
