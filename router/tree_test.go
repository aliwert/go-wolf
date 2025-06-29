package router

import (
	"net/http"
	"testing"

	"github.com/aliwert/go-wolf/pkg/context"
)

func TestNodeAddRoute(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Test simple path insertion
	root.addRoute("/users", handler)
	if root.path != "/users" {
		t.Errorf("Expected path '/users', got '%s'", root.path)
	}

	if root.handle == nil {
		t.Error("Expected handler to be set")
	}
}

func TestNodeAddRouteWithParameter(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Test parameter path insertion
	root.addRoute("/users/:id", handler)

	// Test path lookup
	handle, params, _ := root.getValue("/users/123")
	if handle == nil {
		t.Error("Expected to find handler")
	}

	if params == nil || params["id"] != "123" {
		t.Errorf("Expected param id=123, got %v", params)
	}
}

func TestNodeAddRouteWithWildcard(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "wildcard")
	}

	// Test wildcard path insertion
	root.addRoute("/static/*filepath", handler)

	// Test wildcard lookup
	handle, params, _ := root.getValue("/static/css/main.css")
	if handle == nil {
		t.Error("Expected to find wildcard handler")
	}

	if params == nil || params["filepath"] != "/css/main.css" {
		t.Errorf("Expected param filepath=/css/main.css, got %v", params)
	}
}

func TestNodeGetValueMultipleRoutes(t *testing.T) {
	root := &node{}

	handler1 := func(c *context.Context) error {
		return c.String(http.StatusOK, "users")
	}

	handler2 := func(c *context.Context) error {
		return c.String(http.StatusOK, "user")
	}

	// Add multiple routes
	root.addRoute("/users", handler1)
	root.addRoute("/users/:id", handler2)

	// Test exact match
	handle, params, _ := root.getValue("/users")
	if handle == nil {
		t.Error("Expected to find handler for /users")
	}

	if len(params) > 0 {
		t.Error("Expected no params for exact match")
	}

	// Test parameterized match
	handle, params, _ = root.getValue("/users/123")
	if handle == nil {
		t.Error("Expected to find handler for /users/123")
	}

	if params == nil || params["id"] != "123" {
		t.Errorf("Expected param id=123, got %v", params)
	}
}

func TestNodeGetValueWildcardParams(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "files")
	}

	root.addRoute("/static/*filepath", handler)

	// Test nested file path
	handle, params, _ := root.getValue("/static/css/main.css")
	if handle == nil {
		t.Error("Expected to find handler")
	}

	if params == nil || params["filepath"] != "/css/main.css" {
		t.Errorf("Expected filepath=/css/main.css, got %v", params)
	}
}

func TestNodeGetValueConflictingRoutes(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "handler")
	}

	// Add routes in an order that works with the router's constraints
	root.addRoute("/users/:id", handler)
	root.addRoute("/users/:id/posts", handler)

	// Test parameterized route
	handle, params, _ := root.getValue("/users/123")
	if handle == nil {
		t.Error("Expected to find handler for /users/123")
	}

	if params == nil || params["id"] != "123" {
		t.Errorf("Expected param id=123, got %v", params)
	}

	// Test nested parameterized route
	handle, params, _ = root.getValue("/users/123/posts")
	if handle == nil {
		t.Error("Expected to find handler for /users/123/posts")
	}

	if params == nil || params["id"] != "123" {
		t.Errorf("Expected param id=123, got %v", params)
	}
}

func TestNodeGetValueNotFound(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "found")
	}

	root.addRoute("/users/:id", handler)

	// Test path that doesn't match
	handle, params, _ := root.getValue("/products/123")
	if handle != nil {
		t.Error("Expected not to find handler for unregistered path")
	}

	if params != nil {
		t.Error("Expected no params for unmatched path")
	}
}

func TestNodePriority(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Add routes and check priority increments
	initialPriority := root.priority

	root.addRoute("/users", handler)
	if root.priority <= initialPriority {
		t.Error("Expected priority to increase after adding route")
	}

	prevPriority := root.priority
	root.addRoute("/users/:id", handler)
	if root.priority <= prevPriority {
		t.Error("Expected priority to increase after adding second route")
	}
}

func TestNodeMaxParams(t *testing.T) {
	root := &node{}

	handler := func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Add route with one parameter
	root.addRoute("/users/:id", handler)
	if root.maxParams < 1 {
		t.Error("Expected maxParams to be at least 1")
	}

	// Add route with more parameters
	root.addRoute("/users/:id/posts/:postId", handler)
	if root.maxParams < 2 {
		t.Error("Expected maxParams to be at least 2")
	}
}

func TestCountParams(t *testing.T) {
	tests := []struct {
		path     string
		expected uint8
	}{
		{"/users", 0},
		{"/users/:id", 1},
		{"/users/:id/posts/:postId", 2},
		{"/static/*filepath", 1},
		{"/api/:version/users/:id", 2},
	}

	for _, test := range tests {
		result := countParams(test.path)
		if result != test.expected {
			t.Errorf("countParams(%s) = %d, expected %d", test.path, result, test.expected)
		}
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"abc", "", 0},
		{"", "abc", 0},
		{"abc", "abc", 3},
		{"abc", "ab", 2},
		{"ab", "abc", 2},
		{"abc", "def", 0},
		{"/users", "/users/123", 6},
		{"/users/123", "/users", 6},
	}

	for _, test := range tests {
		result := longestCommonPrefix(test.a, test.b)
		if result != test.expected {
			t.Errorf("longestCommonPrefix(%s, %s) = %d, expected %d", test.a, test.b, result, test.expected)
		}
	}
}

func TestFindWildcard(t *testing.T) {
	tests := []struct {
		path          string
		expectedWild  string
		expectedIndex int
		expectedValid bool
	}{
		{"/users", "", -1, false},
		{"/users/:id", ":id", 7, true},
		{"/static/*filepath", "*filepath", 8, true},
		{"/users/:id/posts", ":id", 7, true},
		{"/bad/:param:invalid", ":param", 5, false},
	}

	for _, test := range tests {
		wildcard, index, valid := findWildcard(test.path)
		if wildcard != test.expectedWild {
			t.Errorf("findWildcard(%s) wildcard = %s, expected %s", test.path, wildcard, test.expectedWild)
		}
		if index != test.expectedIndex {
			t.Errorf("findWildcard(%s) index = %d, expected %d", test.path, index, test.expectedIndex)
		}
		if valid != test.expectedValid {
			t.Errorf("findWildcard(%s) valid = %t, expected %t", test.path, valid, test.expectedValid)
		}
	}
}
