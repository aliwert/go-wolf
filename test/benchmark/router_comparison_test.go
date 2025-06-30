package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/aliwert/go-wolf"
	"github.com/aliwert/go-wolf/pkg/context"
)

// BenchmarkComparison provides benchmarks to compare go-wolf against standard library
// and other popular routers (when available)

// BenchmarkGoWolfVsStdLib compares go-wolf router performance against standard library
func BenchmarkGoWolfVsStdLib(b *testing.B) {
	// go-wolf setup
	wolfApp := wolf.New()
	wolfApp.GET("/", func(c *context.Context) error {
		return c.String(200, "hello world")
	})
	wolfApp.GET("/users/:id", func(c *context.Context) error {
		return c.String(200, c.Param("id"))
	})
	wolfApp.GET("/static/*filepath", func(c *context.Context) error {
		return c.String(200, c.Param("filepath"))
	})

	// Standard library setup
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("user"))
	})
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("static"))
	})

	tests := []struct {
		name string
		path string
	}{
		{"StaticRoute", "/"},
		{"UserRoute", "/users/123"},
		{"StaticFileRoute", "/static/css/main.css"},
	}

	for _, tt := range tests {
		b.Run("GoWolf_"+tt.name, func(b *testing.B) {
			req := httptest.NewRequest("GET", tt.path, nil)
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				wolfApp.ServeHTTP(resp, req)
			}
		})

		b.Run("StdLib_"+tt.name, func(b *testing.B) {
			req := httptest.NewRequest("GET", tt.path, nil)
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				mux.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkRouterScaling tests how performance scales with route count
func BenchmarkRouterScaling(b *testing.B) {
	routeCounts := []int{1, 10, 100, 1000, 5000}

	for _, routeCount := range routeCounts {
		b.Run(fmt.Sprintf("Routes_%d", routeCount), func(b *testing.B) {
			app := wolf.New()

			// Add routes
			for i := 0; i < routeCount; i++ {
				path := fmt.Sprintf("/route_%d", i)
				app.GET(path, func(c *context.Context) error {
					return c.String(200, "ok")
				})
			}

			// Test the last route (worst case scenario)
			testPath := fmt.Sprintf("/route_%d", routeCount-1)
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

// BenchmarkParameterExtraction tests parameter extraction performance
func BenchmarkParameterExtraction(b *testing.B) {
	tests := []struct {
		name       string
		pattern    string
		testPath   string
		paramCount int
	}{
		{"NoParams", "/static/route", "/static/route", 0},
		{"OneParam", "/users/:id", "/users/12345", 1},
		{"TwoParams", "/users/:userId/posts/:postId", "/users/12345/posts/67890", 2},
		{"FiveParams", "/a/:b/c/:d/e/:f/g/:h/i/:j", "/a/1/c/2/e/3/g/4/i/5", 5},
		{"Wildcard", "/files/*path", "/files/documents/report.pdf", 1},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			app := wolf.New()
			app.GET(tt.pattern, func(c *context.Context) error {
				// Extract all parameters to measure extraction cost
				for i := 0; i < tt.paramCount; i++ {
					_ = c.Param(fmt.Sprintf("param%d", i))
				}
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

// BenchmarkMiddlewareOverhead tests middleware performance impact
func BenchmarkMiddlewareOverhead(b *testing.B) {
	middlewareCounts := []int{0, 1, 5, 10, 20}

	for _, count := range middlewareCounts {
		b.Run(fmt.Sprintf("Middleware_%d", count), func(b *testing.B) {
			app := wolf.New()

			// Add middleware
			for i := 0; i < count; i++ {
				app.Use(func(c *context.Context) error {
					// Simple middleware that adds a header
					c.SetHeader(fmt.Sprintf("X-Middleware-%d", i), "true")
					return c.Next()
				})
			}

			app.GET("/", func(c *context.Context) error {
				return c.String(200, "ok")
			})

			req := httptest.NewRequest("GET", "/", nil)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				resp := httptest.NewRecorder()
				app.ServeHTTP(resp, req)
			}
		})
	}
}

// BenchmarkMemoryFootprint measures memory allocations per request
func BenchmarkMemoryFootprint(b *testing.B) {
	app := wolf.New()

	// Add a realistic route with multiple parameters
	app.GET("/api/v1/users/:userId/posts/:postId/comments/:commentId", func(c *context.Context) error {
		userId := c.Param("userId")
		postId := c.Param("postId")
		commentId := c.Param("commentId")

		response := map[string]string{
			"userId":    userId,
			"postId":    postId,
			"commentId": commentId,
		}

		return c.JSON(200, response)
	})

	req := httptest.NewRequest("GET", "/api/v1/users/123/posts/456/comments/789", nil)

	// Force GC before benchmark
	runtime.GC()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}

// BenchmarkConcurrentLoad tests performance under concurrent load
func BenchmarkConcurrentLoad(b *testing.B) {
	app := wolf.New()

	// Add multiple routes
	app.GET("/", func(c *context.Context) error { return c.String(200, "home") })
	app.GET("/api/users/:id", func(c *context.Context) error { return c.String(200, c.Param("id")) })
	app.GET("/static/*file", func(c *context.Context) error { return c.String(200, c.Param("file")) })

	routes := []string{
		"/",
		"/api/users/123",
		"/static/css/main.css",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			route := routes[counter%len(routes)]
			req := httptest.NewRequest("GET", route, nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)
			counter++
		}
	})
}

// BenchmarkLongestPath tests performance with very long paths
func BenchmarkLongestPath(b *testing.B) {
	app := wolf.New()

	// Create a very long path
	longPath := "/very/long/path/with/many/segments/that/could/potentially/slow/down/routing/performance/test/case"
	app.GET(longPath, func(c *context.Context) error {
		return c.String(200, "found")
	})

	req := httptest.NewRequest("GET", longPath, nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}

// BenchmarkWorstCase tests the worst-case routing scenario
func BenchmarkWorstCase(b *testing.B) {
	app := wolf.New()

	// Add many similar routes that would cause maximum tree traversal
	for i := 0; i < 1000; i++ {
		path := fmt.Sprintf("/api/v1/resource_%04d", i)
		app.GET(path, func(c *context.Context) error {
			return c.String(200, "ok")
		})
	}

	// Test the last route (requires traversing most of the tree)
	req := httptest.NewRequest("GET", "/api/v1/resource_0999", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
	}
}
