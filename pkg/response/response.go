package response

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// JSON sends a JSON response
func JSON(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(obj)
}

// JSONPretty sends a pretty-formatted JSON response
func JSONPretty(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(obj)
}

// String sends a plain text response
func String(w http.ResponseWriter, code int, format string, values ...interface{}) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if len(values) > 0 {
		_, err := fmt.Fprintf(w, format, values...)
		return err
	}
	_, err := w.Write([]byte(format))
	return err
}

// HTML sends an HTML response
func HTML(w http.ResponseWriter, code int, html string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_, err := w.Write([]byte(html))
	return err
}

// XML sends an XML response
func XML(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	return encoder.Encode(obj)
}

// YAML sends a YAML response
func YAML(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	w.WriteHeader(code)

	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	return encoder.Encode(obj)
}

// Data sends raw data response
func Data(w http.ResponseWriter, code int, contentType string, data []byte) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err := w.Write(data)
	return err
}

// Stream sends a streaming response
func Stream(w http.ResponseWriter, code int, contentType string, reader io.Reader) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err := io.Copy(w, reader)
	return err
}

// File sends a file response
func File(w http.ResponseWriter, r *http.Request, filePath string) {
	http.ServeFile(w, r, filePath)
}

// Download sends a file as attachment
func Download(w http.ResponseWriter, r *http.Request, filePath, filename string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

// Error sends an error response
func Error(w http.ResponseWriter, code int, message string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	return json.NewEncoder(w).Encode(response)
}

// Success sends a success response
func Success(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}

	return json.NewEncoder(w).Encode(response)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// SetCacheHeaders sets cache control headers
func SetCacheHeaders(w http.ResponseWriter, maxAge int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
	w.Header().Set("Expires", time.Now().Add(time.Duration(maxAge)*time.Second).Format(http.TimeFormat))
}

// SetNoCacheHeaders sets no-cache headers
func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// Attachment sends a file as attachment with proper headers
func Attachment(w http.ResponseWriter, r *http.Request, filePath, filename string) error {
	// Set security headers
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	// Determine content type from file extension
	contentType := "application/octet-stream"
	if ext := filepath.Ext(filename); ext != "" {
		if ct := getContentTypeFromExt(ext); ct != "" {
			contentType = ct
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	http.ServeFile(w, r, filePath)
	return nil
}

// JSONP sends a JSONP response
func JSONP(w http.ResponseWriter, code int, callback string, obj interface{}) error {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.WriteHeader(code)

	if callback == "" {
		callback = "callback"
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s(%s);", callback, data)
	return err
}

// Redirect sends a redirect response
func Redirect(w http.ResponseWriter, r *http.Request, code int, url string) error {
	http.Redirect(w, r, url, code)
	return nil
}

// getContentTypeFromExt returns content type based on file extension
func getContentTypeFromExt(ext string) string {
	types := map[string]string{
		".pdf":  "application/pdf",
		".zip":  "application/zip",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".json": "application/json",
		".xml":  "application/xml",
		".yaml": "application/x-yaml",
		".yml":  "application/x-yaml",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
	}
	return types[ext]
}
