package raptor

// Group :
type Group struct {
	prefix      string
	middlewares []MiddlewareFunc
	raptor      *Raptor
}

// Group :
func (g *Group) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	return &Group{
		prefix:      g.prefix + prefix,
		middlewares: append(g.middlewares, middleware...),
		raptor:      g.raptor,
	}
}

func (g *Group) addRoute(method string, path string, h HandlerFunc, middleware ...MiddlewareFunc) *Group {
	fullPath := path
	if g.prefix != "" {
		fullPath = g.prefix + fullPath
	}
	m := make([]MiddlewareFunc, 0, len(g.middlewares)+len(middleware))
	m = append(m, g.middlewares...)
	m = append(m, middleware...)
	g.raptor.addRoute(method, fullPath, h, m...)
	return g
}

// CONNECT :
func (g *Group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(CONNECT, path, h, m...)
}

// TRACE :
func (g *Group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(TRACE, path, h, m...)
}

// HEAD :
func (g *Group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(HEAD, path, h, m...)
}

// OPTIONS :
func (g *Group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(OPTIONS, path, h, m...)
}

// GET :
func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(GET, path, h, m...)
}

// POST :
func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(POST, path, h, m...)
}

// PUT :
func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(PUT, path, h, m...)
}

// PATCH :
func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(PATCH, path, h, m...)
}

// DELETE :
func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(DELETE, path, h, m...)
}

// Any :
func (g *Group) Any(path string, h HandlerFunc, m ...MiddlewareFunc) {
	methods := [...]string{CONNECT, PROPFIND, TRACE, HEAD, OPTIONS, GET, POST, PUT, PATCH, DELETE}
	for i := 0; i < len(methods); i++ {
		g.addRoute(methods[i], path, h, m...)
	}
}
