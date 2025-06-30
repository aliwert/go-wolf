package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliwert/go-wolf/pkg/context"
	"github.com/stretchr/testify/assert"
)

func simpleHandler(message string) context.HandlerFunc {
	return func(c *context.Context) error {
		return c.String(http.StatusOK, message)
	}
}

func paramHandler(c *context.Context) error {
	return c.String(http.StatusOK, c.Param("id"))
}

func testMiddleware(id string) context.HandlerFunc {
	return func(c *context.Context) error {
		c.SetHeader("X-Middleware", c.Header("X-Middleware")+id)
		return c.Next()
	}
}

// --- Test Cases ---

func TestRouter_BasicRouting(t *testing.T) {
	router := New()
	router.Handle("GET", "/", simpleHandler("root"))
	router.Handle("GET", "/users/:id", paramHandler)

	t.Run("BasicRouting", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		c := context.Acquire()
		defer context.Release(c)
		c.Reset(resp, req)

		router.ServeHTTP(resp, req, c)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "root", resp.Body.String())
	})

	t.Run("ParameterRouting", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/123", nil)
		resp := httptest.NewRecorder()
		c := context.Acquire()
		defer context.Release(c)
		c.Reset(resp, req)

		router.ServeHTTP(resp, req, c)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "123", resp.Body.String())
	})
}

func TestRouter_Middleware(t *testing.T) {
	router := New()
	router.Handle("GET", "/", simpleHandler("handler"), testMiddleware("route"))

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	c := context.Acquire()
	defer context.Release(c)
	c.Reset(resp, req)

	router.ServeHTTP(resp, req, c)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "route", resp.Header().Get("X-Middleware"))
}

func TestRouter_Groups(t *testing.T) {
	router := New()
	adminGroup := router.Group("/admin", testMiddleware("group"))
	adminGroup.GET("/dashboard", simpleHandler("dashboard"))

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	resp := httptest.NewRecorder()
	c := context.Acquire()
	defer context.Release(c)
	c.Reset(resp, req)

	router.ServeHTTP(resp, req, c)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "dashboard", resp.Body.String())
	assert.Equal(t, "group", resp.Header().Get("X-Middleware"))
}

func TestRouter_ErrorHandling(t *testing.T) {
	router := New()
	router.Handle("GET", "/exists", simpleHandler("ok"))

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent", nil)
		resp := httptest.NewRecorder()
		c := context.Acquire()
		defer context.Release(c)
		c.Reset(resp, req)

		router.ServeHTTP(resp, req, c)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/exists", nil)
		resp := httptest.NewRecorder()
		c := context.Acquire()
		defer context.Release(c)
		c.Reset(resp, req)

		router.ServeHTTP(resp, req, c)

		assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	})
}

// Benchmark tests
func BenchmarkRouterStaticRoute(b *testing.B) {
	router := New()
	router.Handle("GET", "/test", func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c := context.Acquire()
		c.Reset(w, req)
		router.ServeHTTP(w, req, c)
		context.Release(c)
	}
}

func BenchmarkRouterParameterRoute(b *testing.B) {
	router := New()
	router.Handle("GET", "/users/:id", func(c *context.Context) error {
		return c.String(http.StatusOK, "user")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c := context.Acquire()
		c.Reset(w, req)
		router.ServeHTTP(w, req, c)
		context.Release(c)
	}
}

func BenchmarkRouterWildcardRoute(b *testing.B) {
	router := New()
	router.Handle("GET", "/static/*filepath", func(c *context.Context) error {
		return c.String(http.StatusOK, "file")
	})

	req := httptest.NewRequest(http.MethodGet, "/static/css/main.css", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c := context.Acquire()
		c.Reset(w, req)
		router.ServeHTTP(w, req, c)
		context.Release(c)
	}
}
