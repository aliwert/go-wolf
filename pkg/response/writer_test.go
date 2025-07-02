package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWriter(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	if writer.Status() != 200 {
		t.Errorf("expected default status 200, got %d", writer.Status())
	}

	if writer.Written() {
		t.Error("expected written to be false initially")
	}

	if writer.Size() != 0 {
		t.Errorf("expected size 0, got %d", writer.Size())
	}
}

func TestWriterWrite(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	data := []byte("hello world")
	n, err := writer.Write(data)

	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}

	if !writer.Written() {
		t.Error("expected written to be true after Write")
	}

	if writer.Size() != len(data) {
		t.Errorf("expected size %d, got %d", len(data), writer.Size())
	}
}

func TestWriterStatus(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.WriteHeader(404)

	if writer.Status() != 404 {
		t.Errorf("expected status 404, got %d", writer.Status())
	}
}

func TestWriterHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetHeader("X-Test", "value")
	writer.AddHeader("X-Multiple", "value1")
	writer.AddHeader("X-Multiple", "value2")

	if writer.Header().Get("X-Test") != "value" {
		t.Errorf("expected X-Test header value, got %s", writer.Header().Get("X-Test"))
	}

	values := writer.Header().Values("X-Multiple")
	if len(values) != 2 {
		t.Errorf("expected 2 X-Multiple values, got %d", len(values))
	}

	writer.DeleteHeader("X-Test")
	if writer.Header().Get("X-Test") != "" {
		t.Error("expected X-Test header to be deleted")
	}
}

func TestWriterContentType(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetContentType("application/json")

	if writer.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", writer.Header().Get("Content-Type"))
	}
}

func TestWriterCache(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetCacheControl(3600)

	cacheControl := writer.Header().Get("Cache-Control")
	if cacheControl != "max-age=3600" {
		t.Errorf("expected Cache-Control max-age=3600, got %s", cacheControl)
	}

	if writer.Header().Get("Expires") == "" {
		t.Error("expected Expires header to be set")
	}
}

func TestWriterNoCache(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetNoCache()

	cacheControl := writer.Header().Get("Cache-Control")
	if cacheControl != "no-cache, no-store, must-revalidate" {
		t.Errorf("unexpected Cache-Control: %s", cacheControl)
	}

	if writer.Header().Get("Pragma") != "no-cache" {
		t.Errorf("expected Pragma no-cache, got %s", writer.Header().Get("Pragma"))
	}

	if writer.Header().Get("Expires") != "0" {
		t.Errorf("expected Expires 0, got %s", writer.Header().Get("Expires"))
	}
}

func TestWriterCORS(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetCORS("https://example.com")

	if writer.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("unexpected CORS origin: %s", writer.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestWriterSecurity(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	writer.SetSecurity()

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expected := range expectedHeaders {
		if writer.Header().Get(header) != expected {
			t.Errorf("expected %s: %s, got %s", header, expected, writer.Header().Get(header))
		}
	}
}

func TestWriterCookie(t *testing.T) {
	w := httptest.NewRecorder()
	writer := NewWriter(w)

	cookie := &http.Cookie{
		Name:  "test",
		Value: "value",
	}

	writer.SetCookie(cookie)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Name != "test" || cookies[0].Value != "value" {
		t.Errorf("expected cookie test=value, got %s=%s", cookies[0].Name, cookies[0].Value)
	}
}
