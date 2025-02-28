package handlers

import (
	"html/template"
	"path/filepath"
	"time"
)

// TemplateRenderer handles rendering of templates
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templatesDir string) (*TemplateRenderer, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("Jan 02, 2006 3:04 PM")
		},
		"formatDate": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("Jan 02, 2006")
		},
		"taskStatusBadge": func(status string) string {
			switch status {
			case "inbox":
				return "neutral"
			case "next":
				return "primary"
			case "waiting":
				return "warning"
			case "someday":
				return "secondary"
			case "done":
				return "success"
			case "project":
				return "info"
			case "scheduled":
				return "accent"
			default:
				return "ghost"
			}
		},
	}

	// Parse templates with the function map
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templatesDir, "layouts/*.html"))
	if err != nil {
		return nil, err
	}

	// Parse partials
	if _, err := tmpl.ParseGlob(filepath.Join(templatesDir, "partials/*.html")); err != nil {
		return nil, err
	}

	// Parse pages
	if _, err := tmpl.ParseGlob(filepath.Join(templatesDir, "pages/*.html")); err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		templates: tmpl,
	}, nil
}

// Render renders a template with the given data
func (tr *TemplateRenderer) Render(name string, data interface{}) (string, error) {
	// TODO: Implement rendering to string if needed
	return "", nil
}