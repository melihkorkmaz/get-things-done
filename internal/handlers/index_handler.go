package handlers

import (
	"net/http"

	"github.com/melihkorkmaz/gtd/internal/models"
)

// IndexHandler handles the home page
type IndexHandler struct {
	store     models.TaskStore
	templates *TemplateRenderer
}

// NewIndexHandler creates a new index handler
func NewIndexHandler(store models.TaskStore, templatesDir string) (*IndexHandler, error) {
	templates, err := NewTemplateRenderer(templatesDir)
	if err != nil {
		return nil, err
	}

	return &IndexHandler{
		store:     store,
		templates: templates,
	}, nil
}

// TaskStats represents statistics about tasks
type TaskStats struct {
	Total    int
	Inbox    int
	Next     int
	Waiting  int
	Someday  int
	Projects int
	Done     int
}

// HomePage renders the home page
func (h *IndexHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	// Get counts for different task statuses
	stats := TaskStats{}
	
	// Get all tasks
	tasks, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Count tasks by status
	stats.Total = len(tasks)
	for _, task := range tasks {
		switch task.Status {
		case models.StatusInbox:
			stats.Inbox++
		case models.StatusNext:
			stats.Next++
		case models.StatusWaiting:
			stats.Waiting++
		case models.StatusSomeday:
			stats.Someday++
		case models.StatusProject:
			stats.Projects++
		case models.StatusDone:
			stats.Done++
		}
	}
	
	data := map[string]interface{}{
		"Title": "Welcome to GTD App",
		"Stats": stats,
	}
	
	w.Header().Set("Content-Type", "text/html")
	h.templates.templates.ExecuteTemplate(w, "base.html", data)
}