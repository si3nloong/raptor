package raptor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/si3nloong/raptor/validator"
	"github.com/valyala/fasthttp"
)

type (
	errorClaim struct {
		Error struct {
			Code    string      `json:"code" xml:"code"`
			Message string      `json:"message" xml:"message"`
			Debug   string      `json:"debug,omitempty" xml:"debug,omitempty"`
			Detail  interface{} `json:"detail,omitempty" xml:"detail,omitempty"`
		} `json:"error" xml:"error"`
	}
)

// APIError :
type APIError struct {
	Inner   error
	Code    string
	Message string
	Detail  interface{}
	isDebug bool
}

// MarshalJSON :
func (e *APIError) MarshalJSON() (b []byte, err error) {
	r := new(errorClaim)
	r.Error.Code = e.Code
	r.Error.Message = e.Message
	if e.isDebug && e.Inner != nil {
		r.Error.Debug = e.Inner.Error()
	}
	r.Error.Detail = e.Detail
	b, err = ffjson.Marshal(r)
	return
}

func (e *APIError) Error() string {
	blr := new(strings.Builder)
	if e.Inner != nil {
		blr.WriteString(fmt.Sprintf("debug=%s, ", e.Inner.Error()))
	}
	blr.WriteString(fmt.Sprintf("code=%s, message=%s", e.Code, e.Message))
	if e.Detail != nil {
		blr.WriteString(fmt.Sprintf(", detail=%v", e.Detail))
	}
	return blr.String()
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
func DefaultErrorHandler(ctx *Context, err error) {
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
