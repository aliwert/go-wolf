package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliwert/go-wolf/pkg/context"
)

func TestRouteBuilder(t *testing.T) {
	router := New()

	// Test route building with method chaining (using the Route struct)
	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Add a route and get Route instance for chaining
	router.Handle("GET", "/test", handler)

	// Test that the route was added correctly
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	context.Release(c)
}

func TestRouteNameBuilder(t *testing.T) {
	router := New()

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Create route info for building
	routeInfo := &RouteInfo{
		Method:  "GET",
		Path:    "/test",
		Handler: handler,
	}

	// Create route and test name setting
	route := &Route{
		info:   routeInfo,
		router: router,
	}

	// Test name chaining
	namedRoute := route.Name("test-route")

	if namedRoute.info.Name != "test-route" {
		t.Errorf("Expected route name 'test-route', got '%s'", namedRoute.info.Name)
	}

	// Test that the route is stored in named routes
	if router.namedRoutes == nil {
		t.Error("Expected namedRoutes to be initialized")
	}

	if router.namedRoutes["test-route"] != routeInfo {
		t.Error("Expected route to be stored in namedRoutes")
	}
}

func TestRouteInfoBuilder(t *testing.T) {
	// Test building RouteInfo struct
	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middleware1 := func(c *context.Context) error {
		c.SetHeader("X-Middleware-1", "true")
		return c.Next()
	}

	middleware2 := func(c *context.Context) error {
		c.SetHeader("X-Middleware-2", "true")
		return c.Next()
	}

	// Build route info
	routeInfo := &RouteInfo{
		Method:     "GET",
		Path:       "/users/:id",
		Name:       "user.show",
		Middleware: []context.HandlerFunc{middleware1, middleware2},
		Constraints: map[string]Constraint{
			"id": IsNumeric,
		},
		Subdomain: "api",
	}

	// Set handler and verify it's set correctly
	routeInfo.Handler = handler

	if routeInfo.Handler == nil {
		t.Errorf("Expected handler to be set")
	}
	if routeInfo.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", routeInfo.Method)
	}

	if routeInfo.Path != "/users/:id" {
		t.Errorf("Expected path '/users/:id', got '%s'", routeInfo.Path)
	}

	if routeInfo.Name != "user.show" {
		t.Errorf("Expected name 'user.show', got '%s'", routeInfo.Name)
	}

	if len(routeInfo.Middleware) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(routeInfo.Middleware))
	}

	if len(routeInfo.Constraints) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(routeInfo.Constraints))
	}

	if routeInfo.Subdomain != "api" {
		t.Errorf("Expected subdomain 'api', got '%s'", routeInfo.Subdomain)
	}
}

func TestRouteGroupBuilder(t *testing.T) {
	router := New()

	middleware1 := func(c *context.Context) error {
		c.SetHeader("X-Group", "api")
		return c.Next()
	}

	middleware2 := func(c *context.Context) error {
		c.SetHeader("X-Version", "v1")
		return c.Next()
	}

	// Build route group
	group := &Group{
		router:     router,
		prefix:     "/api/v1",
		middleware: []context.HandlerFunc{middleware1, middleware2},
	}

	// Test group properties
	if group.router != router {
		t.Error("Expected group to reference router")
	}

	if group.prefix != "/api/v1" {
		t.Errorf("Expected prefix '/api/v1', got '%s'", group.prefix)
	}

	if len(group.middleware) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(group.middleware))
	}
}

func TestRouteChaining(t *testing.T) {
	router := New()

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Create route info
	routeInfo := &RouteInfo{
		Method:  "GET",
		Path:    "/test",
		Handler: handler,
	}

	// Create route
	route := &Route{
		info:   routeInfo,
		router: router,
	}

	// Test method chaining
	chainedRoute := route.Name("test-route")

	// Should return the same route instance for chaining
	if chainedRoute != route {
		t.Error("Expected route chaining to return same instance")
	}

	// Test that properties are set correctly
	if route.info.Name != "test-route" {
		t.Errorf("Expected name 'test-route', got '%s'", route.info.Name)
	}
}

func TestMultipleRouteNames(t *testing.T) {
	router := New()

	handler1 := func(c *context.Context) error {
		return c.String(http.StatusOK, "route1")
	}

	handler2 := func(c *context.Context) error {
		return c.String(http.StatusOK, "route2")
	}

	// Create first route
	routeInfo1 := &RouteInfo{
		Method:  "GET",
		Path:    "/route1",
		Handler: handler1,
	}

	route1 := &Route{
		info:   routeInfo1,
		router: router,
	}

	// Create second route
	routeInfo2 := &RouteInfo{
		Method:  "GET",
		Path:    "/route2",
		Handler: handler2,
	}

	route2 := &Route{
		info:   routeInfo2,
		router: router,
	}

	// Name both routes
	route1.Name("first-route")
	route2.Name("second-route")

	// Test that both routes are stored
	if router.namedRoutes["first-route"] != routeInfo1 {
		t.Error("Expected first route to be stored")
	}

	if router.namedRoutes["second-route"] != routeInfo2 {
		t.Error("Expected second route to be stored")
	}

	// Test that they don't interfere with each other
	if len(router.namedRoutes) != 2 {
		t.Errorf("Expected 2 named routes, got %d", len(router.namedRoutes))
	}
}

