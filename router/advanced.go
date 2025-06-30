package router

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aliwert/go-wolf/pkg/context"
)

// RouteConstraint represents advanced parameter constraints
type RouteConstraint struct {
	Name    string
	Pattern *regexp.Regexp
	Checker func(string) bool
}

// RouteBuilder provides a fluent interface for building routes
type RouteBuilder struct {
	router      *Router
	method      string
	path        string
	handler     context.HandlerFunc
	middleware  []context.HandlerFunc
	constraints map[string]RouteConstraint
	name        string
	subdomain   string
}

// NewRouteBuilder creates a new route builder
func (r *Router) NewRoute() *RouteBuilder {
	return &RouteBuilder{
		router:      r,
		constraints: make(map[string]RouteConstraint),
	}
}

// Method sets the HTTP method for the route
func (rb *RouteBuilder) Method(method string) *RouteBuilder {
	rb.method = method
	return rb
}

// Path sets the path pattern for the route
func (rb *RouteBuilder) Path(path string) *RouteBuilder {
	rb.path = path
	return rb
}

// Handler sets the handler function for the route
func (rb *RouteBuilder) Handler(handler context.HandlerFunc) *RouteBuilder {
	rb.handler = handler
	return rb
}

// Middleware adds middleware to the route
func (rb *RouteBuilder) Middleware(mw ...context.HandlerFunc) *RouteBuilder {
	rb.middleware = append(rb.middleware, mw...)
	return rb
}

// Name sets a name for the route (for URL generation)
func (rb *RouteBuilder) Name(name string) *RouteBuilder {
	rb.name = name
	return rb
}

// Subdomain restricts the route to a specific subdomain
func (rb *RouteBuilder) Subdomain(subdomain string) *RouteBuilder {
	rb.subdomain = subdomain
	return rb
}

// Where adds parameter constraints
func (rb *RouteBuilder) Where(param string, constraint interface{}) *RouteBuilder {
	var rc RouteConstraint
	rc.Name = param

	switch c := constraint.(type) {
	case string:
		// Regex pattern
		rc.Pattern = regexp.MustCompile(c)
		rc.Checker = func(value string) bool {
			return rc.Pattern.MatchString(value)
		}
	case func(string) bool:
		// Custom function
		rc.Checker = c
	case *regexp.Regexp:
		// Compiled regex
		rc.Pattern = c
		rc.Checker = func(value string) bool {
			return rc.Pattern.MatchString(value)
		}
	}

	rb.constraints[param] = rc
	return rb
}

// WhereNumber constrains parameter to be numeric
func (rb *RouteBuilder) WhereNumber(param string) *RouteBuilder {
	return rb.Where(param, IsNumeric)
}

// WhereAlpha constrains parameter to be alphabetic
func (rb *RouteBuilder) WhereAlpha(param string) *RouteBuilder {
	return rb.Where(param, IsAlpha)
}

// WhereAlphaNumeric constrains parameter to be alphanumeric
func (rb *RouteBuilder) WhereAlphaNumeric(param string) *RouteBuilder {
	return rb.Where(param, IsAlphaNumeric)
}

// WhereIn constrains parameter to be one of the provided values
func (rb *RouteBuilder) WhereIn(param string, values ...string) *RouteBuilder {
	valueMap := make(map[string]bool)
	for _, v := range values {
		valueMap[v] = true
	}
	return rb.Where(param, func(value string) bool {
		return valueMap[value]
	})
}

// WhereUUID constrains parameter to be a valid UUID
func (rb *RouteBuilder) WhereUUID(param string) *RouteBuilder {
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return rb.Where(param, uuidPattern)
}

// WhereSlug constrains parameter to be a URL-friendly slug
func (rb *RouteBuilder) WhereSlug(param string) *RouteBuilder {
	return rb.Where(param, IsSlug)
}

