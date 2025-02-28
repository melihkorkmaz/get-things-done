package models

import (
	"errors"
	"time"
)

// TaskStatus defines possible statuses for a task
type TaskStatus string

const (
	// Status constants for GTD workflow
	StatusInbox    TaskStatus = "inbox"    // Uncategorized/unprocessed
	StatusNext     TaskStatus = "next"     // Next actions
	StatusWaiting  TaskStatus = "waiting"  // Waiting for someone else
	StatusScheduled TaskStatus = "scheduled" // Scheduled for a specific time
	StatusSomeday  TaskStatus = "someday"  // Someday/maybe items
	StatusDone     TaskStatus = "done"     // Completed tasks
	StatusProject  TaskStatus = "project"  // Multi-step projects
	StatusReference TaskStatus = "reference" // Reference material
)

// Context represents where a task can be completed
type Context string

// Timeframe represents when a task should be done
type Timeframe string

const (
	TimeframeToday     Timeframe = "today"
	TimeframeThisWeek  Timeframe = "this_week"
	TimeframeNextWeek  Timeframe = "next_week"
	TimeframeSomeday   Timeframe = "someday"
)

// Task represents a GTD task
type Task struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         TaskStatus `json:"status"`
	ProjectID      string     `json:"projectId,omitempty"`     // For tasks that are part of a project
	ParentID       string     `json:"parentId,omitempty"`      // For hierarchical tasks
	Contexts       []Context  `json:"contexts,omitempty"`      // Where this can be done (home, work, phone, etc.)
	Tags           []string   `json:"tags,omitempty"`          // Custom tags for organization
	DueDate        *time.Time `json:"dueDate,omitempty"`       // When this must be completed by
	ScheduledDate  *time.Time `json:"scheduledDate,omitempty"` // When this is scheduled to be done
	TimeEstimate   int        `json:"timeEstimate,omitempty"`  // Estimated minutes to complete
	EnergyRequired string     `json:"energyRequired,omitempty"` // High, medium, low
	Priority       int        `json:"priority,omitempty"`      // 1-3 priority level (1 highest)
	Timeframe      Timeframe  `json:"timeframe,omitempty"`     // When this should be addressed
	IsRecurring    bool       `json:"isRecurring,omitempty"`   // Whether this task recurs
	RecurringRule  string     `json:"recurringRule,omitempty"` // Rule for recurrence (e.g., "daily", "weekly on Monday")
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`     // Soft delete support
}

// NewTask creates a new task with default values (in inbox)
func NewTask(title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          GenerateID(),
		Title:       title,
		Description: description,
		Status:      StatusInbox,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate checks if the task data is valid
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("task title cannot be empty")
	}
	
	// Validate status is one of the defined constants
	validStatus := map[TaskStatus]bool{
		StatusInbox:     true,
		StatusNext:      true,
		StatusWaiting:   true,
		StatusScheduled: true,
		StatusSomeday:   true,
		StatusDone:      true,
		StatusProject:   true,
		StatusReference: true,
	}
	
	if !validStatus[t.Status] {
		return errors.New("invalid task status")
	}
	
	return nil
}

// MarkAsNext moves a task to the Next Actions list
func (t *Task) MarkAsNext() {
	t.Status = StatusNext
	t.UpdatedAt = time.Now()
}

// MarkAsWaiting marks a task as waiting for someone else
func (t *Task) MarkAsWaiting() {
	t.Status = StatusWaiting
	t.UpdatedAt = time.Now()
}

// MarkAsScheduled schedules a task for a specific time
func (t *Task) MarkAsScheduled(scheduledDate time.Time) {
	t.Status = StatusScheduled
	t.ScheduledDate = &scheduledDate
	t.UpdatedAt = time.Now()
}

// MarkAsSomeday moves a task to the Someday/Maybe list
func (t *Task) MarkAsSomeday() {
	t.Status = StatusSomeday
	t.UpdatedAt = time.Now()
}

// MarkAsDone marks a task as completed
func (t *Task) MarkAsDone() {
	t.Status = StatusDone
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkAsProject marks a task as a project (container for other tasks)
func (t *Task) MarkAsProject() {
	t.Status = StatusProject
	t.UpdatedAt = time.Now()
}

// Delete soft-deletes a task
func (t *Task) Delete() {
	now := time.Now()
	t.DeletedAt = &now
	t.UpdatedAt = now
}

// IsDeleted checks if a task has been soft-deleted
func (t *Task) IsDeleted() bool {
	return t.DeletedAt != nil
}

// GenerateID generates a random ID for a task
// Note: In a real application, you might want to use a more sophisticated ID generation method
func GenerateID() string {
	return time.Now().Format("20060102150405")
}