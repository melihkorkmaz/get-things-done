package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
	"github.com/melihkorkmaz/gtd/internal/views/pages"
	"github.com/melihkorkmaz/gtd/internal/views/partials"
)

// ProjectHandler manages project-related HTTP endpoints
type ProjectHandler struct {
	store models.TaskStore
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(store models.TaskStore, templatesDir string) (*ProjectHandler, error) {
	return &ProjectHandler{
		store: store,
	}, nil
}

// RegisterRoutes registers all project-related routes
func (h *ProjectHandler) RegisterRoutes(r chi.Router) {
	// API routes for JSON responses
	r.Route("/api/projects", func(r chi.Router) {
		r.Get("/", h.ListProjectsAPI)
		r.Post("/", h.CreateProjectAPI)
		r.Get("/{id}", h.GetProjectAPI)
		r.Put("/{id}", h.UpdateProjectAPI)
		r.Delete("/{id}", h.DeleteProjectAPI)
		r.Put("/{id}/complete", h.CompleteProjectAPI)
		r.Put("/{id}/archive", h.ArchiveProjectAPI)
		r.Put("/{id}/tasks/{taskId}", h.AddTaskToProjectAPI)
	})

	// HTML routes for server-side rendering
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", h.ListProjectsPage)
		r.Post("/", h.CreateProjectSubmit)
		r.Get("/{id}", h.ViewProjectPage)
		r.Get("/{id}/edit", h.EditProjectForm)
		r.Post("/{id}/tasks", h.AddTaskToProjectSubmit)
	})
}

