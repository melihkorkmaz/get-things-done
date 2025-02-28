package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
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
	// Log when this handler is called
	fmt.Printf("ListTasksPage handler called with path: %s and query: %s\n", r.URL.Path, r.URL.RawQuery)

	// Get status filter if provided
	status := r.URL.Query().Get("status")

	var tasks []*models.Task
	var err error
	var title string

	if status != "" {
		// Filter tasks by status
		fmt.Printf("Filtering tasks by status: %s\n", status)
		tasks, err = h.store.GetByStatus(models.TaskStatus(status))

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
		// Get all tasks
		fmt.Println("Getting all tasks")
		tasks, err = h.store.GetAll()
		title = "All Tasks"
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Found %d tasks\n", len(tasks))

	// ********* TEMPORARY DIRECT HTML APPROACH *********
	// This bypasses the template system to test if the route is working

	w.Header().Set("Content-Type", "text/html")

	taskListHtml := `
<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + title + ` - GTD App</title>
    
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
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
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
                    <li><a href="/tasks?status=someday">Someday/Maybe</a></li>
                    <li><a href="/weekly-review" class="text-accent">Weekly Review</a></li>
                </ul>
            </div>
        </header>

    <div class="container mx-auto">
        <div class="bg-base-100 p-6 rounded shadow">
            <div class="flex justify-between items-center mb-6">
                <h1 class="text-2xl font-bold">` + title + `</h1>
                <div class="flex items-center gap-4">
                    <span class="badge">` + fmt.Sprintf("%d tasks found", len(tasks)) + `</span>
                    <button class="btn btn-primary" onclick="document.getElementById('quick-capture-modal').showModal()">
                        Add Task
                    </button>
                </div>
            </div>
            
            <!-- Quick Capture Modal -->
            <dialog id="quick-capture-modal" class="modal">
                <div class="modal-box">
                    <h3 class="font-bold text-lg">Quick Capture</h3>
                    <p class="py-2">Quickly capture a new task or idea.</p>
                    
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
            
            <div class="space-y-4">
`

	if len(tasks) > 0 {
		for _, task := range tasks {
			taskListHtml += `
            <div class="bg-base-200 p-4 rounded">
                <div class="flex justify-between">
                    <h3 class="font-bold">` + task.Title + `</h3>
                    <span class="badge">` + string(task.Status) + `</span>
                </div>
                <p class="text-sm my-2">` + task.Description + `</p>
                <div class="flex justify-end mt-2">
                    <a href="/tasks/` + task.ID + `" class="btn btn-xs btn-outline">View</a>
                </div>
            </div>`
		}
	} else {
		taskListHtml += `
            <div class="alert">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <span>No tasks found. Create one using the Add Task button.</span>
            </div>`
	}

	taskListHtml += `
            </div>
        </div>
    </div>
</body>
</html>
`
	w.Write([]byte(taskListHtml))

	// ********* ORIGINAL APPROACH - COMMENTED OUT *********
	/*
		data := map[string]interface{}{
			"Title": title,
			"Tasks": tasks,
		}

		w.Header().Set("Content-Type", "text/html")

		// Print debug information
		fmt.Printf("DEBUG: Executing template 'base.html' for ListTasksPage with title: %s and %d tasks\n", title, len(tasks))

		// Debug info about defined templates
		for _, tmpl := range h.templates.templates.Templates() {
			fmt.Printf("DEBUG: Found template: %s\n", tmpl.Name())
		}

		// Execute the template and check for errors
		err = h.templates.templates.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			fmt.Printf("ERROR: Failed to execute template: %v\n", err)
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	*/
}

// NewTaskForm renders the form to create a new task
func (h *TaskHandler) NewTaskForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	h.templates.templates.ExecuteTemplate(w, "task_form", nil)
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

	// Handle contexts
	for _, ctx := range r.Form["contexts[]"] {
		task.Contexts = append(task.Contexts, models.Context(ctx))
	}

	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return a partial HTML response
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		h.templates.templates.ExecuteTemplate(w, "task_row", task)
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

	// ********* TEMPORARY DIRECT HTML APPROACH *********
	// This bypasses the template system to test if the route is working

	w.Header().Set("Content-Type", "text/html")

	taskHtml := `
<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + task.Title + ` - GTD App</title>
    
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
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
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
                    <li><a href="/tasks?status=someday">Someday/Maybe</a></li>
                    <li><a href="/weekly-review" class="text-accent">Weekly Review</a></li>
                </ul>
            </div>
        </header>

    <div class="container mx-auto">
        <div class="bg-base-100 p-6 rounded shadow">
            <div class="flex justify-between items-center mb-6">
                <h1 class="text-2xl font-bold">` + task.Title + `</h1>
                <span class="badge badge-` + string(task.Status) + `">` + string(task.Status) + `</span>
            </div>
            
            <div class="mb-6">
                <p class="text-lg">` + task.Description + `</p>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <div class="card bg-base-200 p-4">
                    <h3 class="font-bold text-lg mb-3">Details</h3>
                    <div class="divider my-1"></div>
                    <div class="flex flex-col gap-2">
                        <div class="flex justify-between">
                            <span class="font-medium">Created:</span>
                            <span>` + task.CreatedAt.Format("Jan 02, 2006 3:04 PM") + `</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="font-medium">Status:</span>
                            <span class="badge badge-` + string(task.Status) + `">` + string(task.Status) + `</span>
                        </div>
                    </div>
                </div>
                
                <div class="card bg-base-200 p-4">
                    <h3 class="font-bold text-lg mb-3">Actions</h3>
                    <div class="divider my-1"></div>
                    
                    <div id="taskActions">
                        <!-- Actions will be populated by JavaScript based on task status -->
                    </div>
                    
                    <script>
                        // Set up task actions based on status when page loads
                        document.addEventListener('DOMContentLoaded', function() {
                            const taskStatus = '` + string(task.Status) + `';
                            const taskId = '` + task.ID + `';
                            const actionsContainer = document.getElementById('taskActions');
                            
                            if (taskStatus === 'inbox') {
                                actionsContainer.innerHTML = 
                                    '<div class="dropdown">' +
                                    '    <label tabindex="0" class="btn btn-primary m-1">Process Task</label>' +
                                    '    <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">' +
                                    '        <li><a href="#" onclick="processTask(\'' + taskId + '\', \'next\')">Mark as Next Action</a></li>' +
                                    '        <li><a href="#" onclick="processTask(\'' + taskId + '\', \'waiting\')">Mark as Waiting For</a></li>' +
                                    '        <li><a href="#" onclick="processTask(\'' + taskId + '\', \'someday\')">Mark as Someday/Maybe</a></li>' +
                                    '        <li><a href="#" onclick="processTask(\'' + taskId + '\', \'project\')">Convert to Project</a></li>' +
                                    '    </ul>' +
                                    '</div>';
                            } else if (taskStatus !== 'done') {
                                actionsContainer.innerHTML =
                                    '<button class="btn btn-success" onclick="processTask(\'' + taskId + '\', \'done\')">' +
                                    '    Mark as Done' +
                                    '</button>';
                            } else {
                                actionsContainer.innerHTML = 
                                    '<div class="alert alert-success">' +
                                    '    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>' +
                                    '    <span>This task has been marked as done!</span>' +
                                    '</div>';
                            }
                        });
                    </script>
                </div>
            </div>
            
            <div class="flex space-x-2 mt-6">
                <a href="/tasks" class="btn">Back to Tasks</a>
                <a href="/tasks/` + task.ID + `/edit" class="btn btn-primary">Edit Task</a>
                <button class="btn btn-error" onclick="deleteTask('` + task.ID + `')">Delete Task</button>
            </div>
            
            <!-- Quick Capture Modal -->
            <dialog id="quick-capture-modal" class="modal">
                <div class="modal-box">
                    <h3 class="font-bold text-lg">Quick Capture</h3>
                    <p class="py-2">Quickly capture a new task or idea.</p>
                    
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
        </div>
    </div>
    
    <!-- JavaScript for task actions -->
    <script>
        function processTask(taskId, status) {
            if (!confirm('Are you sure you want to change this task to ' + status + '?')) {
                return;
            }
            
            fetch('/api/tasks/' + taskId + '/' + status, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => {
                if (response.ok) {
                    window.location.reload();
                } else {
                    alert('Failed to update task status');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while updating the task');
            });
        }
        
        function deleteTask(taskId) {
            if (!confirm('Are you sure you want to delete this task? This cannot be undone.')) {
                return;
            }
            
            fetch('/api/tasks/' + taskId, {
                method: 'DELETE'
            })
            .then(response => {
                if (response.ok) {
                    window.location.href = '/tasks';
                } else {
                    alert('Failed to delete task');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while deleting the task');
            });
        }
    </script>
</body>
</html>
`
	w.Write([]byte(taskHtml))

	// ********* ORIGINAL APPROACH - COMMENTED OUT *********
	/*
		data := map[string]interface{}{
			"Title": task.Title,
			"Task":  task,
		}

		w.Header().Set("Content-Type", "text/html")
		h.templates.templates.ExecuteTemplate(w, "base.html", data)
	*/
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
	h.templates.templates.ExecuteTemplate(w, "task_form", data)
}

