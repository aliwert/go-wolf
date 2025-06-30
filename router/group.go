package router

import "github.com/aliwert/go-wolf/pkg/context"

// Group represents a route group
type Group struct {
	router     *Router
	prefix     string
	middleware []context.HandlerFunc
}

// Group creates a sub-group with additional prefix
func (g *Group) Group(prefix string, middleware ...context.HandlerFunc) *Group {
	return &Group{
		router:     g.router,
		prefix:     g.prefix + prefix,
		middleware: append(g.middleware, middleware...),
	}
}

// Use adds middleware to the group
func (g *Group) Use(middleware ...context.HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// GET adds a GET route to the group
func (g *Group) GET(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("GET", path, handler, middleware...)
}

// POST adds a POST route to the group
func (g *Group) POST(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("POST", path, handler, middleware...)
}

// PUT adds a PUT route to the group
func (g *Group) PUT(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("PUT", path, handler, middleware...)
}

// DELETE adds a DELETE route to the group
func (g *Group) DELETE(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("DELETE", path, handler, middleware...)
}

// PATCH adds a PATCH route to the group
func (g *Group) PATCH(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("PATCH", path, handler, middleware...)
}

// HEAD adds a HEAD route to the group
func (g *Group) HEAD(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("HEAD", path, handler, middleware...)
}

// OPTIONS adds an OPTIONS route to the group
func (g *Group) OPTIONS(path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	g.handle("OPTIONS", path, handler, middleware...)
}

// handle adds a route with the given method to the group
func (g *Group) handle(method, path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	// Combine group middleware with route-specific middleware
	allMiddleware := append(g.middleware, middleware...)
	fullPath := g.prefix + path
	g.router.Handle(method, fullPath, handler, allMiddleware...)
}
