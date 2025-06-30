package router

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/aliwert/go-wolf/pkg/context"
)

// RouteUtils provides utility functions for working with routes
type RouteUtils struct{}

// NewRouteUtils creates a new RouteUtils instance
func NewRouteUtils() *RouteUtils {
	return &RouteUtils{}
}

// MatchPath checks if a given path matches a route pattern
func (ru *RouteUtils) MatchPath(pattern, path string) bool {
	if pattern == path {
		return true
	}

	// Handle parameter matching
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			// Parameter - matches any non-empty value
			if pathParts[i] == "" {
				return false
			}
			continue
		}

		if strings.HasPrefix(part, "*") {
			// Wildcard - matches everything from here
			return true
		}

		if part != pathParts[i] {
			return false
		}
	}

	return true
}

// ExtractParams extracts parameters from a path based on a pattern
func (ru *RouteUtils) ExtractParams(pattern, path string) map[string]string {
	params := make(map[string]string)

	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			paramName := part[1:]
			params[paramName] = pathParts[i]
		} else if strings.HasPrefix(part, "*") {
			paramName := part[1:]
			// Join remaining path parts for wildcard
			params[paramName] = strings.Join(pathParts[i:], "/")
			break
		}
	}

	return params
}

// ValidatePath validates if a path is well-formed
func (ru *RouteUtils) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with '/'")
	}

	// Check for invalid characters
	if strings.Contains(path, "//") {
		return fmt.Errorf("path cannot contain consecutive slashes")
	}

	return nil
}

// NormalizePath normalizes a path by removing trailing slashes and cleaning up
func (ru *RouteUtils) NormalizePath(path string) string {
	if path == "" {
		return "/"
	}

	// Add leading slash if missing
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Remove trailing slash unless it's the root path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	// Clean up consecutive slashes
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	return path
}

// CombinePaths combines two paths safely
func (ru *RouteUtils) CombinePaths(base, path string) string {
	base = ru.NormalizePath(base)
	path = ru.NormalizePath(path)

	if base == "/" {
		return path
	}

	if path == "/" {
		return base
	}

	return base + path
}

// IsValidMethod checks if an HTTP method is valid
func (ru *RouteUtils) IsValidMethod(method string) bool {
	validMethods := []string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
	}

	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}

	return false
}

// ParseRoutePattern parses a route pattern and returns information about it
func (ru *RouteUtils) ParseRoutePattern(pattern string) *RoutePattern {
	parts := strings.Split(pattern, "/")
	params := []string{}
	wildcards := []string{}
	staticParts := []string{}

	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, part[1:])
		} else if strings.HasPrefix(part, "*") {
			wildcards = append(wildcards, part[1:])
		} else if part != "" {
			staticParts = append(staticParts, part)
		}
	}

	return &RoutePattern{
		Pattern:     pattern,
		Params:      params,
		Wildcards:   wildcards,
		StaticParts: staticParts,
		HasParams:   len(params) > 0,
		HasWildcard: len(wildcards) > 0,
	}
}

// RoutePattern contains information about a parsed route pattern
type RoutePattern struct {
	Pattern     string
	Params      []string
	Wildcards   []string
	StaticParts []string
	HasParams   bool
	HasWildcard bool
}

// GenerateURL generates a URL from a pattern and parameters
func (ru *RouteUtils) GenerateURL(pattern string, params map[string]string) (string, error) {
	result := pattern

	// Replace parameters
	for key, value := range params {
		paramPattern := ":" + key
		if !strings.Contains(result, paramPattern) {
			return "", fmt.Errorf("parameter '%s' not found in pattern", key)
		}

		// URL encode the parameter value
		encodedValue := url.QueryEscape(value)
		result = strings.ReplaceAll(result, paramPattern, encodedValue)
	}

	// Check if there are any unreplaced parameters
	if strings.Contains(result, ":") {
		re := regexp.MustCompile(`:(\w+)`)
		matches := re.FindAllStringSubmatch(result, -1)
		if len(matches) > 0 {
			return "", fmt.Errorf("missing parameter: %s", matches[0][1])
		}
	}

	return result, nil
}

// RouteDebugInfo provides debugging information about routes
type RouteDebugInfo struct {
	Method      string
	Pattern     string
	Parameters  []string
	Wildcards   []string
	Middleware  int
	Constraints int
	Named       bool
	Name        string
}

// GetRouteDebugInfo extracts debug information from a RouteInfo
func (ru *RouteUtils) GetRouteDebugInfo(route *RouteInfo) *RouteDebugInfo {
	pattern := ru.ParseRoutePattern(route.Path)

	return &RouteDebugInfo{
		Method:      route.Method,
		Pattern:     route.Path,
		Parameters:  pattern.Params,
		Wildcards:   pattern.Wildcards,
		Middleware:  len(route.Middleware),
		Constraints: len(route.Constraints),
		Named:       route.Name != "",
		Name:        route.Name,
	}
}

