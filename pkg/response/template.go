package response

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// TemplateRenderer handles HTML template rendering
type TemplateRenderer struct {
	TemplateDir string
	Templates   map[string]*template.Template
	FuncMap     template.FuncMap
	Layout      string
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templateDir string) *TemplateRenderer {
	return &TemplateRenderer{
		TemplateDir: templateDir,
		Templates:   make(map[string]*template.Template),
		FuncMap:     make(template.FuncMap),
	}
}

// AddFunc adds a template function
func (tr *TemplateRenderer) AddFunc(name string, fn interface{}) {
	tr.FuncMap[name] = fn
}

// SetLayout sets the layout template
func (tr *TemplateRenderer) SetLayout(layout string) {
	tr.Layout = layout
}

// LoadTemplates loads all templates from the template directory
func (tr *TemplateRenderer) LoadTemplates() error {
	if tr.TemplateDir == "" {
		return fmt.Errorf("template directory not set")
	}

	templates, err := filepath.Glob(filepath.Join(tr.TemplateDir, "*.html"))
	if err != nil {
		return err
	}

	for _, tmpl := range templates {
		name := filepath.Base(tmpl)
		name = strings.TrimSuffix(name, filepath.Ext(name))

		t := template.New(name).Funcs(tr.FuncMap)

		// If layout is set, parse layout first
		if tr.Layout != "" {
			layoutPath := filepath.Join(tr.TemplateDir, tr.Layout)
			t, err = t.ParseFiles(layoutPath, tmpl)
		} else {
			t, err = t.ParseFiles(tmpl)
		}

		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tmpl, err)
		}

		tr.Templates[name] = t
	}

	return nil
}

// Render renders a template with data
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}) error {
	tmpl, exists := tr.Templates[name]
	if !exists {
		return fmt.Errorf("template %s not found", name)
	}

	return tmpl.Execute(w, data)
}

// RenderHTTP renders a template as HTTP response
func (tr *TemplateRenderer) RenderHTTP(w http.ResponseWriter, code int, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	return tr.Render(w, name, data)
}

// Template renders an HTML template
func Template(w http.ResponseWriter, code int, renderer *TemplateRenderer, name string, data interface{}) error {
	return renderer.RenderHTTP(w, code, name, data)
}

// DefaultTemplateRenderer is a default instance
var DefaultTemplateRenderer *TemplateRenderer

// SetDefaultTemplateDir sets the default template directory
func SetDefaultTemplateDir(dir string) error {
	DefaultTemplateRenderer = NewTemplateRenderer(dir)
	return DefaultTemplateRenderer.LoadTemplates()
}

// RenderTemplate renders a template using the default renderer
func RenderTemplate(w http.ResponseWriter, code int, name string, data interface{}) error {
	if DefaultTemplateRenderer == nil {
		return fmt.Errorf("default template renderer not initialized")
	}
	return DefaultTemplateRenderer.RenderHTTP(w, code, name, data)
}
