package raptor

import (
	"encoding/json"
	"encoding/xml"

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

	// ErrorResponse :
	ErrorResponse interface {
		error
		json.Marshaler
		xml.Marshaler
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

// Blob :
func (r *Response) Blob(contentType string, b []byte, statusCode ...int) error {
	r.Response.Header.Set(HeaderContentType, contentType)
	r.SetStatusCode(fasthttp.StatusOK)
	if len(statusCode) > 0 {
		r.SetStatusCode(statusCode[0])
	}
	r.Write(b)
	return nil
}

// XML :
func (r *Response) XML(data interface{}, statusCode ...int) error {
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
	accept := string(r.Request.Header.Peek(HeaderAccept))
	switch accept {
	case MIMEApplicationXML, MIMETextXML:
		return r.XML(data, statusCode)
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

// NoContent : 204 responding with no content
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

func (r *Response) compileError(statusCode int, err error) error {
	r.Response.SetStatusCode(statusCode)
	return err
}

// BadRequest : 400
func (r *Response) BadRequest(err error) error {
	return r.compileError(fasthttp.StatusBadRequest, err)
}

// Unauthorized : 401
func (r *Response) Unauthorized(err error) error {
	return r.compileError(fasthttp.StatusUnauthorized, err)
}

// Forbidden : 403
func (r *Response) Forbidden(err error) error {
	return r.compileError(fasthttp.StatusForbidden, err)
}

// NotFound : 404
func (r *Response) NotFound(err error) error {
	return r.compileError(fasthttp.StatusNotFound, err)
}

// MethodNotAllowed : 405
func (r *Response) MethodNotAllowed(err error) error {
	return r.compileError(fasthttp.StatusMethodNotAllowed, err)
}

// NotAcceptable : 406
func (r *Response) NotAcceptable(err error) error {
	return r.compileError(fasthttp.StatusNotAcceptable, err)
}

// RequestTimeout : 408
func (r *Response) RequestTimeout(err error) error {
	return r.compileError(fasthttp.StatusRequestTimeout, err)
}

// Conflict : 409
func (r *Response) Conflict(err error) error {
	return r.compileError(fasthttp.StatusConflict, err)
}

// Gone : 410
func (r *Response) Gone(err error) error {
	return r.compileError(fasthttp.StatusGone, err)
}

// LengthRequired : 411
func (r *Response) LengthRequired(err error) error {
	return r.compileError(fasthttp.StatusLengthRequired, err)
}

// PreconditionFailed : 412
func (r *Response) PreconditionFailed(err error) error {
	return r.compileError(fasthttp.StatusPreconditionFailed, err)
}

// PayloadTooLarge : 413
func (r *Response) PayloadTooLarge(err error) error {
	return r.compileError(fasthttp.StatusRequestEntityTooLarge, err)
}

// UnsupportedMediaType : 415
func (r *Response) UnsupportedMediaType(err error) error {
	return r.compileError(fasthttp.StatusUnsupportedMediaType, err)
}

// ExpectationFailed : 417
func (r *Response) ExpectationFailed(err error) error {
	return r.compileError(fasthttp.StatusExpectationFailed, err)
}

// UnprocessableEntity : 422
func (r *Response) UnprocessableEntity(err error) error {
	return r.compileError(422, err)
}

// Locked : 423
func (r *Response) Locked(err error) error {
	return r.compileError(423, err)
}

// InternalServerError : 500
func (r *Response) InternalServerError(err error) error {
	return r.compileError(fasthttp.StatusInternalServerError, err)
}

// BadGateway : 501
func (r *Response) BadGateway(err error) error {
	return r.compileError(fasthttp.StatusBadGateway, err)
}

// ServiceUnavailable : 503
func (r *Response) ServiceUnavailable(err error) error {
	return r.compileError(fasthttp.StatusServiceUnavailable, err)
}

// GatewayTimeout : 504
func (r *Response) GatewayTimeout(err error) error {
	return r.compileError(fasthttp.StatusGatewayTimeout, err)
}
