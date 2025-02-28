package handlers

import (
	"fmt"
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

	// Parse templates with the function map
	fmt.Println("Loading templates from:", templatesDir)
	
	// Load layout templates
	layoutPattern := filepath.Join(templatesDir, "layouts/*.html")
	fmt.Println("Loading layouts from:", layoutPattern)
	layoutFiles, err := filepath.Glob(layoutPattern)
	if err != nil {
		return nil, fmt.Errorf("error finding layout templates: %v", err)
	}
	fmt.Println("Found layout files:", layoutFiles)
	
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(layoutPattern)
	if err != nil {
		return nil, err
	}

	// Parse partials
	partialPattern := filepath.Join(templatesDir, "partials/*.html")
	fmt.Println("Loading partials from:", partialPattern)
	partialFiles, err := filepath.Glob(partialPattern)
	if err != nil {
		return nil, fmt.Errorf("error finding partial templates: %v", err)
	}
	fmt.Println("Found partial files:", partialFiles)
	
	if _, err := tmpl.ParseGlob(partialPattern); err != nil {
		return nil, err
	}

	// Parse pages
	pagePattern := filepath.Join(templatesDir, "pages/*.html")
	fmt.Println("Loading pages from:", pagePattern)
	pageFiles, err := filepath.Glob(pagePattern)
	if err != nil {
		return nil, fmt.Errorf("error finding page templates: %v", err)
	}
	fmt.Println("Found page files:", pageFiles)
	
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