package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
)

// TaskHandler manages task-related HTTP endpoints
type TaskHandler struct {
	store models.TaskStore
	tmpl  *template.Template
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(store models.TaskStore, templatesDir string) (*TaskHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &TaskHandler{
		store: store,
		tmpl:  tmpl,
	}, nil
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Contexts    []string `json:"contexts,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// RegisterRoutes registers all task-related routes
func (h *TaskHandler) RegisterRoutes(r chi.Router) {
	// API routes for JSON responses
	r.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", h.ListTasksAPI)
		r.Post("/", h.CreateTaskAPI)
		r.Get("/{id}", h.GetTaskAPI)
		r.Put("/{id}", h.UpdateTaskAPI)
		r.Delete("/{id}", h.DeleteTaskAPI)
	})

	// HTML routes for server-side rendering
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.ListTasksPage)
		r.Get("/new", h.NewTaskForm)
		r.Post("/", h.CreateTaskSubmit)
		r.Get("/{id}", h.ViewTaskPage)
		r.Get("/{id}/edit", h.EditTaskForm)
	})
}

// ListTasksAPI returns a JSON list of all tasks
func (h *TaskHandler) ListTasksAPI(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CreateTaskAPI creates a new task from JSON input
func (h *TaskHandler) CreateTaskAPI(w http.ResponseWriter, r *http.Request) {
	var request CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := models.NewTask(request.Title, request.Description)
	
	// Convert string contexts to Context type
	for _, ctx := range request.Contexts {
		task.Contexts = append(task.Contexts, models.Context(ctx))
	}
	
	// Add tags
	task.Tags = request.Tags

	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTaskAPI returns a single task as JSON
func (h *TaskHandler) GetTaskAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTaskAPI updates a task from JSON input
func (h *TaskHandler) UpdateTaskAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// First get the existing task
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Decode the update request
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Save the updated task
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTaskAPI deletes a task
func (h *TaskHandler) DeleteTaskAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTasksPage renders the tasks list page
func (h *TaskHandler) ListTasksPage(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "All Tasks",
		"Tasks": tasks,
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "tasks.html", data)
}

// NewTaskForm renders the form to create a new task
func (h *TaskHandler) NewTaskForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "task_form.html", nil)
}

// CreateTaskSubmit processes a task creation form submission
func (h *TaskHandler) CreateTaskSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	task := models.NewTask(title, description)
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return a partial HTML response
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		h.tmpl.ExecuteTemplate(w, "task_row.html", task)
		return
	}

	// Otherwise redirect to the tasks list
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

// ViewTaskPage renders a single task view
func (h *TaskHandler) ViewTaskPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": task.Title,
		"Task":  task,
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "task_view.html", data)
}

// EditTaskForm renders the form to edit a task
func (h *TaskHandler) EditTaskForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Edit Task",
		"Task":  task,
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "task_form.html", data)
}