package router

import (
	"testing"

	"github.com/aliwert/go-wolf/pkg/context"
)

func TestRouteUtils_MatchPath(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"/users", "/users", true},
		{"/users", "/posts", false},
		{"/users/:id", "/users/123", true},
		{"/users/:id", "/users", false},
		{"/users/:id", "/users/123/posts", false},
		{"/users/:id/posts/:postId", "/users/123/posts/456", true},
		{"/static/*filepath", "/static/css/main.css", true},
		{"/static/*filepath", "/static", false},
		{"/*all", "/anything/here", true},
	}

	for _, test := range tests {
		result := utils.MatchPath(test.pattern, test.path)
		if result != test.expected {
			t.Errorf("MatchPath(%s, %s) = %t, expected %t", test.pattern, test.path, result, test.expected)
		}
	}
}

func TestRouteUtils_ExtractParams(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		pattern  string
		path     string
		expected map[string]string
	}{
		{"/users/:id", "/users/123", map[string]string{"id": "123"}},
		{"/users/:id/posts/:postId", "/users/123/posts/456", map[string]string{"id": "123", "postId": "456"}},
		{"/static/*filepath", "/static/css/main.css", map[string]string{"filepath": "css/main.css"}},
		{"/users", "/users", map[string]string{}},
		{"/api/:version/users/:id", "/api/v1/users/123", map[string]string{"version": "v1", "id": "123"}},
	}

	for _, test := range tests {
		result := utils.ExtractParams(test.pattern, test.path)
		if len(result) != len(test.expected) {
			t.Errorf("ExtractParams(%s, %s) returned %d params, expected %d", test.pattern, test.path, len(result), len(test.expected))
			continue
		}

		for key, expectedValue := range test.expected {
			if actualValue, exists := result[key]; !exists || actualValue != expectedValue {
				t.Errorf("ExtractParams(%s, %s) param %s = %s, expected %s", test.pattern, test.path, key, actualValue, expectedValue)
			}
		}
	}
}

func TestRouteUtils_ValidatePath(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		path  string
		valid bool
	}{
		{"/users", true},
		{"/users/123", true},
		{"/api/v1/users", true},
		{"", false},
		{"users", false},
		{"/users//123", false},
		{"/", true},
	}

	for _, test := range tests {
		err := utils.ValidatePath(test.path)
		isValid := err == nil
		if isValid != test.valid {
			t.Errorf("ValidatePath(%s) valid = %t, expected %t", test.path, isValid, test.valid)
		}
	}
}

func TestRouteUtils_NormalizePath(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		input    string
		expected string
	}{
		{"", "/"},
		{"/", "/"},
		{"/users", "/users"},
		{"/users/", "/users"},
		{"users", "/users"},
		{"users/", "/users"},
		{"/users//123", "/users/123"},
		{"/api///v1////users", "/api/v1/users"},
	}

	for _, test := range tests {
		result := utils.NormalizePath(test.input)
		if result != test.expected {
			t.Errorf("NormalizePath(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestRouteUtils_CombinePaths(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		base     string
		path     string
		expected string
	}{
		{"/", "/users", "/users"},
		{"/api", "/users", "/api/users"},
		{"/api/", "/users", "/api/users"},
		{"/api", "/users/", "/api/users"},
		{"/api/", "/users/", "/api/users"},
		{"/api/v1", "/", "/api/v1"},
		{"/", "/", "/"},
		{"", "/users", "/users"},
		{"/api", "", "/api"},
	}

	for _, test := range tests {
		result := utils.CombinePaths(test.base, test.path)
		if result != test.expected {
			t.Errorf("CombinePaths(%s, %s) = %s, expected %s", test.base, test.path, result, test.expected)
		}
	}
}

func TestRouteUtils_IsValidMethod(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		method string
		valid  bool
	}{
		{"GET", true},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"PATCH", true},
		{"HEAD", true},
		{"OPTIONS", true},
		{"TRACE", true},
		{"CONNECT", true},
		{"get", false},
		{"INVALID", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsValidMethod(test.method)
		if result != test.valid {
			t.Errorf("IsValidMethod(%s) = %t, expected %t", test.method, result, test.valid)
		}
	}
}

