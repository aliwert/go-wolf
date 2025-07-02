package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// BindJSON binds the request body to a struct using JSON
func BindJSON(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return Validate(obj)
}

// BindXML binds the request body to a struct using XML
func BindXML(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	decoder := xml.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	return Validate(obj)
}

// BindYAML binds the request body to a struct using YAML
func BindYAML(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	decoder := yaml.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	return Validate(obj)
}

// SmartBind automatically detects content type and binds accordingly
func SmartBind(r *http.Request, obj interface{}) error {
	contentType := GetContentType(r)

	switch {
	case IsJSON(r):
		return BindJSON(r, obj)
	case IsXML(r):
		return BindXML(r, obj)
	case IsYAML(r):
		return BindYAML(r, obj)
	case IsForm(r):
		return BindForm(r, obj)
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
}