// MiddlewareChain represents a chain of middleware functions
type MiddlewareChain struct {
	middleware []context.HandlerFunc
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain(middleware ...context.HandlerFunc) *MiddlewareChain {
	return &MiddlewareChain{
		middleware: middleware,
	}
}

// Add adds middleware to the chain
func (mc *MiddlewareChain) Add(middleware ...context.HandlerFunc) *MiddlewareChain {
	mc.middleware = append(mc.middleware, middleware...)
	return mc
}

// Build builds the final handler with all middleware applied
func (mc *MiddlewareChain) Build(handler context.HandlerFunc) context.HandlerFunc {
	result := handler

	// Apply middleware in reverse order
	for i := len(mc.middleware) - 1; i >= 0; i-- {
		mw := mc.middleware[i]
		next := result
		result = func(mw context.HandlerFunc, next context.HandlerFunc) context.HandlerFunc {
			return func(c *context.Context) error {
				c.SetNext(next)
				return mw(c)
			}
		}(mw, next)
	}

	return result
}

// Length returns the number of middleware in the chain
func (mc *MiddlewareChain) Length() int {
	return len(mc.middleware)
}

// PathMatcher provides advanced path matching functionality
type PathMatcher struct {
	caseSensitive bool
	strictSlash   bool
}

// NewPathMatcher creates a new PathMatcher
func NewPathMatcher(caseSensitive, strictSlash bool) *PathMatcher {
	return &PathMatcher{
		caseSensitive: caseSensitive,
		strictSlash:   strictSlash,
	}
}

// Match matches a path against a pattern with the configured options
func (pm *PathMatcher) Match(pattern, path string) bool {
	if !pm.caseSensitive {
		pattern = strings.ToLower(pattern)
		path = strings.ToLower(path)
	}

	if !pm.strictSlash {
		pattern = strings.TrimSuffix(pattern, "/")
		path = strings.TrimSuffix(path, "/")

		// Handle root path
		if pattern == "" {
			pattern = "/"
		}
		if path == "" {
			path = "/"
		}
	}

	utils := NewRouteUtils()
	return utils.MatchPath(pattern, path)
}

// RouteConflictDetector detects conflicts between routes
type RouteConflictDetector struct{}

// NewRouteConflictDetector creates a new conflict detector
func NewRouteConflictDetector() *RouteConflictDetector {
	return &RouteConflictDetector{}
}

// DetectConflicts detects conflicts between two route patterns
func (rcd *RouteConflictDetector) DetectConflicts(pattern1, pattern2 string) bool {
	if pattern1 == pattern2 {
		return true
	}

	utils := NewRouteUtils()
	p1 := utils.ParseRoutePattern(pattern1)
	p2 := utils.ParseRoutePattern(pattern2)

	// If both have wildcards at the same position, they conflict
	if p1.HasWildcard && p2.HasWildcard {
		return true
	}

	// Check if patterns can match the same paths
	parts1 := strings.Split(pattern1, "/")
	parts2 := strings.Split(pattern2, "/")

	if len(parts1) != len(parts2) {
		return false
	}

	for i := 0; i < len(parts1); i++ {
		part1 := parts1[i]
		part2 := parts2[i]

		// If both are static and different, no conflict
		if !strings.HasPrefix(part1, ":") && !strings.HasPrefix(part1, "*") &&
			!strings.HasPrefix(part2, ":") && !strings.HasPrefix(part2, "*") &&
			part1 != part2 {
			return false
		}
	}

	return true
}

// ConstraintValidator provides validation utilities for route constraints
type ConstraintValidator struct{}

// NewConstraintValidator creates a new constraint validator
func NewConstraintValidator() *ConstraintValidator {
	return &ConstraintValidator{}
}

// ValidateParams validates parameters against their constraints
func (cv *ConstraintValidator) ValidateParams(params map[string]string, constraints map[string]Constraint) error {
	for paramName, constraint := range constraints {
		value, exists := params[paramName]
		if !exists {
			return fmt.Errorf("parameter '%s' is required", paramName)
		}

		if !constraint(value) {
			return fmt.Errorf("parameter '%s' with value '%s' failed validation", paramName, value)
		}
	}

	return nil
}

// RouteStatistics provides statistics about routes
type RouteStatistics struct {
	TotalRoutes        int
	StaticRoutes       int
	ParametricRoutes   int
	WildcardRoutes     int
	NamedRoutes        int
	MiddlewareCount    int
	ConstraintCount    int
	MethodDistribution map[string]int
}

// CalculateStatistics calculates statistics for a set of routes
func CalculateStatistics(routes []*RouteInfo) *RouteStatistics {
	stats := &RouteStatistics{
		MethodDistribution: make(map[string]int),
	}

	utils := NewRouteUtils()

	for _, route := range routes {
		stats.TotalRoutes++

		if route.Name != "" {
			stats.NamedRoutes++
		}

		stats.MiddlewareCount += len(route.Middleware)
		stats.ConstraintCount += len(route.Constraints)

		stats.MethodDistribution[route.Method]++

		pattern := utils.ParseRoutePattern(route.Path)
		if pattern.HasWildcard {
			stats.WildcardRoutes++
		} else if pattern.HasParams {
			stats.ParametricRoutes++
		} else {
			stats.StaticRoutes++
		}
	}

	return stats
}

// String returns a string representation of the statistics
func (rs *RouteStatistics) String() string {
	return fmt.Sprintf(
		"Routes: %d total (%d static, %d parametric, %d wildcard), %d named, %d middleware, %d constraints",
		rs.TotalRoutes,
		rs.StaticRoutes,
		rs.ParametricRoutes,
		rs.WildcardRoutes,
		rs.NamedRoutes,
		rs.MiddlewareCount,
		rs.ConstraintCount,
	)
}
