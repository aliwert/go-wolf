package request

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Test structs
type User struct {
	Name     string `json:"name" form:"name" query:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email"`
	Age      int    `json:"age" form:"age" query:"age" validate:"min=0,max=120"`
	Active   bool   `json:"active" form:"active" query:"active"`
	Website  string `json:"website" form:"website" query:"website" validate:"url"`
	Username string `json:"username" form:"username" query:"username" validate:"required,alphanumeric,min=3"`
}

func TestBindJSON(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		expected    User
	}{
		{
			name: "valid JSON",
			body: `{"name":"John","email":"john@example.com","age":30,"active":true,"username":"john123"}`,
			expected: User{
				Name:     "John",
				Email:    "john@example.com",
				Age:      30,
				Active:   true,
				Username: "john123",
			},
		},
		{
			name:        "invalid JSON",
			body:        `{"name":"John","email":}`,
			expectError: true,
		},
		{
			name:        "validation error",
			body:        `{"name":"","email":"invalid","age":-1}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			var user User
			err := BindJSON(req, &user)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if user.Name != tt.expected.Name {
				t.Errorf("expected name %s, got %s", tt.expected.Name, user.Name)
			}
			if user.Email != tt.expected.Email {
				t.Errorf("expected email %s, got %s", tt.expected.Email, user.Email)
			}
			if user.Age != tt.expected.Age {
				t.Errorf("expected age %d, got %d", tt.expected.Age, user.Age)
			}
			if user.Active != tt.expected.Active {
				t.Errorf("expected active %v, got %v", tt.expected.Active, user.Active)
			}
		})
	}
}

func TestBindQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectError bool
		expected    User
	}{
		{
			name:  "valid query parameters",
			query: "name=John&email=john@example.com&age=30&active=true&username=john123",
			expected: User{
				Name:     "John",
				Email:    "john@example.com",
				Age:      30,
				Active:   true,
				Username: "john123",
			},
		},
		{
			name:        "validation error",
			query:       "name=&email=invalid&age=-1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?"+tt.query, nil)

			var user User
			err := BindQuery(req, &user)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if user.Name != tt.expected.Name {
				t.Errorf("expected name %s, got %s", tt.expected.Name, user.Name)
			}
		})
	}
}

func TestBindForm(t *testing.T) {
	form := url.Values{}
	form.Add("name", "John")
	form.Add("email", "john@example.com")
	form.Add("age", "30")
	form.Add("active", "true")
	form.Add("username", "john123")

	req := httptest.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var user User
	err := BindForm(req, &user)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Name != "John" {
		t.Errorf("expected name John, got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", user.Email)
	}
	if user.Age != 30 {
		t.Errorf("expected age 30, got %d", user.Age)
	}
	if !user.Active {
		t.Error("expected active to be true")
	}
	if user.Username != "john123" {
		t.Errorf("expected username john123, got %s", user.Username)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		user        User
		expectError bool
		errorCount  int
	}{
		{
			name: "valid user",
			user: User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				Website:  "https://example.com",
				Username: "johndoe123",
			},
			expectError: false,
		},
		{
			name: "missing required fields",
			user: User{
				Age: 30,
			},
			expectError: true,
			errorCount:  3, // name, email, username required
		},
		{
			name: "invalid email",
			user: User{
				Name:     "John",
				Email:    "invalid-email",
				Username: "johndoe",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "invalid age",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Age:      -1,
				Username: "johndoe",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "invalid URL",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Age:      30,
				Website:  "not-a-url",
				Username: "johndoe",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "non-alphanumeric username",
			user: User{
				Name:     "John",
				Email:    "john@example.com",
				Age:      30,
				Username: "john@doe",
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("expected validation error but got none")
					return
				}

				if ve, ok := err.(ValidationErrors); ok {
					if len(ve) != tt.errorCount {
						t.Errorf("expected %d validation errors, got %d: %v", tt.errorCount, len(ve), ve)
					}
				} else {
					t.Errorf("expected ValidationErrors, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestRequestWrapper(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/test?name=John&age=30", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})

	wrapper := New(req)

	// Test query parameters
	if name := wrapper.QueryParam("name"); name != "John" {
		t.Errorf("expected name John, got %s", name)
	}

	if age, err := wrapper.QueryParamInt("age"); err != nil || age != 30 {
		t.Errorf("expected age 30, got %d (error: %v)", age, err)
	}

	// Test headers
	if ua := wrapper.UserAgent(); ua != "test-agent" {
		t.Errorf("expected user agent test-agent, got %s", ua)
	}

	// Test client IP
	if ip := wrapper.ClientIP(); ip != "192.168.1.1" {
		t.Errorf("expected client IP 192.168.1.1, got %s", ip)
	}

	// Test cookies
	if sessionValue, err := wrapper.CookieValue("session"); err != nil || sessionValue != "abc123" {
		t.Errorf("expected session abc123, got %s (error: %v)", sessionValue, err)
	}
}

func TestContentTypeDetection(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		isJSON      bool
		isForm      bool
		isXML       bool
	}{
		{
			name:        "JSON content type",
			contentType: "application/json",
			isJSON:      true,
		},
		{
			name:        "Form content type",
			contentType: "application/x-www-form-urlencoded",
			isForm:      true,
		},
		{
			name:        "XML content type",
			contentType: "application/xml",
			isXML:       true,
		},
		{
			name:        "Multipart form",
			contentType: "multipart/form-data",
			isForm:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", nil)
			req.Header.Set("Content-Type", tt.contentType)

			if IsJSON(req) != tt.isJSON {
				t.Errorf("IsJSON() = %v, want %v", IsJSON(req), tt.isJSON)
			}
			if IsForm(req) != tt.isForm {
				t.Errorf("IsForm() = %v, want %v", IsForm(req), tt.isForm)
			}
			if IsXML(req) != tt.isXML {
				t.Errorf("IsXML() = %v, want %v", IsXML(req), tt.isXML)
			}
		})
	}
}

func TestGetAuth(t *testing.T) {
	tests := []struct {
		name          string
		setupRequest  func() *http.Request
		expectedType  string
		expectedValid bool
		expectedToken string
	}{
		{
			name: "Bearer token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer abc123")
				return req
			},
			expectedType:  "bearer",
			expectedValid: true,
			expectedToken: "abc123",
		},
		{
			name: "JWT token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
				return req
			},
			expectedType:  "jwt",
			expectedValid: true,
			expectedToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
		},
		{
			name: "API Key in header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-API-Key", "api-key-123")
				return req
			},
			expectedType:  "api-key",
			expectedValid: true,
			expectedToken: "api-key-123",
		},
		{
			name: "API Key in query",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test?api_key=query-key-456", nil)
				return req
			},
			expectedType:  "api-key",
			expectedValid: true,
			expectedToken: "query-key-456",
		},
		{
			name: "Basic auth",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Basic dXNlcjpwYXNz") // user:pass in base64
				return req
			},
			expectedType:  "basic",
			expectedValid: true,
		},
		{
			name: "No auth",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			expectedType:  "",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := New(tt.setupRequest())
			auth := req.GetAuth()

			if auth.Type != tt.expectedType {
				t.Errorf("expected auth type %s, got %s", tt.expectedType, auth.Type)
			}

			if auth.Valid != tt.expectedValid {
				t.Errorf("expected auth valid %t, got %t", tt.expectedValid, auth.Valid)
			}

			if tt.expectedToken != "" && auth.Token != tt.expectedToken {
				t.Errorf("expected token %s, got %s", tt.expectedToken, auth.Token)
			}
		})
	}
}

func TestAuthHelpers(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer abc123")

	r := New(req)

	if !r.HasAuth() {
		t.Error("expected HasAuth to return true")
	}

	if r.AuthType() != "bearer" {
		t.Errorf("expected auth type bearer, got %s", r.AuthType())
	}

	if r.AuthToken() != "abc123" {
		t.Errorf("expected token abc123, got %s", r.AuthToken())
	}

	if !r.IsAuthType("bearer") {
		t.Error("expected IsAuthType('bearer') to return true")
	}

	if r.IsAuthType("basic") {
		t.Error("expected IsAuthType('basic') to return false")
	}
}
