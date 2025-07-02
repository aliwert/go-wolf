package response

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type TestData struct {
	Name  string `json:"name" xml:"name" yaml:"name"`
	Value int    `json:"value" xml:"value" yaml:"value"`
}

func TestJSON(t *testing.T) {
	data := TestData{Name: "test", Value: 123}
	w := httptest.NewRecorder()

	err := JSON(w, 200, data)
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected JSON content type, got %s", contentType)
	}

	var result TestData
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Name != data.Name || result.Value != data.Value {
		t.Errorf("expected %+v, got %+v", data, result)
	}
}

func TestJSONPretty(t *testing.T) {
	data := TestData{Name: "test", Value: 123}
	w := httptest.NewRecorder()

	err := JSONPretty(w, 200, data)
	if err != nil {
		t.Fatalf("JSONPretty() error = %v", err)
	}

	body := w.Body.String()
	if !strings.Contains(body, "\n") || !strings.Contains(body, "  ") {
		t.Error("expected pretty-formatted JSON with indentation")
	}
}

func TestXML(t *testing.T) {
	data := TestData{Name: "test", Value: 123}
	w := httptest.NewRecorder()

	err := XML(w, 200, data)
	if err != nil {
		t.Fatalf("XML() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") {
		t.Errorf("expected XML content type, got %s", contentType)
	}

	var result TestData
	if err := xml.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal XML response: %v", err)
	}

	if result.Name != data.Name || result.Value != data.Value {
		t.Errorf("expected %+v, got %+v", data, result)
	}
}

func TestYAML(t *testing.T) {
	data := TestData{Name: "test", Value: 123}
	w := httptest.NewRecorder()

	err := YAML(w, 200, data)
	if err != nil {
		t.Fatalf("YAML() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/x-yaml") {
		t.Errorf("expected YAML content type, got %s", contentType)
	}

	var result TestData
	if err := yaml.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal YAML response: %v", err)
	}

	if result.Name != data.Name || result.Value != data.Value {
		t.Errorf("expected %+v, got %+v", data, result)
	}
}

func TestString(t *testing.T) {
	w := httptest.NewRecorder()

	err := String(w, 200, "Hello, %s!", "World")
	if err != nil {
		t.Fatalf("String() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("expected text/plain content type, got %s", contentType)
	}

	body := w.Body.String()
	expected := "Hello, World!"
	if body != expected {
		t.Errorf("expected %s, got %s", expected, body)
	}
}

func TestHTML(t *testing.T) {
	w := httptest.NewRecorder()
	htmlContent := "<h1>Hello, World!</h1>"

	err := HTML(w, 200, htmlContent)
	if err != nil {
		t.Fatalf("HTML() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected text/html content type, got %s", contentType)
	}

	body := w.Body.String()
	if body != htmlContent {
		t.Errorf("expected %s, got %s", htmlContent, body)
	}
}

func TestDataResponse(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("binary data")
	contentType := "application/octet-stream"

	err := Data(w, 200, contentType, data)
	if err != nil {
		t.Fatalf("Data() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != contentType {
		t.Errorf("expected content type %s, got %s", contentType, w.Header().Get("Content-Type"))
	}

	if !bytes.Equal(w.Body.Bytes(), data) {
		t.Errorf("expected %v, got %v", data, w.Body.Bytes())
	}
}

func TestStream(t *testing.T) {
	w := httptest.NewRecorder()
	data := "streaming data"
	reader := strings.NewReader(data)

	err := Stream(w, 200, "text/plain", reader)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != data {
		t.Errorf("expected %s, got %s", data, body)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	code := 400
	message := "Bad Request"

	err := Error(w, code, message)
	if err != nil {
		t.Fatalf("Error() error = %v", err)
	}

	if w.Code != code {
		t.Errorf("expected status %d, got %d", code, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected JSON content type, got %s", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	errorObj, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error object in response")
	}

	if errorObj["code"] != float64(code) {
		t.Errorf("expected error code %d, got %v", code, errorObj["code"])
	}

	if errorObj["message"] != message {
		t.Errorf("expected error message %s, got %v", message, errorObj["message"])
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := TestData{Name: "test", Value: 123}

	err := Success(w, 200, data)
	if err != nil {
		t.Fatalf("Success() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal success response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("expected success to be true")
	}

	if response["data"] == nil {
		t.Error("expected data field in response")
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	err := NoContent(w)
	if err != nil {
		t.Fatalf("NoContent() error = %v", err)
	}

	if w.Code != 204 {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	if w.Body.Len() != 0 {
		t.Error("expected empty body for 204 No Content")
	}
}

func TestJSONP(t *testing.T) {
	w := httptest.NewRecorder()
	data := TestData{Name: "test", Value: 123}
	callback := "myCallback"

	err := JSONP(w, 200, callback, data)
	if err != nil {
		t.Fatalf("JSONP() error = %v", err)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/javascript") {
		t.Errorf("expected JavaScript content type, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.HasPrefix(body, callback+"(") {
		t.Errorf("expected JSONP response to start with %s(, got %s", callback, body)
	}

	if !strings.HasSuffix(body, ");") {
		t.Error("expected JSONP response to end with );")
	}
}

func TestSetCacheHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	maxAge := 3600

	SetCacheHeaders(w, maxAge)

	cacheControl := w.Header().Get("Cache-Control")
	expected := "max-age=3600"
	if cacheControl != expected {
		t.Errorf("expected Cache-Control %s, got %s", expected, cacheControl)
	}

	if w.Header().Get("Expires") == "" {
		t.Error("expected Expires header to be set")
	}
}

func TestSetNoCacheHeaders(t *testing.T) {
	w := httptest.NewRecorder()

	SetNoCacheHeaders(w)

	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "no-cache") {
		t.Errorf("expected no-cache in Cache-Control, got %s", cacheControl)
	}

	pragma := w.Header().Get("Pragma")
	if pragma != "no-cache" {
		t.Errorf("expected Pragma no-cache, got %s", pragma)
	}

	expires := w.Header().Get("Expires")
	if expires != "0" {
		t.Errorf("expected Expires 0, got %s", expires)
	}
}

func TestGetContentTypeFromExt(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".json", "application/json"},
		{".xml", "application/xml"},
		{".pdf", "application/pdf"},
		{".png", "image/png"},
		{".jpg", "image/jpeg"},
		{".txt", "text/plain"},
		{".unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := getContentTypeFromExt(tt.ext)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
