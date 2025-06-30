package tests

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aliwert/go-wolf"
	"github.com/aliwert/go-wolf/pkg/context"
)

// BenchmarkRouterSuite contains comprehensive router performance tests
// These benchmarks measure the performance of different routing scenarios

// BenchmarkStaticRoutes tests performance of static route matching
func BenchmarkStaticRoutes(b *testing.B) {
	tests := []struct {
		name      string
		numRoutes int
	}{
		{"10Routes", 10},
		{"100Routes", 100},
		{"1000Routes", 1000},
		{"10000Routes", 10000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			app := wolf.New()

			// Add many static routes
			for i := 0; i < tt.numRoutes; i++ {
				path := fmt.Sprintf("/route%d", i)
				app.GET(path, func(c *context.Context) error {
					return c.String(200, "ok")
				})
			}

			// Test the middle route to ensure we're not just hitting the first one
			testPath := fmt.Sprintf("/route%d", tt.numRoutes/2)
			req := httptest.NewRequest("GET", testPath, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkParametricRoutes tests performance of parametric route matching
func BenchmarkParametricRoutes(b *testing.B) {
	tests := []struct {
		name     string
		pattern  string
		testPath string
	}{
		{"SingleParam", "/users/:id", "/users/123"},
		{"TwoParams", "/users/:id/posts/:postId", "/users/123/posts/456"},
		{"ThreeParams", "/api/:version/users/:id/posts/:postId", "/api/v1/users/123/posts/456"},
		{"DeepNested", "/a/:b/c/:d/e/:f/g/:h", "/a/1/c/2/e/3/g/4"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			app := wolf.New()
			app.GET(tt.pattern, func(c *context.Context) error {
				return c.String(200, "ok")
			})

			req := httptest.NewRequest("GET", tt.testPath, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkWildcardRoutes tests performance of wildcard route matching
func BenchmarkWildcardRoutes(b *testing.B) {
	tests := []struct {
		name     string
		pattern  string
		testPath string
	}{
		{"ShortPath", "/static/*filepath", "/static/file.css"},
		{"MediumPath", "/static/*filepath", "/static/css/components/button.css"},
		{"LongPath", "/static/*filepath", "/static/assets/images/icons/social/facebook/large.png"},
		{"VeryLongPath", "/static/*filepath", "/static/vendor/library/dist/assets/css/components/ui/forms/inputs/special/date-picker.min.css"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			app := wolf.New()
			app.GET(tt.pattern, func(c *context.Context) error {
				return c.String(200, "ok")
			})

			req := httptest.NewRequest("GET", tt.testPath, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkMixedRoutes tests performance with a realistic mix of route types
func BenchmarkMixedRoutes(b *testing.B) {
	app := wolf.New()

	// Add static routes
	app.GET("/", func(c *context.Context) error { return c.String(200, "home") })
	app.GET("/about", func(c *context.Context) error { return c.String(200, "about") })
	app.GET("/contact", func(c *context.Context) error { return c.String(200, "contact") })
	app.GET("/pricing", func(c *context.Context) error { return c.String(200, "pricing") })
	app.GET("/features", func(c *context.Context) error { return c.String(200, "features") })

	// Add parametric routes
	app.GET("/users/:id", func(c *context.Context) error { return c.String(200, "user") })
	app.GET("/users/:id/profile", func(c *context.Context) error { return c.String(200, "profile") })
	app.GET("/users/:id/posts/:postId", func(c *context.Context) error { return c.String(200, "post") })
	app.GET("/api/v1/users/:id", func(c *context.Context) error { return c.String(200, "api user") })
	app.GET("/api/v1/posts/:id/comments/:commentId", func(c *context.Context) error { return c.String(200, "comment") })

	// Add wildcard routes
	app.GET("/static/*filepath", func(c *context.Context) error { return c.String(200, "static") })
	app.GET("/uploads/*filepath", func(c *context.Context) error { return c.String(200, "upload") })

	// Test routes to benchmark
	testRoutes := []string{
		"/",
		"/about",
		"/users/123",
		"/users/456/posts/789",
		"/api/v1/users/123",
		"/static/css/main.css",
		"/uploads/images/photo.jpg",
	}

	for i, testPath := range testRoutes {
		b.Run(fmt.Sprintf("Route%d_%s", i+1, testPath[1:]), func(b *testing.B) {
			req := httptest.NewRequest("GET", testPath, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for j := 0; j < b.N; j++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkRouteConflicts tests performance when routes have conflicts
func BenchmarkRouteConflicts(b *testing.B) {
	app := wolf.New()

	// Add conflicting routes (router should handle them correctly)
	app.GET("/users/profile", func(c *context.Context) error { return c.String(200, "static profile") })
	app.GET("/users/:id", func(c *context.Context) error { return c.String(200, "user by id") })
	app.GET("/users/:id/posts", func(c *context.Context) error { return c.String(200, "user posts") })
	app.GET("/users/:id/posts/recent", func(c *context.Context) error { return c.String(200, "recent posts") })

	tests := []struct {
		name string
		path string
	}{
		{"StaticMatch", "/users/profile"},
		{"ParamMatch", "/users/123"},
		{"ParamWithSuffix", "/users/123/posts"},
		{"StaticSuffix", "/users/123/posts/recent"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			req := httptest.NewRequest("GET", tt.path, nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkHTTPMethods tests performance across different HTTP methods
func BenchmarkHTTPMethods(b *testing.B) {
	app := wolf.New()

	// Add routes for different methods
	app.GET("/resource", func(c *context.Context) error { return c.String(200, "get") })
	app.POST("/resource", func(c *context.Context) error { return c.String(201, "post") })
	app.PUT("/resource", func(c *context.Context) error { return c.String(200, "put") })
	app.PATCH("/resource", func(c *context.Context) error { return c.String(200, "patch") })
	app.DELETE("/resource", func(c *context.Context) error { return c.String(204, "delete") })

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		b.Run(method, func(b *testing.B) {
			req := httptest.NewRequest(method, "/resource", nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkMemoryAllocation focuses on memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	app := wolf.New()
	app.GET("/users/:id/posts/:postId/comments/:commentId", func(c *context.Context) error {
		// Access parameters to ensure they're parsed
		id := c.Param("id")
		postId := c.Param("postId")
		commentId := c.Param("commentId")

		// Use the parameters to prevent optimization
		result := fmt.Sprintf("%s-%s-%s", id, postId, commentId)
		return c.String(200, result)
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456/comments/789", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}

// BenchmarkConcurrentRouting tests performance under concurrent load
func BenchmarkConcurrentRouting(b *testing.B) {
	app := wolf.New()

	// Add various routes
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/route%d", i)
		app.GET(path, func(c *context.Context) error {
			return c.String(200, "ok")
		})
	}

	// Test concurrent access to different routes
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			path := fmt.Sprintf("/route%d", counter%100)
			req := httptest.NewRequest("GET", path, nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)
			counter++
		}
	})
}

// BenchmarkLargePathParameters tests performance with large parameter values
func BenchmarkLargePathParameters(b *testing.B) {
	app := wolf.New()
	app.GET("/data/:id", func(c *context.Context) error {
		id := c.Param("id")
		return c.String(200, id)
	})

	// Generate large parameter value
	largeId := strings.Repeat("a", 1000)
	req := httptest.NewRequest("GET", "/data/"+largeId, nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}

// BenchmarkDeepNesting tests performance with deeply nested routes
func BenchmarkDeepNesting(b *testing.B) {
	app := wolf.New()

	// Create deeply nested route structure
	depth := 20
	path := "/"
	for i := 0; i < depth; i++ {
		path += fmt.Sprintf("level%d/", i)
	}
	path += ":id"

	app.GET(path, func(c *context.Context) error {
		return c.String(200, c.Param("id"))
	})

	// Create test path
	testPath := "/"
	for i := 0; i < depth; i++ {
		testPath += fmt.Sprintf("level%d/", i)
	}
	testPath += "finalvalue"

	req := httptest.NewRequest("GET", testPath, nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}
