package raptor

import (
	"encoding/xml"

	json "github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

type (
	itemClaim struct {
		XMLName xml.Name    `json:"-" xml:"item" form:"-"`
		Item    interface{} `json:"item" xml:"item" form:"item"`
	}

	itemsClaim struct {
		XMLName xml.Name    `json:"-" xml:"items" form:"-"`
		Items   interface{} `json:"items" xml:"items" form:"items"`
		Meta    struct {
			Count  uint   `json:"count,omitempty"`
			Total  uint   `json:"total,omitempty"`
			Cursor string `json:"cursor,omitempty"`
		} `json:"meta,omitempty"`
	}

	errorClaim struct {
		Error struct {
			Code        string      `json:"code" xml:"code"`
			Message     interface{} `json:"message" xml:"message"`
			Debug       string      `json:"debug,omitempty" xml:"debug"`
			Description interface{} `json:"description,omitempty" xml:"description"`
		} `json:"error" xml:"error"`
	}
)

// Paginator :
type Paginator interface {
	Count() uint
	NextCursor() string
}

// Response :
type Response struct {
	*fasthttp.RequestCtx
}

// SetStatusCode :
func (r *Response) SetStatusCode(statusCode int) *Response {
	r.Response.SetStatusCode(statusCode)
	return r
}

// XML :
func (r *Response) XML(data interface{}, statusCode ...int) error {
	r.SetStatusCode(fasthttp.StatusOK)
	if len(statusCode) > 0 {
		r.SetStatusCode(statusCode[0])
	}
	r.Response.Header.Set(HeaderContentType, MIMEApplicationXML)
	b, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	r.Write(b)
	return nil
}

// JSON :
func (r *Response) JSON(data interface{}, statusCode ...int) error {
	r.SetStatusCode(fasthttp.StatusOK)
	if len(statusCode) > 0 {
		r.SetStatusCode(statusCode[0])
	}
	r.Response.Header.Set(HeaderContentType, MIMEApplicationJSON)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Write(b)
	return nil
}

func (r *Response) compileResponse(data interface{}, statusCode int) error {
	switch string(r.Request.Header.Peek(HeaderAccept)) {
	case MIMEApplicationXML:
		return r.XML(data, statusCode)
	case MIMEApplicationJSON:
		return r.JSON(data, statusCode)
	default:
	}
	return r.JSON(data, statusCode)
}

// Custom :
func (r *Response) Custom(data interface{}, statusCode ...int) error {
	code := fasthttp.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return r.compileResponse(data, code)
}

// Success :
func (r *Response) Success(data interface{}) error {
	return r.compileResponse(data, fasthttp.StatusOK)
}

// NoContent : responding with no content
func (r *Response) NoContent() error {
	r.Response.SetStatusCode(fasthttp.StatusNoContent)
	r.Response.SetBody([]byte(nil))
	return nil
}

// Item :
func (r *Response) Item(data interface{}) error {
	i := new(itemClaim)
	i.Item = data
	return r.Success(i)
}

// Collection :
func (r *Response) Collection(data interface{}) error {
	i := new(itemsClaim)
	i.Items = data
	return r.Success(i)
}

// Paginate :
func (r *Response) Paginate(p Paginator, data interface{}) error {
	i := new(itemsClaim)
	i.Items = data
	i.Meta.Count = p.Count()
	i.Meta.Cursor = p.NextCursor()
	return r.Success(i)
}

func (r *Response) compileError(statusCode int, code string, err error) error {
	e := new(errorClaim)
	e.Error.Code = code
	e.Error.Debug = err.Error()
	return r.compileResponse(e, statusCode)
}

// BadRequest :
func (r *Response) BadRequest(code string, err error) error {
	return r.compileError(fasthttp.StatusBadRequest, code, err)
}

// NotFound :
func (r *Response) NotFound(code string, err error) error {
	return r.compileError(fasthttp.StatusNotFound, code, err)
}

// Forbidden :
func (r *Response) Forbidden(code string, err error) error {
	return r.compileError(fasthttp.StatusForbidden, code, err)
}

// Unauthorized :
func (r *Response) Unauthorized(code string, err error) error {
	return r.compileError(fasthttp.StatusUnauthorized, code, err)
}

// Conflict :
func (r *Response) Conflict(code string, err error) error {
	return r.compileError(fasthttp.StatusConflict, code, err)
}

// Gone :
func (r *Response) Gone(code string, err error) error {
	return r.compileError(fasthttp.StatusGone, code, err)
}

// ExpectationFailed :
func (r *Response) ExpectationFailed(code string, err error) error {
	return r.compileError(fasthttp.StatusExpectationFailed, code, err)
}

// InternalServerError :
func (r *Response) InternalServerError(code string, err error) error {
	return r.compileError(fasthttp.StatusInternalServerError, code, err)
}