func TestRouteUtils_ParseRoutePattern(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		pattern         string
		expectParams    []string
		expectWild      []string
		expectStatic    []string
		expectHasParams bool
		expectHasWild   bool
	}{
		{
			pattern:         "/users",
			expectParams:    []string{},
			expectWild:      []string{},
			expectStatic:    []string{"users"},
			expectHasParams: false,
			expectHasWild:   false,
		},
		{
			pattern:         "/users/:id",
			expectParams:    []string{"id"},
			expectWild:      []string{},
			expectStatic:    []string{"users"},
			expectHasParams: true,
			expectHasWild:   false,
		},
		{
			pattern:         "/static/*filepath",
			expectParams:    []string{},
			expectWild:      []string{"filepath"},
			expectStatic:    []string{"static"},
			expectHasParams: false,
			expectHasWild:   true,
		},
		{
			pattern:         "/api/:version/users/:id",
			expectParams:    []string{"version", "id"},
			expectWild:      []string{},
			expectStatic:    []string{"api", "users"},
			expectHasParams: true,
			expectHasWild:   false,
		},
	}

	for _, test := range tests {
		result := utils.ParseRoutePattern(test.pattern)

		if len(result.Params) != len(test.expectParams) {
			t.Errorf("ParseRoutePattern(%s) params count = %d, expected %d", test.pattern, len(result.Params), len(test.expectParams))
		}

		if len(result.Wildcards) != len(test.expectWild) {
			t.Errorf("ParseRoutePattern(%s) wildcards count = %d, expected %d", test.pattern, len(result.Wildcards), len(test.expectWild))
		}

		if len(result.StaticParts) != len(test.expectStatic) {
			t.Errorf("ParseRoutePattern(%s) static parts count = %d, expected %d", test.pattern, len(result.StaticParts), len(test.expectStatic))
		}

		if result.HasParams != test.expectHasParams {
			t.Errorf("ParseRoutePattern(%s) HasParams = %t, expected %t", test.pattern, result.HasParams, test.expectHasParams)
		}

		if result.HasWildcard != test.expectHasWild {
			t.Errorf("ParseRoutePattern(%s) HasWildcard = %t, expected %t", test.pattern, result.HasWildcard, test.expectHasWild)
		}
	}
}

func TestRouteUtils_GenerateURL(t *testing.T) {
	utils := NewRouteUtils()

	tests := []struct {
		pattern   string
		params    map[string]string
		expected  string
		shouldErr bool
	}{
		{
			pattern:   "/users/:id",
			params:    map[string]string{"id": "123"},
			expected:  "/users/123",
			shouldErr: false,
		},
		{
			pattern:   "/users/:id/posts/:postId",
			params:    map[string]string{"id": "123", "postId": "456"},
			expected:  "/users/123/posts/456",
			shouldErr: false,
		},
		{
			pattern:   "/users/:id",
			params:    map[string]string{},
			expected:  "",
			shouldErr: true,
		},
		{
			pattern:   "/users/:id",
			params:    map[string]string{"id": "test space"},
			expected:  "/users/test%20space",
			shouldErr: false,
		},
		{
			pattern:   "/users",
			params:    map[string]string{},
			expected:  "/users",
			shouldErr: false,
		},
	}

	for _, test := range tests {
		result, err := utils.GenerateURL(test.pattern, test.params)

		if test.shouldErr {
			if err == nil {
				t.Errorf("GenerateURL(%s, %v) expected error but got none", test.pattern, test.params)
			}
		} else {
			if err != nil {
				t.Errorf("GenerateURL(%s, %v) unexpected error: %v", test.pattern, test.params, err)
			} else if result != test.expected {
				t.Errorf("GenerateURL(%s, %v) = %s, expected %s", test.pattern, test.params, result, test.expected)
			}
		}
	}
}

