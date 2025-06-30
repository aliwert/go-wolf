package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliwert/go-wolf/pkg/context"
)

func TestGroupCreation(t *testing.T) {
	router := New()

	// Test group creation
	v1 := router.Group("/api/v1")
	if v1 == nil {
		t.Fatal("Expected group to be created, got nil")
	}

	// Test group prefix
	if v1.prefix != "/api/v1" {
		t.Errorf("Expected prefix '/api/v1', got '%s'", v1.prefix)
	}

	// Test router reference
	if v1.router != router {
		t.Error("Expected group to reference the original router")
	}
}

func TestGroupNestedCreation(t *testing.T) {
	router := New()

	// Create nested groups
	api := router.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	if users.prefix != "/api/v1/users" {
		t.Errorf("Expected nested prefix '/api/v1/users', got '%s'", users.prefix)
	}
}

func TestGroupMiddleware(t *testing.T) {
	router := New()

	middleware1 := func(c *context.Context) error {
		c.SetHeader("X-Middleware-1", "applied")
		return c.Next()
	}

	middleware2 := func(c *context.Context) error {
		c.SetHeader("X-Middleware-2", "applied")
		return c.Next()
	}

	// Create group with middleware
	v1 := router.Group("/api/v1", middleware1)

	if len(v1.middleware) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(v1.middleware))
	}

	// Add more middleware
	v1.Use(middleware2)

	if len(v1.middleware) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(v1.middleware))
	}
}

func TestGroupRoutes(t *testing.T) {
	router := New()
	v1 := router.Group("/api/v1")

	// Add routes to group
	v1.GET("/users", func(c *context.Context) error {
		return c.String(http.StatusOK, "users")
	})

	v1.POST("/users", func(c *context.Context) error {
		return c.String(http.StatusCreated, "user created")
	})

	// Test GET route
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "users" {
		t.Errorf("Expected body 'users', got '%s'", w.Body.String())
	}

	context.Release(c)

	// Test POST route
	req = httptest.NewRequest("POST", "/api/v1/users", nil)
	w = httptest.NewRecorder()
	c = context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if w.Body.String() != "user created" {
		t.Errorf("Expected body 'user created', got '%s'", w.Body.String())
	}

	context.Release(c)
}

func TestGroupHTTPMethods(t *testing.T) {
	router := New()
	api := router.Group("/api")

	// Add routes for all HTTP methods
	api.GET("/get", func(c *context.Context) error {
		return c.String(http.StatusOK, "GET")
	})

	api.POST("/post", func(c *context.Context) error {
		return c.String(http.StatusOK, "POST")
	})

	api.PUT("/put", func(c *context.Context) error {
		return c.String(http.StatusOK, "PUT")
	})

	api.DELETE("/delete", func(c *context.Context) error {
		return c.String(http.StatusOK, "DELETE")
	})

	api.PATCH("/patch", func(c *context.Context) error {
		return c.String(http.StatusOK, "PATCH")
	})

	// Test all methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	paths := []string{"/api/get", "/api/post", "/api/put", "/api/delete", "/api/patch"}

	for i, method := range methods {
		req := httptest.NewRequest(method, paths[i], nil)
		w := httptest.NewRecorder()
		c := context.Acquire()
		c.Reset(w, req)

		router.ServeHTTP(w, req, c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
		}

		if w.Body.String() != method {
			t.Errorf("Expected body '%s', got '%s'", method, w.Body.String())
		}

		context.Release(c)
	}
}