// Build finalizes and registers the route
func (rb *RouteBuilder) Build() *Route {
	if rb.method == "" || rb.path == "" || rb.handler == nil {
		panic("method, path, and handler are required")
	}

	// Create route info
	info := &RouteInfo{
		Method:     rb.method,
		Path:       rb.path,
		Name:       rb.name,
		Handler:    rb.handler,
		Middleware: rb.middleware,
	}

	// Store constraints in the router
	if len(rb.constraints) > 0 {
		if rb.router.constraints == nil {
			rb.router.constraints = make(map[string]map[string]Constraint)
		}
		if rb.router.constraints[rb.path] == nil {
			rb.router.constraints[rb.path] = make(map[string]Constraint)
		}
		for param, constraint := range rb.constraints {
			rb.router.constraints[rb.path][param] = constraint.Checker
		}
	}

	// Register the route
	rb.router.registerAdvancedRoute(info)

	return &Route{
		info:   info,
		router: rb.router,
	}
}

// registerAdvancedRoute registers a route with advanced features
func (r *Router) registerAdvancedRoute(info *RouteInfo) {
	// Store route info
	if r.routes == nil {
		r.routes = make([]*RouteInfo, 0)
	}
	r.routes = append(r.routes, info)

	// Store named route
	if info.Name != "" {
		if r.namedRoutes == nil {
			r.namedRoutes = make(map[string]*RouteInfo)
		}
		r.namedRoutes[info.Name] = info
	}

	// Register with the underlying router
	r.Handle(info.Method, info.Path, info.Handler, info.Middleware...)
}

// URL generates a URL for a named route
func (r *Router) URL(name string, params map[string]string) (string, error) {
	if r.namedRoutes == nil {
		return "", fmt.Errorf("no named routes registered")
	}

	route, exists := r.namedRoutes[name]
	if !exists {
		return "", fmt.Errorf("route '%s' not found", name)
	}

	path := route.Path
	for key, value := range params {
		placeholder := ":" + key
		path = strings.Replace(path, placeholder, value, -1)
	}

	return path, nil
}

// GetRoutes returns all registered routes
func (r *Router) GetRoutes() []*RouteInfo {
	return r.routes
}

// GetNamedRoutes returns all named routes
func (r *Router) GetNamedRoutes() map[string]*RouteInfo {
	return r.namedRoutes
}

// SetNotFoundHandler sets a custom 404 handler
func (r *Router) SetNotFoundHandler(handler context.HandlerFunc) {
	r.notFoundHandler = handler
}

// SetMethodNotAllowedHandler sets a custom 405 handler
func (r *Router) SetMethodNotAllowedHandler(handler context.HandlerFunc) {
	r.methodNotAllowedHandler = handler
}

// Resource creates RESTful routes for a resource
func (r *Router) Resource(name string, controller interface{}) {
	// Create standard RESTful routes:
	// GET    /resource          -> Index
	// GET    /resource/create   -> Create (form)
	// POST   /resource          -> Store
	// GET    /resource/:id      -> Show
	// GET    /resource/:id/edit -> Edit (form)
	// PUT    /resource/:id      -> Update
	// DELETE /resource/:id      -> Destroy

	basePath := "/" + name

	// Index - GET /resource
	if handler := getControllerMethod(controller, "Index"); handler != nil {
		r.Handle("GET", basePath, handler)
	}

	// Create form - GET /resource/create
	if handler := getControllerMethod(controller, "Create"); handler != nil {
		r.Handle("GET", basePath+"/create", handler)
	}

	// Store - POST /resource
	if handler := getControllerMethod(controller, "Store"); handler != nil {
		r.Handle("POST", basePath, handler)
	}

	// Show - GET /resource/:id
	if handler := getControllerMethod(controller, "Show"); handler != nil {
		r.Handle("GET", basePath+"/:id", handler)
	}

	// Edit form - GET /resource/:id/edit
	if handler := getControllerMethod(controller, "Edit"); handler != nil {
		r.Handle("GET", basePath+"/:id/edit", handler)
	}

	// Update - PUT /resource/:id
	if handler := getControllerMethod(controller, "Update"); handler != nil {
		r.Handle("PUT", basePath+"/:id", handler)
	}

	// Destroy - DELETE /resource/:id
	if handler := getControllerMethod(controller, "Destroy"); handler != nil {
		r.Handle("DELETE", basePath+"/:id", handler)
	}
}

