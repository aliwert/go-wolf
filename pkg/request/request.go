// Package request provides utilities for binding, validating, and processing HTTP request data
package request

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Request wraps http.Request with additional functionality
type Request struct {
	*http.Request
	parsedForm      bool
	parsedMultipart bool
}

// New creates a new Request wrapper
func New(r *http.Request) *Request {
	return &Request{Request: r}
}

// Body returns the request body as bytes
func (r *Request) Body() ([]byte, error) {
	if r.Request.Body == nil {
		return nil, nil
	}
	return io.ReadAll(r.Request.Body)
}

// QueryParam returns a query parameter value
func (r *Request) QueryParam(key string) string {
	return r.URL.Query().Get(key)
}

// QueryParamDefault returns a query parameter value or default
func (r *Request) QueryParamDefault(key, defaultValue string) string {
	if value := r.QueryParam(key); value != "" {
		return value
	}
	return defaultValue
}

// QueryParamInt returns a query parameter as integer
func (r *Request) QueryParamInt(key string) (int, error) {
	value := r.QueryParam(key)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// QueryParamIntDefault returns a query parameter as integer or default
func (r *Request) QueryParamIntDefault(key string, defaultValue int) int {
	value, err := r.QueryParamInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// QueryParamBool returns a query parameter as boolean
func (r *Request) QueryParamBool(key string) (bool, error) {
	value := r.QueryParam(key)
	if value == "" {
		return false, nil
	}
	return strconv.ParseBool(value)
}

// QueryParamBoolDefault returns a query parameter as boolean or default
func (r *Request) QueryParamBoolDefault(key string, defaultValue bool) bool {
	value, err := r.QueryParamBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// QueryParams returns all query parameters
func (r *Request) QueryParams() url.Values {
	return r.URL.Query()
}

// FormValue returns a form value
func (r *Request) FormValue(key string) string {
	if !r.parsedForm {
		r.ParseForm()
		r.parsedForm = true
	}
	return r.Request.FormValue(key)
}

// FormValueDefault returns a form value or default
func (r *Request) FormValueDefault(key, defaultValue string) string {
	if value := r.FormValue(key); value != "" {
		return value
	}
	return defaultValue
}

// FormValues returns all form values
func (r *Request) FormValues() url.Values {
	if !r.parsedForm {
		r.ParseForm()
		r.parsedForm = true
	}
	return r.Form
}

// PostFormValue returns a POST form value (excluding query params)
func (r *Request) PostFormValue(key string) string {
	if !r.parsedForm {
		r.ParseForm()
		r.parsedForm = true
	}
	return r.Request.PostFormValue(key)
}

// PostFormValues returns all POST form values
func (r *Request) PostFormValues() url.Values {
	if !r.parsedForm {
		r.ParseForm()
		r.parsedForm = true
	}
	return r.PostForm
}

// FileHeader returns the file header for a multipart form file
func (r *Request) FileHeader(key string) (*multipart.FileHeader, error) {
	if !r.parsedMultipart {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
			return nil, err
		}
		r.parsedMultipart = true
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

// Files returns all file headers for a multipart form
func (r *Request) Files() map[string][]*multipart.FileHeader {
	if !r.parsedMultipart {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
			return nil
		}
		r.parsedMultipart = true
	}

	if r.MultipartForm == nil {
		return nil
	}

	return r.MultipartForm.File
}

// HeaderValue returns a header value
func (r *Request) HeaderValue(key string) string {
	return r.Header.Get(key)
}

// HeaderValues returns all header values for a key
func (r *Request) HeaderValues(key string) []string {
	return r.Header.Values(key)
}

// HasHeader checks if a header exists
func (r *Request) HasHeader(key string) bool {
	_, exists := r.Header[key]
	return exists
}

// Cookie returns a cookie value
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	return r.Request.Cookie(name)
}

// CookieValue returns a cookie value as string
func (r *Request) CookieValue(name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// CookieValueDefault returns a cookie value or default
func (r *Request) CookieValueDefault(name, defaultValue string) string {
	value, err := r.CookieValue(name)
	if err != nil {
		return defaultValue
	}
	return value
}

// ClientIP returns the client's IP address
func (r *Request) ClientIP() string {
	// Check for X-Forwarded-For header first
	if ip := r.HeaderValue("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		if idx := len(ip); idx > 0 {
			if commaIdx := 0; commaIdx < idx {
				for i, char := range ip {
					if char == ',' {
						commaIdx = i
						break
					}
				}
				if commaIdx > 0 {
					return ip[:commaIdx]
				}
			}
		}
		return ip
	}

	// Check for X-Real-IP header
	if ip := r.HeaderValue("X-Real-IP"); ip != "" {
		return ip
	}

	// Check for X-Client-IP header
	if ip := r.HeaderValue("X-Client-IP"); ip != "" {
		return ip
	}

	// Fall back to remote address
	if r.RemoteAddr != "" {
		// Remove port from address
		for i := len(r.RemoteAddr) - 1; i >= 0; i-- {
			if r.RemoteAddr[i] == ':' {
				return r.RemoteAddr[:i]
			}
		}
		return r.RemoteAddr
	}

	return ""
}

// UserAgent returns the User-Agent header
func (r *Request) UserAgent() string {
	return r.HeaderValue("User-Agent")
}

// Referer returns the Referer header
func (r *Request) Referer() string {
	return r.HeaderValue("Referer")
}

// IsAjax checks if the request is an AJAX request
func (r *Request) IsAjax() bool {
	return r.HeaderValue("X-Requested-With") == "XMLHttpRequest"
}

// IsSecure checks if the request is HTTPS
func (r *Request) IsSecure() bool {
	return r.TLS != nil || r.HeaderValue("X-Forwarded-Proto") == "https"
}

// Scheme returns the request scheme (http or https)
func (r *Request) Scheme() string {
	if r.IsSecure() {
		return "https"
	}
	return "http"
}

// BaseURL returns the base URL of the request
func (r *Request) BaseURL() string {
	return r.Scheme() + "://" + r.Host
}

// FullURL returns the full URL of the request
func (r *Request) FullURL() string {
	return r.BaseURL() + r.RequestURI
}

// ContentLength returns the content length
func (r *Request) ContentLength() int64 {
	return r.Request.ContentLength
}

// ContentType returns the content type
func (r *Request) ContentType() string {
	return GetContentType(r.Request)
}

// Accept returns the Accept header
func (r *Request) Accept() string {
	return r.HeaderValue("Accept")
}

// AcceptEncoding returns the Accept-Encoding header
func (r *Request) AcceptEncoding() string {
	return r.HeaderValue("Accept-Encoding")
}

// AcceptLanguage returns the Accept-Language header
func (r *Request) AcceptLanguage() string {
	return r.HeaderValue("Accept-Language")
}

// Authorization returns the Authorization header
func (r *Request) Authorization() string {
	return r.HeaderValue("Authorization")
}

// BearerToken extracts Bearer token from Authorization header
func (r *Request) BearerToken() string {
	auth := r.Authorization()
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

// Auth provides comprehensive authentication information
type Auth struct {
	Type     string                 // "bearer", "basic", "api-key", "jwt", "oauth", "custom"
	Token    string                 // The actual token/credential
	Username string                 // For basic auth
	Password string                 // For basic auth
	Claims   map[string]interface{} // For JWT/custom claims
	Valid    bool                   // Whether the auth is valid
	Metadata map[string]string      // Additional auth metadata
}

// GetAuth extracts and analyzes authentication information
func (r *Request) GetAuth() *Auth {
	auth := &Auth{
		Claims:   make(map[string]interface{}),
		Metadata: make(map[string]string),
	}

	authHeader := r.Authorization()
	if authHeader == "" {
		// Check for API key in query params or headers
		if apiKey := r.QueryParam("api_key"); apiKey != "" {
			auth.Type = "api-key"
			auth.Token = apiKey
			auth.Valid = true
			auth.Metadata["source"] = "query"
			return auth
		}

		if apiKey := r.HeaderValue("X-API-Key"); apiKey != "" {
			auth.Type = "api-key"
			auth.Token = apiKey
			auth.Valid = true
			auth.Metadata["source"] = "header"
			return auth
		}

		// Check for session cookie
		if sessionCookie, err := r.Cookie("session"); err == nil {
			auth.Type = "session"
			auth.Token = sessionCookie.Value
			auth.Valid = true
			auth.Metadata["source"] = "cookie"
			return auth
		}

		auth.Valid = false
		return auth
	}

	// Parse Authorization header
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		auth.Type = "bearer"
		auth.Token = authHeader[7:]
		auth.Valid = true

		// Try to detect if it's a JWT
		if r.isJWT(auth.Token) {
			auth.Type = "jwt"
			auth.Claims = r.parseJWTClaims(auth.Token)
		}
	} else if len(authHeader) > 6 && authHeader[:6] == "Basic " {
		auth.Type = "basic"
		if username, password, ok := r.Request.BasicAuth(); ok {
			auth.Username = username
			auth.Password = password
			auth.Valid = true
		}
	} else if len(authHeader) > 7 && authHeader[:7] == "Digest " {
		auth.Type = "digest"
		auth.Token = authHeader[7:]
		auth.Valid = true
	} else if len(authHeader) > 6 && authHeader[:6] == "OAuth " {
		auth.Type = "oauth"
		auth.Token = authHeader[6:]
		auth.Valid = true
	} else {
		auth.Type = "custom"
		auth.Token = authHeader
		auth.Valid = true
	}

	return auth
}

// isJWT checks if a token looks like a JWT
func (r *Request) isJWT(token string) bool {
	parts := 0
	for _, char := range token {
		if char == '.' {
			parts++
		}
	}
	return parts == 2 // JWT has 3 parts separated by 2 dots
}

// parseJWTClaims extracts claims from JWT (basic implementation)
func (r *Request) parseJWTClaims(token string) map[string]interface{} {
	claims := make(map[string]interface{})

	// Split token into parts
	parts := make([]string, 0, 3)
	start := 0
	for i, char := range token {
		if char == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	if start < len(token) {
		parts = append(parts, token[start:])
	}

	if len(parts) != 3 {
		return claims
	}

	// For a complete JWT parser, you'd decode the payload (parts[1])
	// This is a simplified version that just marks common claims
	claims["header"] = parts[0]
	claims["payload"] = parts[1]
	claims["signature"] = parts[2]
	claims["raw_token"] = token

	return claims
}

// HasAuth checks if the request has any authentication
func (r *Request) HasAuth() bool {
	auth := r.GetAuth()
	return auth.Valid
}

// AuthType returns the authentication type
func (r *Request) AuthType() string {
	auth := r.GetAuth()
	return auth.Type
}

// AuthToken returns the authentication token
func (r *Request) AuthToken() string {
	auth := r.GetAuth()
	return auth.Token
}

// IsAuthType checks if the request uses a specific auth type
func (r *Request) IsAuthType(authType string) bool {
	return r.AuthType() == authType
}

// AcceptsJSON checks if the client accepts JSON responses
func (r *Request) AcceptsJSON() bool {
	accept := r.Accept()
	return contains(accept, "application/json") || contains(accept, "*/*")
}

// AcceptsXML checks if the client accepts XML responses
func (r *Request) AcceptsXML() bool {
	accept := r.Accept()
	return contains(accept, "application/xml") || contains(accept, "text/xml")
}

// AcceptsYAML checks if the client accepts YAML responses
func (r *Request) AcceptsYAML() bool {
	accept := r.Accept()
	return contains(accept, "application/x-yaml") || contains(accept, "text/yaml")
}

// AcceptsHTML checks if the client accepts HTML responses
func (r *Request) AcceptsHTML() bool {
	accept := r.Accept()
	return contains(accept, "text/html")
}

// AcceptsPlainText checks if the client accepts plain text responses
func (r *Request) AcceptsPlainText() bool {
	accept := r.Accept()
	return contains(accept, "text/plain")
}

// SmartBind automatically detects content type and binds accordingly
func (r *Request) SmartBind(obj interface{}) error {
	switch {
	case IsJSON(r.Request):
		return BindJSON(r.Request, obj)
	case IsXML(r.Request):
		return BindXML(r.Request, obj)
	case IsYAML(r.Request):
		return BindYAML(r.Request, obj)
	case IsForm(r.Request):
		return BindForm(r.Request, obj)
	default:
		// Try query parameters as fallback
		return BindQuery(r.Request, obj)
	}
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