// SearchTasksAPI returns a JSON list of tasks matching the search query
func (h *TaskHandler) SearchTasksAPI(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	tasks, err := h.store.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// SearchTasksPage renders the search results page for tasks
func (h *TaskHandler) SearchTasksPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var tasks []*models.Task
	var err error

	if query != "" {
		tasks, err = h.store.Search(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data := map[string]interface{}{
		"Title":        "Search Results",
		"Tasks":        tasks,
		"SearchQuery":  query,
		"ResultsCount": len(tasks),
	}

	// If this is an HTMX request, return just the search results
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		h.templates.templates.ExecuteTemplate(w, "search_results", data)
		return
	}

	// Otherwise return the full page
	w.Header().Set("Content-Type", "text/html")
	h.templates.templates.ExecuteTemplate(w, "base.html", data)
}

// QuickCaptureAPI handles the quick capture task creation
func (h *TaskHandler) QuickCaptureAPI(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Create a new task
	task := models.NewTask(title, description)

	// Save to the store
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return HTML response for HTMX
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		successHTML := `
		<div class="alert alert-success shadow-lg">
			<div>
				<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
				<div>
					<span class="font-bold">Task captured successfully!</span>
					<p class="text-sm">Added to your inbox for processing later.</p>
				</div>
			</div>
		</div>
		`
		w.Write([]byte(successHTML))
		return
	}

	// JSON response if not HTMX
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Task captured successfully",
		"task":    task,
	})
}
