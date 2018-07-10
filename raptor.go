package raptor

import (
	"io/ioutil"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

func init() {
	color.Blue(`
		Version: 0.1

		'########:::::'###::::'########::'########::'#######::'########::
		##.... ##:::'## ##::: ##.... ##:... ##..::'##.... ##: ##.... ##:
		##:::: ##::'##:. ##:: ##:::: ##:::: ##:::: ##:::: ##: ##:::: ##:
		########::'##:::. ##: ########::::: ##:::: ##:::: ##: ########::
		##.. ##::: #########: ##.....:::::: ##:::: ##:::: ##: ##.. ##:::
		##::. ##:: ##.... ##: ##::::::::::: ##:::: ##:::: ##: ##::. ##::
		##:::. ##: ##:::: ##: ##::::::::::: ##::::. #######:: ##:::. ##:
		..:::::..::..:::::..::..::::::::::::..::::::.......:::..:::::..::
	`)
}

// HTTP REQUEST METHODS :
const (
	CONNECT  = "CONNECT"
	DELETE   = "DELETE"
	GET      = "GET"
	HEAD     = "HEAD"
	OPTIONS  = "OPTIONS"
	PATCH    = "PATCH"
	POST     = "POST"
	PROPFIND = "PROPFIND"
	PUT      = "PUT"
	TRACE    = "TRACE"
)

// Raptor :
type Raptor struct {
	router       *fasthttprouter.Router
	middlewares  []MiddlewareFunc
	ErrorHandler func(c *Context, err error)
	Debug        bool
}

// HandlerFunc :
type HandlerFunc func(*Context) error

// MiddlewareFunc :
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// New :
func New() *Raptor {
	r := &Raptor{
		router:      fasthttprouter.New(),
		middlewares: make([]MiddlewareFunc, 0),
	}
	r.ErrorHandler = r.DefaultErrorHandler
	return r
}

// Handler :
func (r *Raptor) Handler() HandlerFunc {
	return func(c *Context) error {
		r.router.Handler(c.RequestCtx)
		return nil
	}
}

// Use :
func (r *Raptor) Use(middleware ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware...)
}

// Static :
func (r *Raptor) Static(prefix, path string) *Raptor {
	r.GET(prefix, func(c *Context) error {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		c.Write(b)
		return nil
	})
	return r
}

// Group :
func (r *Raptor) Group(path string, middleware ...MiddlewareFunc) *Group {
	g := new(Group)
	g.prefix = path
	g.middlewares = append(g.middlewares, middleware...)
	g.raptor = r
	return g
}

// HEAD :
func (r *Raptor) HEAD(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(HEAD, path, handler, middleware...)
}

// OPTIONS :
func (r *Raptor) OPTIONS(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(OPTIONS, path, handler, middleware...)
}

// GET :
func (r *Raptor) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(GET, path, handler, middleware...)
}

// POST :
func (r *Raptor) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(POST, path, handler, middleware...)
}

// PUT :
func (r *Raptor) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(PUT, path, handler, middleware...)
}

// PATCH :
func (r *Raptor) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(PATCH, path, handler, middleware...)
}

// DELETE :
func (r *Raptor) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(DELETE, path, handler, middleware...)
}

func (r *Raptor) mergeHandler(handler HandlerFunc, middlewares ...MiddlewareFunc) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		h := handler
		c := &Context{ctx}
		for _, m := range middlewares {
			h = m(h)
		}
		if err := h(c); err != nil {
			r.ErrorHandler(c, err)
		}
	})
}

func (r *Raptor) addRoute(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	m := make([]MiddlewareFunc, 0, len(r.middlewares)+len(middleware))
	m = append(m, r.middlewares...)
	m = append(m, middleware...)

	cb := r.mergeHandler(handler, m...)
	switch method {
	case OPTIONS:
		r.router.OPTIONS(path, cb)
	case GET:
		r.router.GET(path, cb)
	case POST:
		r.router.POST(path, cb)
	case PUT:
		r.router.PUT(path, cb)
	case PATCH:
		r.router.PATCH(path, cb)
	case DELETE:
		r.router.DELETE(path, cb)
	default:
		r.router.Handle(method, path, cb)
	}

	return
}

// Start :
func (r *Raptor) Start(port string, handler ...HandlerFunc) error {
	cb := r.router.Handler
	if len(handler) > 0 {
		cb = func(ctx *fasthttp.RequestCtx) {
			c := &Context{ctx}
			if err := handler[0](c); err != nil {
				r.ErrorHandler(c, err)
				return
			}
		}
	}
	log.Println("fasthttp server started on", port)
	return fasthttp.ListenAndServe(port, cb)
}
