package response

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// ContentTypeGuesser determines content type from file extension
var contentTypes = map[string]string{
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "application/x-yaml",
	".yml":  "application/x-yaml",
	".html": "text/html",
	".htm":  "text/html",
	".txt":  "text/plain",
	".css":  "text/css",
	".js":   "application/javascript",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".zip":  "application/zip",
	".gz":   "application/gzip",
	".tar":  "application/x-tar",
	".mp3":  "audio/mpeg",
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
}

// GuessContentType guesses content type from file extension
func GuessContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	// Try mime package as fallback
	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}

	return "application/octet-stream"
}

// SetDownloadHeaders sets headers for file downloads
func SetDownloadHeaders(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Transfer-Encoding", "binary")
}

// SetInlineHeaders sets headers for inline file display
func SetInlineHeaders(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
}

// SetContentLength sets the Content-Length header
func SetContentLength(w http.ResponseWriter, length int64) {
	w.Header().Set("Content-Length", strconv.FormatInt(length, 10))
}

// IsValidStatusCode checks if a status code is valid
func IsValidStatusCode(code int) bool {
	return code >= 100 && code < 600
}

// StatusText returns the status text for a given status code
func StatusText(code int) string {
	return http.StatusText(code)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

// SetSecurityHeaders sets common security headers
func SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}

// SetCORSHeaders sets CORS headers
func SetCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "86400")
}
