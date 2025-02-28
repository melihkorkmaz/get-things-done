package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
)

// ProjectHandler manages project-related HTTP endpoints
type ProjectHandler struct {
	store     models.TaskStore
	templates *TemplateRenderer
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(store models.TaskStore, templatesDir string) (*ProjectHandler, error) {
	templates, err := NewTemplateRenderer(templatesDir)
	if err != nil {
		return nil, err
	}

	return &ProjectHandler{
		store:     store,
		templates: templates,
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
	// Get filter and sort parameters
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	// Get all tasks with project status
	projects, err := h.getProjects(filter, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// getProjects retrieves projects with optional filtering and sorting
func (h *ProjectHandler) getProjects(filter, sort string) ([]*models.Task, error) {
	// First get all tasks with project status
	tasks, err := h.store.GetByStatus(models.StatusProject)
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

	// Create a new task with project status
	project := models.NewTask(request.Title, request.Description)
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
	// Get filter and sort parameters
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	// Get all tasks with project status
	projects, err := h.getProjects(filter, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Enhance projects with additional metadata
	enhancedProjects := make([]map[string]interface{}, 0, len(projects))
	for _, project := range projects {
		// Get tasks for this project
		projectTasks, err := h.getProjectTasks(project.ID)
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
		
		enhancedProject := map[string]interface{}{
			"ID":                  project.ID,
			"Title":               project.Title,
			"Description":         project.Description,
			"Status":              project.Status,
			"DueDate":             project.DueDate,
			"Contexts":            project.Contexts,
			"Tags":                project.Tags,
			"CreatedAt":           project.CreatedAt,
			"TaskCount":           len(projectTasks),
			"CompletedTaskCount":  completedCount,
			"CompletionPercentage": completionPercentage,
		}
		
		enhancedProjects = append(enhancedProjects, enhancedProject)
	}

	data := map[string]interface{}{
		"Title":    "Projects",
		"Projects": enhancedProjects,
	}

	w.Header().Set("Content-Type", "text/html")
	// Direct HTML rendering for now, similar to our task handler approach
	fmt.Printf("DEBUG: Rendering projects page with %d projects\n", len(enhancedProjects))
	
	// For debugging
	for i, p := range enhancedProjects {
		fmt.Printf("Project %d: %s, Tasks: %d, Completion: %d%%\n", 
			i+1, p["Title"], p["TaskCount"], p["CompletionPercentage"])
	}
	
	// Since our template-based approach is having issues, let's generate HTML directly
	projectsHtml := generateProjectsPageHtml(enhancedProjects)
	w.Write([]byte(projectsHtml))
}

// getProjectTasks retrieves all tasks associated with a project
func (h *ProjectHandler) getProjectTasks(projectID string) ([]*models.Task, error) {
	// Get all tasks
	allTasks, err := h.store.GetAll()
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
	
	// Create a new task with project status
	project := models.NewTask(title, description)
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
	projectTasks, err := h.getProjectTasks(project.ID)
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
	availableTasks, err := h.getAvailableTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Enhance project with additional metadata
	enhancedProject := map[string]interface{}{
		"ID":                  project.ID,
		"Title":               project.Title,
		"Description":         project.Description,
		"Status":              project.Status,
		"DueDate":             project.DueDate,
		"Contexts":            project.Contexts,
		"Tags":                project.Tags,
		"CreatedAt":           project.CreatedAt,
		"TaskCount":           len(projectTasks),
		"CompletedTaskCount":  completedCount,
		"CompletionPercentage": completionPercentage,
	}

	data := map[string]interface{}{
		"Title":          project.Title,
		"Project":        enhancedProject,
		"Tasks":          projectTasks,
		"AvailableTasks": availableTasks,
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Printf("DEBUG: Rendering project detail page for project %s with %d tasks\n", 
		project.ID, len(projectTasks))
		
	// Since our template-based approach is having issues, let's generate HTML directly
	// This function isn't implemented yet, but would follow the same pattern as generateProjectsPageHtml
	// For now, show a basic page with project details
	projectDetailHtml := `
<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + project.Title + ` - GTD App</title>
    
    <!-- DaisyUI with Tailwind CSS -->
    <link href="https://cdn.jsdelivr.net/npm/daisyui@3.9.4/dist/full.css" rel="stylesheet" type="text/css" />
    <script src="https://cdn.tailwindcss.com"></script>
    
    <!-- Alpine.js -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.13.3/dist/cdn.min.js"></script>
    
    <!-- HTMX -->
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    
    <!-- Custom CSS -->
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body class="min-h-screen bg-base-200">
    <div class="container mx-auto p-4">
        <header class="navbar bg-base-100 rounded-box shadow-lg mb-6">
            <div class="flex-1">
                <a href="/" class="btn btn-ghost text-xl">GTD App</a>
            </div>
            <div class="flex-none">
                <ul class="menu menu-horizontal px-1">
                    <li><a href="/tasks">All Tasks</a></li>
                    <li><a href="/tasks?status=inbox">Inbox</a></li>
                    <li><a href="/tasks?status=next">Next Actions</a></li>
                    <li><a href="/tasks?status=waiting">Waiting For</a></li>
                    <li><a href="/projects" class="text-primary font-medium">Projects</a></li>
                    <li><a href="/tasks?status=someday">Someday/Maybe</a></li>
                    <li><a href="/weekly-review" class="text-accent">Weekly Review</a></li>
                </ul>
            </div>
        </header>

        <main>
            <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <div class="flex justify-between items-center mb-6">
                        <h2 class="card-title text-2xl">` + project.Title + `</h2>
                        <div class="badge badge-primary">` + string(project.Status) + `</div>
                    </div>
                    
                    <p class="my-4">` + project.Description + `</p>
                    
                    <div class="flex justify-between items-center mb-4">
                        <h3 class="text-xl font-bold">Project Tasks</h3>
                        <button class="btn btn-primary">Add Task</button>
                    </div>
                    
                    <div class="overflow-x-auto">
                        <table class="table w-full">
                            <thead>
                                <tr>
                                    <th>Task</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
`
	
	// Add tasks to the project detail page
	if len(projectTasks) > 0 {
		for _, task := range projectTasks {
			projectDetailHtml += `
                                <tr>
                                    <td>` + task.Title + `</td>
                                    <td><div class="badge badge-primary">` + string(task.Status) + `</div></td>
                                    <td>
                                        <a href="/tasks/` + task.ID + `" class="btn btn-xs">View</a>
                                    </td>
                                </tr>
`
		}
	} else {
		projectDetailHtml += `
                                <tr>
                                    <td colspan="3" class="text-center">No tasks added to this project yet</td>
                                </tr>
`
	}
	
	projectDetailHtml += `
                            </tbody>
                        </table>
                    </div>
                    
                    <div class="mt-4">
                        <a href="/projects" class="btn">Back to Projects</a>
                    </div>
                </div>
            </div>
        </main>
    </div>
</body>
</html>
`
	
	w.Write([]byte(projectDetailHtml))
}

// getAvailableTasks retrieves tasks that aren't already assigned to a project
func (h *ProjectHandler) getAvailableTasks() ([]*models.Task, error) {
	// Get all tasks
	allTasks, err := h.store.GetAll()
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

	data := map[string]interface{}{
		"Title":   "Edit Project",
		"Project": project,
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.templates.ExecuteTemplate(w, "project_form", data)
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
	
	// Create a new task
	task := models.NewTask(title, description)
	
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
func generateProjectsPageHtml(projects []map[string]interface{}) string {
	html := `
<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Projects - GTD App</title>
    
    <!-- DaisyUI with Tailwind CSS -->
    <link href="https://cdn.jsdelivr.net/npm/daisyui@3.9.4/dist/full.css" rel="stylesheet" type="text/css" />
    <script src="https://cdn.tailwindcss.com"></script>
    
    <!-- Alpine.js -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.13.3/dist/cdn.min.js"></script>
    
    <!-- HTMX -->
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    
    <!-- Custom CSS -->
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body class="min-h-screen bg-base-200">
    <div class="container mx-auto p-4">
        <header class="navbar bg-base-100 rounded-box shadow-lg mb-6">
            <div class="flex-1">
                <a href="/" class="btn btn-ghost text-xl">GTD App</a>
            </div>
            <div class="flex-none">
                <!-- Search box -->
                <div class="form-control mx-2">
                    <form action="/tasks/search" method="GET" class="flex">
                        <div class="relative">
                            <input type="text" name="q" placeholder="Search tasks..." 
                                   class="input input-bordered w-24 md:w-auto" />
                            <div class="search-indicator">
                                <span class="loading loading-spinner loading-xs"></span>
                            </div>
                        </div>
                        <button type="submit" class="btn btn-primary">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0118 0z" />
                            </svg>
                        </button>
                    </form>
                </div>
                <!-- Quick Capture Button -->
                <button class="btn btn-success btn-sm mx-2" onclick="document.getElementById('quick-capture-modal').showModal()">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                    </svg>
                    Capture
                </button>
                <ul class="menu menu-horizontal px-1">
                    <li><a href="/tasks">All Tasks</a></li>
                    <li><a href="/tasks?status=inbox">Inbox</a></li>
                    <li><a href="/tasks?status=next">Next Actions</a></li>
                    <li><a href="/tasks?status=waiting">Waiting For</a></li>
                    <li><a href="/projects" class="text-primary font-medium">Projects</a></li>
                    <li><a href="/tasks?status=someday">Someday/Maybe</a></li>
                    <li><a href="/weekly-review" class="text-accent">Weekly Review</a></li>
                </ul>
            </div>
        </header>

        <main>
            <div id="main-content">
                <div class="card bg-base-100 shadow-xl">
                    <div class="card-body">
                        <div class="flex justify-between items-center mb-6">
                            <h2 class="card-title text-2xl">Projects</h2>
                            <button class="btn btn-primary" onclick="document.getElementById('new-project-modal').showModal()">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                                </svg>
                                New Project
                            </button>
                        </div>
                        
                        <!-- Projects filter options -->
                        <div class="flex flex-wrap gap-2 mb-4">
                            <div class="form-control">
                                <div class="input-group">
                                    <span>Filter:</span>
                                    <select class="select select-bordered" id="project-filter">
                                        <option value="all">All Projects</option>
                                        <option value="active">Active</option>
                                        <option value="completed">Completed</option>
                                        <option value="on-hold">On Hold</option>
                                    </select>
                                </div>
                            </div>
                            
                            <div class="form-control">
                                <div class="input-group">
                                    <span>Sort:</span>
                                    <select class="select select-bordered" id="project-sort">
                                        <option value="created-desc">Newest First</option>
                                        <option value="created-asc">Oldest First</option>
                                        <option value="progress-asc">Least Progress</option>
                                        <option value="progress-desc">Most Progress</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        
                        <!-- Projects List -->
                        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4" id="projects-container">
`

	if len(projects) > 0 {
		for _, project := range projects {
			html += generateProjectCardHtml(project)
		}
	} else {
		html += `
                            <div class="col-span-3 alert">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                                <span>No projects found. Create your first project with the "New Project" button.</span>
                            </div>
`
	}

	html += `
                        </div>
                    </div>
                </div>
            </div>
        </main>

        <footer class="footer footer-center p-4 bg-base-100 text-base-content mt-6 rounded-box">
            <div>
                <p>Copyright © 2025 - All rights reserved</p>
            </div>
        </footer>
    </div>

    <!-- New Project Modal -->
    <dialog id="new-project-modal" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">Create New Project</h3>
            <p class="py-2">Enter the details for your new project.</p>
            
            <form method="POST" action="/projects" id="new-project-form">
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Project Name</span>
                    </label>
                    <input type="text" name="title" placeholder="Enter project name..." class="input input-bordered" required />
                </div>
                
                <div class="form-control mt-2">
                    <label class="label">
                        <span class="label-text">Description</span>
                    </label>
                    <textarea name="description" placeholder="Enter project description..." class="textarea textarea-bordered" rows="3"></textarea>
                </div>
                
                <div class="form-control mt-2">
                    <label class="label">
                        <span class="label-text">Due Date (Optional)</span>
                    </label>
                    <input type="date" name="due_date" class="input input-bordered" />
                </div>
                
                <div class="form-control mt-4">
                    <button type="submit" class="btn btn-primary">Create Project</button>
                </div>
            </form>
            
            <div class="modal-action">
                <form method="dialog">
                    <button class="btn">Close</button>
                </form>
            </div>
        </div>
    </dialog>

    <!-- Quick capture modal -->
    <dialog id="quick-capture-modal" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">Quick Capture</h3>
            <p class="py-2">Quickly capture a new task or idea. It will be added to your inbox for processing later.</p>
            
            <form method="POST" action="/tasks">
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Task/Idea Title</span>
                    </label>
                    <input type="text" name="title" placeholder="Enter title..." class="input input-bordered" required />
                </div>
                
                <div class="form-control mt-2">
                    <label class="label">
                        <span class="label-text">Description (optional)</span>
                    </label>
                    <textarea name="description" placeholder="Enter description..." class="textarea textarea-bordered" rows="3"></textarea>
                </div>
                
                <div class="form-control mt-4">
                    <button type="submit" class="btn btn-primary">Capture</button>
                </div>
            </form>
            
            <div class="modal-action">
                <form method="dialog">
                    <button class="btn">Close</button>
                </form>
            </div>
        </div>
    </dialog>

    <script>
        // Project filtering and sorting
        document.addEventListener('DOMContentLoaded', function() {
            const filterSelect = document.getElementById('project-filter');
            const sortSelect = document.getElementById('project-sort');
            
            if (filterSelect && sortSelect) {
                filterSelect.addEventListener('change', updateProjects);
                sortSelect.addEventListener('change', updateProjects);
            }
            
            function updateProjects() {
                const filter = filterSelect.value;
                const sort = sortSelect.value;
                
                fetch('/api/projects?filter=' + filter + '&sort=' + sort)
                    .then(response => response.json())
                    .then(data => {
                        // Update the projects container with the filtered/sorted projects
                        const container = document.getElementById('projects-container');
                        // This would be replaced with a proper rendering of projects
                        console.log("Projects updated:", data);
                    })
                    .catch(error => {
                        console.error('Error fetching projects:', error);
                    });
            }
        });
    </script>
</body>
</html>
`

	return html
}

// generateProjectCardHtml creates HTML for a single project card
func generateProjectCardHtml(project map[string]interface{}) string {
	id := fmt.Sprintf("%v", project["ID"])
	title := fmt.Sprintf("%v", project["Title"])
	description := fmt.Sprintf("%v", project["Description"])
	status := fmt.Sprintf("%v", project["Status"])
	taskCount := fmt.Sprintf("%v", project["TaskCount"])
	completionPercentage := fmt.Sprintf("%v", project["CompletionPercentage"])
	
	// Format the date from a time.Time value if it exists
	createdAtStr := ""
	if createdAt, ok := project["CreatedAt"].(time.Time); ok {
		createdAtStr = createdAt.Format("Jan 02, 2006")
	}
	
	// Check if due date exists and format it
	dueDateStr := ""
	dueDate, hasDueDate := project["DueDate"].(*time.Time)
	if hasDueDate && dueDate != nil {
		dueDateStr = dueDate.Format("Jan 02, 2006")
	}
	
	card := `
                <div class="card bg-base-100 shadow-md">
                    <div class="card-body p-4">
                        <div class="flex justify-between items-start">
                            <h3 class="card-title">` + title + `</h3>
                            <div class="badge badge-primary">` + status + `</div>
                        </div>
                        
                        <p class="text-sm my-2 line-clamp-2">` + description + `</p>
                        
                        <!-- Project progress -->
                        <div class="mt-3">
                            <div class="flex justify-between mb-1">
                                <span class="text-xs font-medium">Progress</span>
                                <span class="text-xs font-medium">` + completionPercentage + `%</span>
                            </div>
                            <div class="w-full bg-gray-200 rounded-full h-2.5">
                                <div class="bg-primary h-2.5 rounded-full" style="width: ` + completionPercentage + `%"></div>
                            </div>
                        </div>
                        
                        <!-- Project stats -->
                        <div class="flex flex-wrap gap-2 mt-3 text-xs text-gray-500">
                            <div class="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                                </svg>
                                ` + taskCount + ` Tasks
                            </div>
`

	if dueDateStr != "" {
		card += `
                            <div class="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                </svg>
                                Due: ` + dueDateStr + `
                            </div>
`
	}

	if createdAtStr != "" {
		card += `
                            <div class="flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                Created: ` + createdAtStr + `
                            </div>
`
	}

	card += `
                        </div>
                        
                        <!-- Actions -->
                        <div class="card-actions justify-end mt-2">
                            <a href="/projects/` + id + `" class="btn btn-xs btn-outline">View</a>
                            <button class="btn btn-xs btn-outline" onclick="addTaskToProject('` + id + `')">Add Task</button>
                        </div>
                    </div>
                </div>
`

	return card
}