// ListProjectsAPI returns a JSON list of all projects
func (h *ProjectHandler) ListProjectsAPI(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get filter and sort parameters
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	// Get all tasks with project status for this user
	projects, err := h.getProjects(filter, sort, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// getProjects retrieves projects with optional filtering and sorting
func (h *ProjectHandler) getProjects(filter, sort string, userID string) ([]*models.Task, error) {
	// First get all tasks with project status for this user
	tasks, err := h.store.GetByStatusAndUserID(models.StatusProject, userID)
	if err != nil {
		return nil, err
	}

	// Filter projects if needed
	var filteredProjects []*models.Task
	if filter == "" || filter == "all" {
		filteredProjects = tasks
	} else {
		// In a real implementation, we would filter by project properties
		// For now, we'll just return all projects
		filteredProjects = tasks
	}

	// Sort projects if needed
	// In a real implementation, we would sort by created date, progress, etc.
	// For now, we'll just return the projects as is

	return filteredProjects, nil
}

// CreateProjectAPI creates a new project from JSON input
func (h *ProjectHandler) CreateProjectAPI(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"dueDate,omitempty"`
		Contexts    []string  `json:"contexts,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a new task with project status
	project := models.NewTask(request.Title, request.Description, user.ID)
	project.MarkAsProject()

	// Set due date if provided
	if !request.DueDate.IsZero() {
		project.DueDate = &request.DueDate
	}

	// Convert contexts to Context type
	for _, ctx := range request.Contexts {
		project.Contexts = append(project.Contexts, models.Context(ctx))
	}

	// Set tags
	project.Tags = request.Tags

	if err := h.store.Save(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// GetProjectAPI returns a single project as JSON
func (h *ProjectHandler) GetProjectAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// UpdateProjectAPI updates a project from JSON input
func (h *ProjectHandler) UpdateProjectAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// First get the existing project
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Decode the update request
	if err := json.NewDecoder(r.Body).Decode(project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure it remains a project
	project.Status = models.StatusProject

	// Save the updated project
	if err := h.store.Save(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// DeleteProjectAPI deletes a project
func (h *ProjectHandler) DeleteProjectAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// First check if it's a project
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Delete the project
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CompleteProjectAPI marks a project as complete
func (h *ProjectHandler) CompleteProjectAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// First get the project
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Mark as done
	project.MarkAsDone()

	// Save the updated project
	if err := h.store.Save(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// ArchiveProjectAPI archives a project
func (h *ProjectHandler) ArchiveProjectAPI(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// First get the project
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// For now, we'll just mark it as done since we don't have a separate archive status
	// In a real implementation, you might want to add an "archived" flag to the Task model
	project.MarkAsDone()

	// Save the updated project
	if err := h.store.Save(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// AddTaskToProjectAPI adds a task to a project
func (h *ProjectHandler) AddTaskToProjectAPI(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	// Get project
	project, err := h.store.Get(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Get task
	task, err := h.store.Get(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add task to project
	task.ProjectID = project.ID

	// Save the updated task
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// ListProjectsPage renders the projects list page
func (h *ProjectHandler) ListProjectsPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		// User should be authenticated at this point due to middleware,
		// but this is an extra safety check
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get filter and sort parameters
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	// Get all tasks with project status for this user
	projects, err := h.getProjects(filter, sort, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enhance projects with additional metadata
	enhancedProjects := make([]partials.ProjectInfo, 0, len(projects))
	for _, project := range projects {
		// Get tasks for this project
		projectTasks, err := h.getProjectTasks(project.ID, user.ID)
		if err != nil {
			fmt.Printf("Error getting tasks for project %s: %v\n", project.ID, err)
			continue
		}

		// Count completed tasks
		completedCount := 0
		for _, task := range projectTasks {
			if task.Status == models.StatusDone {
				completedCount++
			}
		}

		// Calculate completion percentage
		completionPercentage := 0
		if len(projectTasks) > 0 {
			completionPercentage = (completedCount * 100) / len(projectTasks)
		}

		// Convert Context type to string slice
		contexts := make([]string, len(project.Contexts))
		for i, ctx := range project.Contexts {
			contexts[i] = string(ctx)
		}

		enhancedProject := partials.ProjectInfo{
			ID:                   project.ID,
			Title:                project.Title,
			Description:          project.Description,
			Status:               string(project.Status),
			DueDate:              project.DueDate,
			Contexts:             contexts,
			Tags:                 project.Tags,
			CreatedAt:            project.CreatedAt,
			TaskCount:            len(projectTasks),
			CompletedTaskCount:   completedCount,
			CompletionPercentage: completionPercentage,
		}

		enhancedProjects = append(enhancedProjects, enhancedProject)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Printf("DEBUG: Rendering projects page with %d projects\n", len(enhancedProjects))

	// Add user to context and use the base template
	ctx := context.WithValue(r.Context(), "user", user)
	
	// Render the page
	component := pages.ProjectsPage(enhancedProjects)
	if err := component.Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getProjectTasks retrieves all tasks associated with a project
func (h *ProjectHandler) getProjectTasks(projectID string, userID string) ([]*models.Task, error) {
	// Get all tasks for this user
	allTasks, err := h.store.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter tasks for this project
	var projectTasks []*models.Task
	for _, task := range allTasks {
		if task.ProjectID == projectID {
			projectTasks = append(projectTasks, task)
		}
	}

	return projectTasks, nil
}

// CreateProjectSubmit handles form submission for creating a new project
func (h *ProjectHandler) CreateProjectSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	dueDateStr := r.FormValue("due_date")

	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Create a new task with project status
	project := models.NewTask(title, description, user.ID)
	project.MarkAsProject()

	// Set due date if provided
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			http.Error(w, "Invalid due date format", http.StatusBadRequest)
			return
		}
		project.DueDate = &dueDate
	}

	// Save project
	if err := h.store.Save(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the project view page
	http.Redirect(w, r, "/projects/"+project.ID, http.StatusSeeOther)
}

// ViewProjectPage renders a single project view
func (h *ProjectHandler) ViewProjectPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		// User should be authenticated at this point due to middleware,
		// but this is an extra safety check
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	id := chi.URLParam(r, "id")
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Get tasks for this project
	projectTasks, err := h.getProjectTasks(project.ID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Count completed tasks
	completedCount := 0
	for _, task := range projectTasks {
		if task.Status == models.StatusDone {
			completedCount++
		}
	}

	// Calculate completion percentage
	completionPercentage := 0
	if len(projectTasks) > 0 {
		completionPercentage = (completedCount * 100) / len(projectTasks)
	}

	// Get available tasks (not assigned to any project)
	availableTasks, err := h.getAvailableTasks(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert Context type to string slice
	contexts := make([]string, len(project.Contexts))
	for i, ctx := range project.Contexts {
		contexts[i] = string(ctx)
	}

	// Create project info for templ
	enhancedProject := partials.ProjectInfo{
		ID:                   project.ID,
		Title:                project.Title,
		Description:          project.Description,
		Status:               string(project.Status),
		DueDate:              project.DueDate,
		Contexts:             contexts,
		Tags:                 project.Tags,
		CreatedAt:            project.CreatedAt,
		TaskCount:            len(projectTasks),
		CompletedTaskCount:   completedCount,
		CompletionPercentage: completionPercentage,
	}

	// Convert tasks
	templTasks := make([]partials.TaskInfo, len(projectTasks))
	for i, task := range projectTasks {
		taskContexts := make([]string, len(task.Contexts))
		for j, ctx := range task.Contexts {
			taskContexts[j] = string(ctx)
		}

		templTasks[i] = partials.TaskInfo{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			DueDate:     task.DueDate,
			Contexts:    taskContexts,
			Tags:        task.Tags,
			CreatedAt:   task.CreatedAt,
			ProjectID:   task.ProjectID,
		}
	}

	// Convert available tasks
	templAvailableTasks := make([]pages.AvailableTask, len(availableTasks))
	for i, task := range availableTasks {
		templAvailableTasks[i] = pages.AvailableTask{
			ID:    task.ID,
			Title: task.Title,
		}
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Printf("DEBUG: Rendering project detail page for project %s with %d tasks\n",
		project.ID, len(projectTasks))

	// Add user to context and use the base template
	ctx := context.WithValue(r.Context(), "user", user)
	
	// Render the page
	component := pages.ProjectDetailPage(enhancedProject, templTasks, templAvailableTasks)
	if err := component.Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getAvailableTasks retrieves tasks that aren't already assigned to a project
func (h *ProjectHandler) getAvailableTasks(userID string) ([]*models.Task, error) {
	// Get all tasks for this user
	allTasks, err := h.store.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter tasks that don't have a project ID and are not projects themselves
	var availableTasks []*models.Task
	for _, task := range allTasks {
		if task.ProjectID == "" && task.Status != models.StatusProject {
			availableTasks = append(availableTasks, task)
		}
	}

	return availableTasks, nil
}

// EditProjectForm renders the form to edit a project
func (h *ProjectHandler) EditProjectForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	// Convert Context type to string slice (for future use)
	contexts := make([]string, len(project.Contexts))
	for i, ctx := range project.Contexts {
		contexts[i] = string(ctx)
	}

	// TODO: Implement edit project form
	// We'll need project info when we create the edit form

	w.Header().Set("Content-Type", "text/html")

	// For now, we'll redirect to the project view page since we don't have an edit template yet
	http.Redirect(w, r, "/projects/"+project.ID, http.StatusSeeOther)
}

// AddTaskToProjectSubmit handles form submission for adding a task to a project
func (h *ProjectHandler) AddTaskToProjectSubmit(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	// Get project
	project, err := h.store.Get(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verify that it's a project
	if project.Status != models.StatusProject {
		http.Error(w, "Not a project", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	status := r.FormValue("status")
	dueDateStr := r.FormValue("due_date")

	// Get user from context if authenticated
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Create a new task
	task := models.NewTask(title, description, user.ID)

	// Set status based on form input
	switch status {
	case "next":
		task.MarkAsNext()
	case "waiting":
		task.MarkAsWaiting()
	}

	// Set project ID
	task.ProjectID = projectID

	// Set due date if provided
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			http.Error(w, "Invalid due date format", http.StatusBadRequest)
			return
		}
		task.DueDate = &dueDate
	}

	// Save task
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the project view
	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}

// generateProjectsPageHtml creates the HTML for the projects page
