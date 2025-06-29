package router

import (
	"net/http"

	"github.com/aliwert/go-wolf/pkg/context"
)

// Router represents the HTTP router
type Router struct {
	trees                   map[string]*node
	routes                  []*RouteInfo
	namedRoutes             map[string]*RouteInfo
	notFoundHandler         context.HandlerFunc
	methodNotAllowedHandler context.HandlerFunc
	constraints             map[string]map[string]Constraint // path -> param -> constraint
}

// RouteInfo represents information about a registered route
type RouteInfo struct {
	Method      string
	Path        string
	Name        string
	Handler     context.HandlerFunc
	Middleware  []context.HandlerFunc
	Constraints map[string]Constraint
	Subdomain   string
}

// Route represents a route with additional metadata
type Route struct {
	info   *RouteInfo
	router *Router
}

// Name sets the name for this route
func (r *Route) Name(name string) *Route {
	r.info.Name = name
	if r.router.namedRoutes == nil {
		r.router.namedRoutes = make(map[string]*RouteInfo)
	}
	r.router.namedRoutes[name] = r.info
	return r
}

// New creates a new router
func New() *Router {
	return &Router{
		trees: make(map[string]*node),
	}
}

// Handle registers a new request handle with the given path and method
func (r *Router) Handle(method, path string, handler context.HandlerFunc, middleware ...context.HandlerFunc) {
	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handler == nil {
		panic("handler must not be nil")
	}

	// Get or create tree for method
	root := r.trees[method]
	if root == nil {
		root = &node{}
		r.trees[method] = root
	}

	// Build middleware chain
	finalHandler := handler
	for i := len(middleware) - 1; i >= 0; i-- {
		mw := middleware[i]
		next := finalHandler
		finalHandler = func(mw context.HandlerFunc, next context.HandlerFunc) context.HandlerFunc {
			return func(c *context.Context) error {
				c.SetNext(next)
				return mw(c)
			}
		}(mw, next)
	}

	root.addRoute(path, finalHandler)
}

// Group creates a new route group with the given prefix
func (r *Router) Group(prefix string, middleware ...context.HandlerFunc) *Group {
	return &Group{
		router:     r,
		prefix:     prefix,
		middleware: middleware,
	}
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request, c *context.Context) {
	method := req.Method
	path := req.URL.Path

	if root := r.trees[method]; root != nil {
		if handle, params, _ := root.getValue(path); handle != nil {
			if params != nil {
				c.SetParams(params)
			}
			if err := handle(c); err != nil {
				if errorHandler := c.GetErrorHandler(); errorHandler != nil {
					errorHandler(c, err)
				}
			}
			return
		}
	}

	// Handle 405 Method Not Allowed
	for m := range r.trees {
		if m != method {
			if root := r.trees[m]; root != nil {
				if handle, _, _ := root.getValue(path); handle != nil {
					c.Writer.WriteHeader(http.StatusMethodNotAllowed)
					c.Writer.Write([]byte("Method Not Allowed"))
					return
				}
			}
		}
	}

	// Handle 404 Not Found
	c.Writer.WriteHeader(http.StatusNotFound)
	c.Writer.Write([]byte("Not Found"))
}

// RouterOptions holds router configuration
type RouterOptions struct {
	NotFoundHandler         context.HandlerFunc
	MethodNotAllowedHandler context.HandlerFunc
	EnableCaching           bool
	CacheSize               int
}

// Utility functions for the radix tree

// countParams counts the number of parameters in a path
func countParams(path string) uint8 {
	var n uint8
	for i := range []byte(path) {
		switch path[i] {
		case ':', '*':
			n++
		}
	}
	return n
}

// longestCommonPrefix finds the longest common prefix between two strings
func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// findWildcard finds wildcard segment and checks if it is valid
func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' && c != '*' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}
