package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
	"github.com/melihkorkmaz/gtd/internal/views/pages"
	"github.com/melihkorkmaz/gtd/internal/views/partials"
)

// TaskHandler manages task-related HTTP endpoints
type TaskHandler struct {
	store     models.TaskStore
	templates *TemplateRenderer
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(store models.TaskStore, templatesDir string) (*TaskHandler, error) {
	templates, err := NewTemplateRenderer(templatesDir)
	if err != nil {
		return nil, err
	}

	return &TaskHandler{
		store:     store,
		templates: templates,
	}, nil
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Contexts    []string `json:"contexts,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Helper function to convert Task model to TaskCardInfo for templates
func getTaskCardInfo(task *models.Task) partials.TaskCardInfo {
	// Convert Context type to string slice
	contexts := make([]string, len(task.Contexts))
	for i, ctx := range task.Contexts {
		contexts[i] = string(ctx)
	}

	return partials.TaskCardInfo{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		DueDate:     task.DueDate,
		Contexts:    contexts,
		Tags:        task.Tags,
		CreatedAt:   task.CreatedAt,
		ProjectID:   task.ProjectID,
	}
}

// RegisterRoutes registers all task-related routes
func (h *TaskHandler) RegisterRoutes(r chi.Router) {
	// API routes for JSON responses
	r.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", h.ListTasksAPI)
		r.Post("/", h.CreateTaskAPI)
		r.Post("/quick-capture", h.QuickCaptureAPI)
		r.Get("/search", h.SearchTasksAPI)
		r.Get("/{id}", h.GetTaskAPI)
		r.Put("/{id}", h.UpdateTaskAPI)
		r.Delete("/{id}", h.DeleteTaskAPI)
	})

	// HTML routes for server-side rendering
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.ListTasksPage)
		r.Get("/search", h.SearchTasksPage)
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
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := models.NewTask(request.Title, request.Description, user.ID)

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
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		// User should be authenticated at this point due to middleware,
		// but this is an extra safety check
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Log when this handler is called
	fmt.Printf("ListTasksPage handler called with path: %s and query: %s\n", r.URL.Path, r.URL.RawQuery)

	// Get status filter if provided
	status := r.URL.Query().Get("status")

	var tasks []*models.Task
	var err error
	var title string

	if status != "" {
		// Filter tasks by status and user ID
		fmt.Printf("Filtering tasks by status: %s for user: %s\n", status, user.ID)
		tasks, err = h.store.GetByStatusAndUserID(models.TaskStatus(status), user.ID)

		// Set title based on status
		switch status {
		case "inbox":
			title = "Inbox"
		case "next":
			title = "Next Actions"
		case "waiting":
			title = "Waiting For"
		case "someday":
			title = "Someday/Maybe"
		case "done":
			title = "Completed Tasks"
		case "project":
			title = "Projects"
		default:
			title = "Tasks - " + status
		}
	} else {
		// Get all tasks for this user
		fmt.Printf("Getting all tasks for user: %s\n", user.ID)
		tasks, err = h.store.GetAllByUserID(user.ID)
		title = "All Tasks"
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Found %d tasks\n", len(tasks))

	// Convert tasks to template-friendly format
	taskInfos := make([]partials.TaskCardInfo, len(tasks))
	for i, task := range tasks {
		taskInfos[i] = getTaskCardInfo(task)
	}

	// Add user to context and use the base template
	ctx := context.WithValue(r.Context(), "user", user)

	// Render the page
	tasksPage := pages.TasksListPage(title, taskInfos)
	w.Header().Set("Content-Type", "text/html")
	tasksPage.Render(ctx, w)
}

// ViewTaskPage renders a single task view
func (h *TaskHandler) ViewTaskPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		// User should be authenticated at this point due to middleware,
		// but this is an extra safety check
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get task ID from URL
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert task to template-friendly format
	taskInfo := getTaskCardInfo(task)

	// Add user to context and use the base template
	ctx := context.WithValue(r.Context(), "user", user)

	// Render the page
	w.Header().Set("Content-Type", "text/html")
	pages.TaskDetailPage(taskInfo).Render(ctx, w)
}

// NewTaskForm renders the form for creating a new task
func (h *TaskHandler) NewTaskForm(w http.ResponseWriter, r *http.Request) {
	// This is a placeholder; we'll implement the templ version later
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div class='form-container'><h3>Create Task Form</h3><p>This form will be implemented using templ soon.</p></div>"))
}

// EditTaskForm renders the form for editing a task
func (h *TaskHandler) EditTaskForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.store.Get(id) // Just verify task exists, we'll use it in the actual implementation
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// This is a placeholder; we'll implement the templ version later
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div class='form-container'><h3>Edit Task Form</h3><p>This form will be implemented using templ soon.</p></div>"))
}

// CreateTaskSubmit handles form submission for creating a task
func (h *TaskHandler) CreateTaskSubmit(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	task := models.NewTask(title, description, user.ID)

	// TODO: Parse and add contexts and tags from form

	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

// QuickCaptureAPI handles quick capture submissions via AJAX
func (h *TaskHandler) QuickCaptureAPI(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	task := models.NewTask(title, description, user.ID)
	task.Status = models.StatusInbox // Quick capture always goes to inbox

	if err := h.store.Save(task); err != nil {
		http.Error(w, "Failed to save task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div class='alert alert-success'>Task captured successfully!</div>"))
}

// SearchTasksAPI searches for tasks and returns JSON results
func (h *TaskHandler) SearchTasksAPI(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		// No query, return empty results
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	tasks, err := h.store.SearchByUserID(query, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// SearchTasksPage renders the search results page for tasks
func (h *TaskHandler) SearchTasksPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		// User should be authenticated at this point due to middleware,
		// but this is an extra safety check
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get search query
	query := r.URL.Query().Get("q")

	var tasks []*models.Task
	var err error

	if query != "" {
		// Search tasks for this user
		tasks, err = h.store.SearchByUserID(query, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Convert tasks to template-friendly format
	taskInfos := make([]partials.TaskCardInfo, len(tasks))
	for i, task := range tasks {
		taskInfos[i] = getTaskCardInfo(task)
	}

	// Create search results data
	searchResults := partials.SearchResultsData{
		SearchQuery:  query,
		ResultsCount: len(tasks),
		Tasks:        taskInfos,
	}

	// Render search results using the templ component
	w.Header().Set("Content-Type", "text/html")
	partials.SearchResults(searchResults).Render(r.Context(), w)
}