// ResourceController defines the interface for resource controllers
type ResourceController interface {
	Index(*context.Context) error   // GET /resource
	Create(*context.Context) error  // GET /resource/create
	Store(*context.Context) error   // POST /resource
	Show(*context.Context) error    // GET /resource/:id
	Edit(*context.Context) error    // GET /resource/:id/edit
	Update(*context.Context) error  // PUT /resource/:id
	Destroy(*context.Context) error // DELETE /resource/:id
}

// getControllerMethod uses reflection to get a method from a controller
func getControllerMethod(controller interface{}, methodName string) context.HandlerFunc {
	if controller == nil {
		return nil
	}

	// Get the reflect value and type
	v := reflect.ValueOf(controller)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if it's a struct
	if v.Kind() != reflect.Struct {
		return nil
	}

	// Get the method
	method := reflect.ValueOf(controller).MethodByName(methodName)
	if !method.IsValid() {
		return nil
	}

	// Check if method has correct signature: func(*context.Context) error
	methodType := method.Type()
	if methodType.NumIn() != 1 || methodType.NumOut() != 1 {
		return nil
	}

	// Check parameter type
	paramType := methodType.In(0)
	contextType := reflect.TypeOf((*context.Context)(nil))
	if paramType != contextType {
		return nil
	}

	// Check return type
	returnType := methodType.Out(0)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if returnType != errorType {
		return nil
	}

	// Return a wrapper function
	return func(c *context.Context) error {
		results := method.Call([]reflect.Value{reflect.ValueOf(c)})
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}
		return nil
	}
}

// ResourceOptions allows customization of resource routes
type ResourceOptions struct {
	Only   []string // only create these routes
	Except []string // create all routes except these
	Prefix string   // path prefix for the resource
	Name   string   // name prefix for route names
}

// ResourceWithOptions creates RESTful routes with custom options
func (r *Router) ResourceWithOptions(name string, controller interface{}, options ResourceOptions) {
	basePath := "/" + name
	if options.Prefix != "" {
		basePath = "/" + strings.Trim(options.Prefix, "/") + basePath
	}

	routeMap := map[string]struct {
		method     string
		path       string
		methodName string
	}{
		"index":   {"GET", basePath, "Index"},
		"create":  {"GET", basePath + "/create", "Create"},
		"store":   {"POST", basePath, "Store"},
		"show":    {"GET", basePath + "/:id", "Show"},
		"edit":    {"GET", basePath + "/:id/edit", "Edit"},
		"update":  {"PUT", basePath + "/:id", "Update"},
		"destroy": {"DELETE", basePath + "/:id", "Destroy"},
	}

	// Determine which routes to create
	createRoute := func(routeName string) bool {
		// If Only is specified, only create those routes
		if len(options.Only) > 0 {
			for _, only := range options.Only {
				if only == routeName {
					return true
				}
			}
			return false
		}

		// If Except is specified, create all except those
		if len(options.Except) > 0 {
			for _, except := range options.Except {
				if except == routeName {
					return false
				}
			}
		}

		return true
	}

	// Create routes
	for routeName, route := range routeMap {
		if createRoute(routeName) {
			if handler := getControllerMethod(controller, route.methodName); handler != nil {
				routeBuilder := r.NewRoute().
					Method(route.method).
					Path(route.path).
					Handler(handler)

				if options.Name != "" {
					routeBuilder.Name(options.Name + "." + routeName)
				} else {
					routeBuilder.Name(name + "." + routeName)
				}

				routeBuilder.Build()
			}
		}
	}
}
