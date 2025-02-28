package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/melihkorkmaz/gtd/internal/models"
)

// StatusChangeResponse is the response for status change endpoints
type StatusChangeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	TaskID  string `json:"taskId"`
	Status  string `json:"status"`
}

// RegisterTaskStatusRoutes registers routes for task status transitions
func (h *TaskHandler) RegisterTaskStatusRoutes(r chi.Router) {
	r.Put("/api/tasks/{id}/next", h.MarkTaskAsNext)
	r.Put("/api/tasks/{id}/waiting", h.MarkTaskAsWaiting)
	r.Put("/api/tasks/{id}/someday", h.MarkTaskAsSomeday)
	r.Put("/api/tasks/{id}/done", h.MarkTaskAsDone)
	r.Put("/api/tasks/{id}/project", h.MarkTaskAsProject)
	r.Put("/api/tasks/{id}/scheduled", h.ScheduleTask)
}

// MarkTaskAsNext marks a task as a next action
func (h *TaskHandler) MarkTaskAsNext(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.MarkAsNext()
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task marked as Next Action")
}

// MarkTaskAsWaiting marks a task as waiting for someone else
func (h *TaskHandler) MarkTaskAsWaiting(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.MarkAsWaiting()
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task marked as Waiting For")
}

// MarkTaskAsSomeday marks a task as a someday/maybe item
func (h *TaskHandler) MarkTaskAsSomeday(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.MarkAsSomeday()
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task marked as Someday/Maybe")
}

// MarkTaskAsDone marks a task as done
func (h *TaskHandler) MarkTaskAsDone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.MarkAsDone()
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task marked as Done")
}

// MarkTaskAsProject marks a task as a project
func (h *TaskHandler) MarkTaskAsProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.MarkAsProject()
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task converted to Project")
}

// ScheduleRequest represents the request to schedule a task
type ScheduleRequest struct {
	Date string `json:"date"` // Format: "2006-01-02"
}

// ScheduleTask schedules a task for a specific date
func (h *TaskHandler) ScheduleTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scheduleDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	task.MarkAsScheduled(scheduleDate)
	if err := h.store.Save(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendStatusChangeResponse(w, task, "Task scheduled")
}

// Helper function to send a standard response for status changes
func sendStatusChangeResponse(w http.ResponseWriter, task *models.Task, message string) {
	w.Header().Set("Content-Type", "application/json")
	resp := StatusChangeResponse{
		Success: true,
		Message: message,
		TaskID:  task.ID,
		Status:  string(task.Status),
	}
	json.NewEncoder(w).Encode(resp)
}
