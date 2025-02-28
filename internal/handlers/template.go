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
				return "inbox"
			case "next":
				return "next"
			case "waiting":
				return "waiting"
			case "someday":
				return "someday"
			case "done":
				return "done"
			case "project":
				return "project"
			case "reference":
				return "reference"
			case "scheduled":
				return "scheduled"
			default:
				return "neutral"
			}
		},
	}


	// Load layout templates
	layoutPattern := filepath.Join(templatesDir, "layouts/*.html")
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(layoutPattern)
	if err != nil {
		return nil, err
	}

	// Parse partials
	partialPattern := filepath.Join(templatesDir, "partials/*.html")
	if _, err := tmpl.ParseGlob(partialPattern); err != nil {
		return nil, err
	}

	// Parse pages
	pagePattern := filepath.Join(templatesDir, "pages/*.html")
	if _, err := tmpl.ParseGlob(pagePattern); err != nil {
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