func TestMiddlewareChain(t *testing.T) {
	order := []string{}

	middleware1 := func(c *context.Context) error {
		order = append(order, "mw1-start")
		err := c.Next()
		order = append(order, "mw1-end")
		return err
	}

	middleware2 := func(c *context.Context) error {
		order = append(order, "mw2-start")
		err := c.Next()
		order = append(order, "mw2-end")
		return err
	}

	handler := func(c *context.Context) error {
		order = append(order, "handler")
		return nil
	}

	chain := NewMiddlewareChain(middleware1, middleware2)
	finalHandler := chain.Build(handler)

	// Create a mock context
	c := &context.Context{}

	// Execute the chain
	finalHandler(c)

	// Check execution order
	expected := []string{"mw1-start", "mw2-start", "handler", "mw2-end", "mw1-end"}
	if len(order) != len(expected) {
		t.Errorf("Expected %d execution steps, got %d", len(expected), len(order))
		return
	}

	for i, step := range expected {
		if i >= len(order) || order[i] != step {
			t.Errorf("Step %d: expected %s, got %s", i, step, order[i])
		}
	}
}

func TestPathMatcher(t *testing.T) {
	// Case sensitive, strict slash
	matcher1 := NewPathMatcher(true, true)

	tests1 := []struct {
		pattern string
		path    string
		match   bool
	}{
		{"/Users", "/users", false},  // case sensitive
		{"/users/", "/users", false}, // strict slash
		{"/users", "/users", true},
	}

	for _, test := range tests1 {
		result := matcher1.Match(test.pattern, test.path)
		if result != test.match {
			t.Errorf("PathMatcher(case=true, strict=true).Match(%s, %s) = %t, expected %t",
				test.pattern, test.path, result, test.match)
		}
	}

	// Case insensitive, non-strict slash
	matcher2 := NewPathMatcher(false, false)

	tests2 := []struct {
		pattern string
		path    string
		match   bool
	}{
		{"/Users", "/users", true},  // case insensitive
		{"/users/", "/users", true}, // non-strict slash
		{"/users", "/users/", true}, // non-strict slash
	}

	for _, test := range tests2 {
		result := matcher2.Match(test.pattern, test.path)
		if result != test.match {
			t.Errorf("PathMatcher(case=false, strict=false).Match(%s, %s) = %t, expected %t",
				test.pattern, test.path, result, test.match)
		}
	}
}

func TestRouteConflictDetector(t *testing.T) {
	detector := NewRouteConflictDetector()

	tests := []struct {
		pattern1  string
		pattern2  string
		conflicts bool
	}{
		{"/users", "/users", true},                // exact same
		{"/users/:id", "/users/:userId", true},    // both have params at same position
		{"/users/:id", "/users/admin", true},      // param vs static at same position
		{"/users/admin", "/users/:id", true},      // static vs param at same position
		{"/users", "/posts", false},               // different static paths
		{"/users/:id", "/users/:id/posts", false}, // different lengths
		{"/static/*file", "/static/*path", true},  // both wildcards
	}

	for _, test := range tests {
		result := detector.DetectConflicts(test.pattern1, test.pattern2)
		if result != test.conflicts {
			t.Errorf("DetectConflicts(%s, %s) = %t, expected %t",
				test.pattern1, test.pattern2, result, test.conflicts)
		}
	}
}

func TestConstraintValidator(t *testing.T) {
	validator := NewConstraintValidator()

	constraints := map[string]Constraint{
		"id":    IsNumeric,
		"email": IsEmail,
		"name":  MinLength(2),
	}

	tests := []struct {
		params map[string]string
		valid  bool
		desc   string
	}{
		{
			params: map[string]string{
				"id":    "123",
				"email": "test@example.com",
				"name":  "John",
			},
			valid: true,
			desc:  "all valid",
		},
		{
			params: map[string]string{
				"id":    "abc",
				"email": "test@example.com",
				"name":  "John",
			},
			valid: false,
			desc:  "invalid id",
		},
		{
			params: map[string]string{
				"id":    "123",
				"email": "invalid-email",
				"name":  "John",
			},
			valid: false,
			desc:  "invalid email",
		},
		{
			params: map[string]string{
				"id":    "123",
				"email": "test@example.com",
				"name":  "J",
			},
			valid: false,
			desc:  "name too short",
		},
		{
			params: map[string]string{
				"email": "test@example.com",
				"name":  "John",
			},
			valid: false,
			desc:  "missing id",
		},
	}

	for _, test := range tests {
		err := validator.ValidateParams(test.params, constraints)
		isValid := err == nil

		if isValid != test.valid {
			t.Errorf("ValidateParams(%s) valid = %t, expected %t", test.desc, isValid, test.valid)
		}
	}
}

