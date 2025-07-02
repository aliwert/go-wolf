package response

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CompressionMiddleware provides gzip compression for responses
type CompressionMiddleware struct {
	level int
}

// NewCompressionMiddleware creates a new compression middleware
func NewCompressionMiddleware(level int) *CompressionMiddleware {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}
	return &CompressionMiddleware{level: level}
}

// gzipWriter wraps http.ResponseWriter to provide gzip compression
type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

func (gw gzipWriter) Write(data []byte) (int, error) {
	return gw.Writer.Write(data)
}

// ShouldCompress determines if the response should be compressed
func (cm *CompressionMiddleware) ShouldCompress(r *http.Request, contentType string) bool {
	// Check Accept-Encoding header
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		return false
	}

	// Only compress text-based content types
	compressibleTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-yaml",
	}

	for _, t := range compressibleTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}

	return false
}

// Wrap wraps an http.ResponseWriter with gzip compression
func (cm *CompressionMiddleware) Wrap(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	// Check if we should compress
	contentType := w.Header().Get("Content-Type")
	if !cm.ShouldCompress(r, contentType) {
		return w
	}

	// Set gzip headers
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Vary", "Accept-Encoding")

	// Create gzip writer
	gz, err := gzip.NewWriterLevel(w, cm.level)
	if err != nil {
		return w // Fallback to uncompressed
	}

	return &gzipWriter{
		Writer:         gz,
		ResponseWriter: w,
	}
}

// SecurityMiddleware adds security headers to responses
type SecurityMiddleware struct {
	CSPPolicy string
	HSTS      bool
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{
		CSPPolicy: "default-src 'self'",
		HSTS:      true,
	}
}

// SetCSP sets the Content Security Policy
func (sm *SecurityMiddleware) SetCSP(policy string) {
	sm.CSPPolicy = policy
}

// SetHSTS enables or disables HSTS
func (sm *SecurityMiddleware) SetHSTS(enabled bool) {
	sm.HSTS = enabled
}

// Wrap adds security headers to the response
func (sm *SecurityMiddleware) Wrap(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	// Add security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	if sm.CSPPolicy != "" {
		w.Header().Set("Content-Security-Policy", sm.CSPPolicy)
	}

	if sm.HSTS && r.TLS != nil {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	return w
}

// CORSMiddleware handles Cross-Origin Resource Sharing
type CORSMiddleware struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		MaxAge:         86400,
	}
}

// SetAllowedOrigins sets allowed origins
func (cm *CORSMiddleware) SetAllowedOrigins(origins ...string) {
	cm.AllowedOrigins = origins
}

// SetAllowedMethods sets allowed methods
func (cm *CORSMiddleware) SetAllowedMethods(methods ...string) {
	cm.AllowedMethods = methods
}

// SetAllowedHeaders sets allowed headers
func (cm *CORSMiddleware) SetAllowedHeaders(headers ...string) {
	cm.AllowedHeaders = headers
}

// Wrap adds CORS headers to the response
func (cm *CORSMiddleware) Wrap(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	origin := r.Header.Get("Origin")

	// Check if origin is allowed
	if cm.isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", strings.Join(cm.AllowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(cm.AllowedHeaders, ", "))

	if cm.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cm.MaxAge))
	}

	return w
}

// isOriginAllowed checks if an origin is allowed
func (cm *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowed := range cm.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