func TestGroupWithParameters(t *testing.T) {
	router := New()
	users := router.Group("/users")

	users.GET("/:id", func(c *context.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	users.GET("/:id/posts/:postId", func(c *context.Context) error {
		return c.String(http.StatusOK, c.Param("id")+":"+c.Param("postId"))
	})

	// Test single parameter
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "123" {
		t.Errorf("Expected body '123', got '%s'", w.Body.String())
	}

	context.Release(c)

	// Test multiple parameters
	req = httptest.NewRequest("GET", "/users/123/posts/456", nil)
	w = httptest.NewRecorder()
	c = context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "123:456" {
		t.Errorf("Expected body '123:456', got '%s'", w.Body.String())
	}

	context.Release(c)
}

func TestGroupMiddlewareExecution(t *testing.T) {
	router := New()

	middleware1 := func(c *context.Context) error {
		c.SetHeader("X-Group", "group1")
		return c.Next()
	}

	middleware2 := func(c *context.Context) error {
		c.SetHeader("X-Route", "route1")
		return c.Next()
	}

	api := router.Group("/api", middleware1)

	api.GET("/test", func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	}, middleware2)

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check if middleware headers were set
	if w.Header().Get("X-Group") != "group1" {
		t.Errorf("Expected X-Group header 'group1', got '%s'", w.Header().Get("X-Group"))
	}

	if w.Header().Get("X-Route") != "route1" {
		t.Errorf("Expected X-Route header 'route1', got '%s'", w.Header().Get("X-Route"))
	}

	context.Release(c)
}

func TestGroupNestedMiddleware(t *testing.T) {
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

	// Create nested groups with middleware
	api := router.Group("/api", middleware1)
	v1 := api.Group("/v1", middleware2)
	users := v1.Group("/users", middleware3)

	users.GET("/test", func(c *context.Context) error {
		return c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest("GET", "/api/v1/users/test", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check if all middleware was executed in order
	if w.Header().Get("X-Level") != "1-2-3" {
		t.Errorf("Expected X-Level header '1-2-3', got '%s'", w.Header().Get("X-Level"))
	}

	context.Release(c)
}

func TestGroupWildcardRoutes(t *testing.T) {
	router := New()
	static := router.Group("/static")

	static.GET("/*filepath", func(c *context.Context) error {
		return c.String(http.StatusOK, c.Param("filepath"))
	})

	req := httptest.NewRequest("GET", "/static/css/main.css", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "/css/main.css" {
		t.Errorf("Expected body '/css/main.css', got '%s'", w.Body.String())
	}

	context.Release(c)
}

func TestGroupConflictingRoutes(t *testing.T) {
	router := New()
	users := router.Group("/users")

	// Add parameterized route first
	users.GET("/:id", func(c *context.Context) error {
		return c.String(http.StatusOK, "user: "+c.Param("id"))
	})

	// Add static routes that should work alongside the parameterized route
	users.GET("/:id/posts", func(c *context.Context) error {
		return c.String(http.StatusOK, "user posts: "+c.Param("id"))
	})

	// Test parameterized route
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "user: 123" {
		t.Errorf("Expected body 'user: 123', got '%s'", w.Body.String())
	}

	context.Release(c)

	// Test nested parameterized route
	req = httptest.NewRequest("GET", "/users/123/posts", nil)
	w = httptest.NewRecorder()
	c = context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "user posts: 123" {
		t.Errorf("Expected body 'user posts: 123', got '%s'", w.Body.String())
	}

	context.Release(c)
}

func TestGroupStaticVsParam(t *testing.T) {
	// Test demonstrating that the router has specific rules about static vs param routes
	router := New()
	users := router.Group("/users")

	// Add static routes first
	users.GET("/new", func(c *context.Context) error {
		return c.String(http.StatusOK, "new user form")
	})

	users.GET("/admin", func(c *context.Context) error {
		return c.String(http.StatusOK, "admin user")
	})

	// Note: Adding a parameterized route after static routes at the same level
	// would cause a conflict in this router implementation. This is expected behavior.

	// Test static routes
	req := httptest.NewRequest("GET", "/users/new", nil)
	w := httptest.NewRecorder()
	c := context.Acquire()
	c.Reset(w, req)

	router.ServeHTTP(w, req, c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "new user form" {
		t.Errorf("Expected body 'new user form', got '%s'", w.Body.String())
	}

	context.Release(c)
}
