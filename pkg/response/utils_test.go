package response

import (
	"net/http/httptest"
	"testing"
)

func TestGuessContentType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.json", "application/json"},
		{"test.xml", "application/xml"},
		{"test.yaml", "application/x-yaml"},
		{"test.yml", "application/x-yaml"},
		{"test.html", "text/html"},
		{"test.txt", "text/plain"},
		{"test.pdf", "application/pdf"},
		{"test.png", "image/png"},
		{"test.jpg", "image/jpeg"},
		{"test.unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := GuessContentType(tt.filename)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSetDownloadHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	filename := "test.pdf"

	SetDownloadHeaders(w, filename)

	contentDisposition := w.Header().Get("Content-Disposition")
	expected := `attachment; filename="test.pdf"`
	if contentDisposition != expected {
		t.Errorf("expected Content-Disposition %s, got %s", expected, contentDisposition)
	}

	if w.Header().Get("Content-Description") != "File Transfer" {
		t.Errorf("expected Content-Description 'File Transfer', got %s", w.Header().Get("Content-Description"))
	}
}

func TestSetInlineHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	filename := "document.pdf"

	SetInlineHeaders(w, filename)

	contentDisposition := w.Header().Get("Content-Disposition")
	expected := `inline; filename="document.pdf"`
	if contentDisposition != expected {
		t.Errorf("expected Content-Disposition %s, got %s", expected, contentDisposition)
	}
}

func TestSetContentLength(t *testing.T) {
	w := httptest.NewRecorder()
	length := int64(1024)

	SetContentLength(w, length)

	contentLength := w.Header().Get("Content-Length")
	if contentLength != "1024" {
		t.Errorf("expected Content-Length 1024, got %s", contentLength)
	}
}

func TestIsValidStatusCode(t *testing.T) {
	tests := []struct {
		code  int
		valid bool
	}{
		{200, true},
		{404, true},
		{500, true},
		{99, false},
		{600, false},
		{0, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.code)), func(t *testing.T) {
			result := IsValidStatusCode(tt.code)
			if result != tt.valid {
				t.Errorf("expected %t for code %d, got %t", tt.valid, tt.code, result)
			}
		})
	}
}

func TestStatusText(t *testing.T) {
	tests := []struct {
		code int
		text string
	}{
		{200, "OK"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := StatusText(tt.code)
			if result != tt.text {
				t.Errorf("expected %s for code %d, got %s", tt.text, tt.code, result)
			}
		})
	}
}

func TestStatusCodeHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(w *httptest.ResponseRecorder)
		expected int
	}{
		{
			name:     "BadRequest",
			fn:       func(w *httptest.ResponseRecorder) { BadRequest(w) },
			expected: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.fn(w)
			if w.Code != tt.expected {
				t.Errorf("expected status %d, got %d", tt.expected, w.Code)
			}
		})
	}
}
