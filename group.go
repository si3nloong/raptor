package raptor

// Group :
type Group struct {
	prefix      string
	middlewares []MiddlewareFunc
	router      *Router
}

func (g *Group) addRoute(method string, path string, h HandlerFunc, middleware ...MiddlewareFunc) *Group {
	fullPath := path
	if g.prefix != "" {
		fullPath = g.prefix + fullPath
	}
	m := make([]MiddlewareFunc, 0, len(g.middlewares)+len(middleware))
	m = append(m, g.middlewares...)
	m = append(m, middleware...)
	return g
}

// CONNECT :
func (g *Group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(CONNECT, path, h, m...)
	return
}

// TRACE :
func (g *Group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(TRACE, path, h, m...)
	return
}

// HEAD :
func (g *Group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(HEAD, path, h, m...)
	return
}

// OPTIONS :
func (g *Group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(OPTIONS, path, h, m...)
	return
}

// GET :
func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(GET, path, h, m...)
	return
}

// POST :
func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(POST, path, h, m...)
	return
}

// PUT :
func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(PUT, path, h, m...)
	return
}

// PATCH :
func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(PATCH, path, h, m...)
	return
}

// DELETE :
func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(DELETE, path, h, m...)
	return
}

// Any :
func (g *Group) Any(path string, h HandlerFunc, m ...MiddlewareFunc) {
	g.addRoute(GET, path, h, m...)
	return
}
