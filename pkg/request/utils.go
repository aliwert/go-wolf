// Package request provides utilities for binding and validating HTTP request data
package request

import (
	"net/http"
	"strings"
)

// GetContentType returns the content type of the request
func GetContentType(r *http.Request) string {
	return r.Header.Get("Content-Type")
}

// IsJSON checks if the request content type is JSON
func IsJSON(r *http.Request) bool {
	contentType := GetContentType(r)
	return strings.Contains(contentType, "application/json")
}

// IsForm checks if the request content type is form data
func IsForm(r *http.Request) bool {
	contentType := GetContentType(r)
	return strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data")
}

// IsXML checks if the request content type is XML
func IsXML(r *http.Request) bool {
	contentType := GetContentType(r)
	return strings.Contains(contentType, "application/xml") ||
		strings.Contains(contentType, "text/xml")
}

// IsYAML checks if the request content type is YAML
func IsYAML(r *http.Request) bool {
	contentType := GetContentType(r)
	return strings.Contains(contentType, "application/x-yaml") ||
		strings.Contains(contentType, "application/yaml") ||
		strings.Contains(contentType, "text/yaml")
}