func TestCalculateStatistics(t *testing.T) {
	routes := []*RouteInfo{
		{
			Method:      "GET",
			Path:        "/users",
			Name:        "users.index",
			Middleware:  []context.HandlerFunc{func(c *context.Context) error { return nil }},
			Constraints: map[string]Constraint{},
		},
		{
			Method:      "GET",
			Path:        "/users/:id",
			Middleware:  []context.HandlerFunc{},
			Constraints: map[string]Constraint{"id": IsNumeric},
		},
		{
			Method:      "POST",
			Path:        "/users",
			Middleware:  []context.HandlerFunc{},
			Constraints: map[string]Constraint{},
		},
		{
			Method:      "GET",
			Path:        "/static/*filepath",
			Middleware:  []context.HandlerFunc{},
			Constraints: map[string]Constraint{},
		},
	}

	stats := CalculateStatistics(routes)

	if stats.TotalRoutes != 4 {
		t.Errorf("Expected 4 total routes, got %d", stats.TotalRoutes)
	}

	if stats.StaticRoutes != 1 {
		t.Errorf("Expected 1 static route, got %d", stats.StaticRoutes)
	}

	if stats.ParametricRoutes != 1 {
		t.Errorf("Expected 1 parametric route, got %d", stats.ParametricRoutes)
	}

	if stats.WildcardRoutes != 1 {
		t.Errorf("Expected 1 wildcard route, got %d", stats.WildcardRoutes)
	}

	if stats.NamedRoutes != 1 {
		t.Errorf("Expected 1 named route, got %d", stats.NamedRoutes)
	}

	if stats.MiddlewareCount != 1 {
		t.Errorf("Expected 1 middleware, got %d", stats.MiddlewareCount)
	}

	if stats.ConstraintCount != 1 {
		t.Errorf("Expected 1 constraint, got %d", stats.ConstraintCount)
	}

	if stats.MethodDistribution["GET"] != 3 {
		t.Errorf("Expected 3 GET routes, got %d", stats.MethodDistribution["GET"])
	}

	if stats.MethodDistribution["POST"] != 1 {
		t.Errorf("Expected 1 POST route, got %d", stats.MethodDistribution["POST"])
	}
}

// Benchmark tests
func BenchmarkRouteUtils_MatchPath(b *testing.B) {
	utils := NewRouteUtils()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.MatchPath("/users/:id", "/users/123")
	}
}

func BenchmarkRouteUtils_ExtractParams(b *testing.B) {
	utils := NewRouteUtils()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.ExtractParams("/users/:id/posts/:postId", "/users/123/posts/456")
	}
}

func BenchmarkMiddlewareChain_Build(b *testing.B) {
	middleware1 := func(c *context.Context) error { return c.Next() }
	middleware2 := func(c *context.Context) error { return c.Next() }
	middleware3 := func(c *context.Context) error { return c.Next() }
	handler := func(c *context.Context) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewMiddlewareChain(middleware1, middleware2, middleware3)
		chain.Build(handler)
	}
}

func BenchmarkConstraintValidator_ValidateParams(b *testing.B) {
	validator := NewConstraintValidator()
	constraints := map[string]Constraint{
		"id":    IsNumeric,
		"email": IsEmail,
	}
	params := map[string]string{
		"id":    "123",
		"email": "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateParams(params, constraints)
	}
}
