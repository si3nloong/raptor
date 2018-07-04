package raptor

import (
	"io/ioutil"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

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

// Router :
type Router struct {
	router      *fasthttprouter.Router
	middlewares []MiddlewareFunc
}

// NewRouter :
func NewRouter() *Router {
	return &Router{
		router:      fasthttprouter.New(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

// Handler :
func (r *Router) Handler() HandlerFunc {
	return func(c *Context) error {
		r.router.Handler(c.RequestCtx)
		return nil
	}
}

// Use :
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware...)
}

// Static :
func (r *Router) Static(prefix, path string) *Router {
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
func (r *Router) Group(path string, middleware ...MiddlewareFunc) *Group {
	g := new(Group)
	g.prefix = path
	g.middlewares = append(g.middlewares, middleware...)
	g.router = r
	return g
}

// HEAD :
func (r *Router) HEAD(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(HEAD, path, handler, middleware...)
}

// OPTIONS :
func (r *Router) OPTIONS(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(OPTIONS, path, handler, middleware...)
}

// GET :
func (r *Router) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(GET, path, handler, middleware...)
}

// POST :
func (r *Router) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(POST, path, handler, middleware...)
}

// PUT :
func (r *Router) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(PUT, path, handler, middleware...)
}

// PATCH :
func (r *Router) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(PATCH, path, handler, middleware...)
}

// DELETE :
func (r *Router) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	r.addRoute(DELETE, path, handler, middleware...)
}

func (r *Router) mergeHandler(handler HandlerFunc, middlewares ...MiddlewareFunc) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		h := handler
		c := &Context{ctx}
		for _, m := range middlewares {
			h = m(h)
		}
		if err := h(c); err != nil {
		}
	})
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
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
