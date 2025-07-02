// Package response provides utilities for handling HTTP responses
package response

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// Writer wraps http.ResponseWriter with additional functionality
type Writer struct {
	http.ResponseWriter
	statusCode int
	written    bool
	size       int
	mu         sync.RWMutex
}

// NewWriter creates a new Response wrapper
func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: w,
		statusCode:     200, // Default status code
	}
}

// WriteHeader captures the status code
func (w *Writer) WriteHeader(code int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.written {
		return // Already written
	}

	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the response size and sets written flag
func (w *Writer) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.written {
		// If WriteHeader wasn't called explicitly, use default status
		w.ResponseWriter.WriteHeader(w.statusCode)
		w.written = true
	}

	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

// Status returns the HTTP status code
func (w *Writer) Status() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.statusCode
}

// Size returns the response size in bytes
func (w *Writer) Size() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.size
}

// Written returns whether the response has been written
func (w *Writer) Written() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.written
}

// Hijack implements http.Hijacker interface
func (w *Writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
}

// Flush implements http.Flusher interface
func (w *Writer) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// CloseNotify implements http.CloseNotifier interface
func (w *Writer) CloseNotify() <-chan bool {
	if notifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return notifier.CloseNotify()
	}
	// Return a channel that never sends
	return make(<-chan bool)
}

// Push implements http.Pusher interface for HTTP/2 server push
func (w *Writer) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("the ResponseWriter doesn't support HTTP/2 server push")
}

// SetCookie sets a cookie
func (w *Writer) SetCookie(cookie *http.Cookie) {
	http.SetCookie(w.ResponseWriter, cookie)
}

// SetHeader sets a response header
func (w *Writer) SetHeader(key, value string) {
	w.Header().Set(key, value)
}

// AddHeader adds a response header
func (w *Writer) AddHeader(key, value string) {
	w.Header().Add(key, value)
}

// DeleteHeader deletes a response header
func (w *Writer) DeleteHeader(key string) {
	w.Header().Del(key)
}

// SetContentType sets the Content-Type header
func (w *Writer) SetContentType(contentType string) {
	w.SetHeader("Content-Type", contentType)
}

// SetCacheControl sets cache control headers
func (w *Writer) SetCacheControl(maxAge int) {
	w.SetHeader("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
	w.SetHeader("Expires", time.Now().Add(time.Duration(maxAge)*time.Second).Format(http.TimeFormat))
}

// SetNoCache sets no-cache headers
func (w *Writer) SetNoCache() {
	w.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	w.SetHeader("Pragma", "no-cache")
	w.SetHeader("Expires", "0")
}

// SetCORS sets CORS headers
func (w *Writer) SetCORS(origin string) {
	w.SetHeader("Access-Control-Allow-Origin", origin)
	w.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// SetSecurity sets security headers
func (w *Writer) SetSecurity() {
	w.SetHeader("X-Content-Type-Options", "nosniff")
	w.SetHeader("X-Frame-Options", "DENY")
	w.SetHeader("X-XSS-Protection", "1; mode=block")
	w.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}