func TestGroupNestedBuilder(t *testing.T) {
	router := New()

	middleware1 := func(c *context.Context) error {
		c.SetHeader("X-Level", "1")
		return c.Next()
	}

	middleware2 := func(c *context.Context) error {
		c.SetHeader("X-Level", c.Header("X-Level")+"-2")
		return c.Next()
	}

	middleware3 := func(c *context.Context) error {
		c.SetHeader("X-Level", c.Header("X-Level")+"-3")
		return c.Next()
	}

	// Build nested group structure
	api := router.Group("/api", middleware1)
	v1 := api.Group("/v1", middleware2)
	users := v1.Group("/users", middleware3)

	// Test that middleware is properly combined
	if len(api.middleware) != 1 {
		t.Errorf("Expected api to have 1 middleware, got %d", len(api.middleware))
	}

	if len(v1.middleware) != 2 {
		t.Errorf("Expected v1 to have 2 middleware, got %d", len(v1.middleware))
	}

	if len(users.middleware) != 3 {
		t.Errorf("Expected users to have 3 middleware, got %d", len(users.middleware))
	}

	// Test that prefixes are properly combined
	if users.prefix != "/api/v1/users" {
		t.Errorf("Expected nested prefix '/api/v1/users', got '%s'", users.prefix)
	}
}

func TestConstraintBuilder(t *testing.T) {
	// Test building complex constraints
	constraint := And(
		Or(IsAlpha, IsNumeric),
		MinLength(2),
		MaxLength(10),
	)

	// Test the built constraint
	tests := []struct {
		input    string
		expected bool
	}{
		{"ab", true},           // alpha, length 2
		{"12", true},           // numeric, length 2
		{"a", false},           // too short
		{"abcdefghijk", false}, // too long
		{"ab1", false},         // mixed alpha-numeric
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Complex constraint(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestRouteOptionsBuilder(t *testing.T) {
	// Test building RouterOptions
	notFoundHandler := func(c *context.Context) error {
		return c.String(http.StatusNotFound, "Custom Not Found")
	}

	methodNotAllowedHandler := func(c *context.Context) error {
		return c.String(http.StatusMethodNotAllowed, "Custom Method Not Allowed")
	}

	options := &RouterOptions{
		NotFoundHandler:         notFoundHandler,
		MethodNotAllowedHandler: methodNotAllowedHandler,
		EnableCaching:           true,
		CacheSize:               1000,
	}

	// Test that all options are set correctly
	if options.NotFoundHandler == nil {
		t.Error("Expected NotFoundHandler to be set")
	}

	if options.MethodNotAllowedHandler == nil {
		t.Error("Expected MethodNotAllowedHandler to be set")
	}

	if !options.EnableCaching {
		t.Error("Expected EnableCaching to be true")
	}

	if options.CacheSize != 1000 {
		t.Errorf("Expected CacheSize 1000, got %d", options.CacheSize)
	}
}

func TestMethodChainBuilder(t *testing.T) {
	router := New()

	// Test building routes with method chain pattern
	group := router.Group("/api")

	// Simulate a builder pattern for routes
	routes := []struct {
		method  string
		path    string
		handler context.HandlerFunc
	}{
		{"GET", "/users", func(c *context.Context) error {
			return c.String(http.StatusOK, "get users")
		}},
		{"POST", "/users", func(c *context.Context) error {
			return c.String(http.StatusCreated, "create user")
		}},
		{"PUT", "/users/:id", func(c *context.Context) error {
			return c.String(http.StatusOK, "update user")
		}},
		{"DELETE", "/users/:id", func(c *context.Context) error {
			return c.String(http.StatusNoContent, "delete user")
		}},
	}

	// Build all routes
	for _, route := range routes {
		switch route.method {
		case "GET":
			group.GET(route.path, route.handler)
		case "POST":
			group.POST(route.path, route.handler)
		case "PUT":
			group.PUT(route.path, route.handler)
		case "DELETE":
			group.DELETE(route.path, route.handler)
		}
	}

	// Test that routes work
	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "get users" {
		t.Errorf("Expected body 'get users', got '%s'", w.Body.String())
	}

	context.Release(c)
}

// Benchmark tests for builder patterns
func BenchmarkRouteBuilder(b *testing.B) {
	router := New()

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		routeInfo := &RouteInfo{
			Method:  "GET",
			Path:    "/test",
			Handler: handler,
		}

		route := &Route{
			info:   routeInfo,
			router: router,
		}

		route.Name("test-route")
	}
}

func BenchmarkGroupBuilder(b *testing.B) {
	router := New()

	middleware := func(c *context.Context) error {
		return c.Next()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		api := router.Group("/api", middleware)
		v1 := api.Group("/v1", middleware)
		v1.Group("/users", middleware)
	}
}

func BenchmarkConstraintBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		constraint := And(
			Or(IsAlpha, IsNumeric),
			MinLength(2),
			MaxLength(10),
		)
		constraint("test")
	}
}
