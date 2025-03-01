package models

import (
	"errors"
	"strings"
	"sync"
)

// TaskStore defines the interface for task storage operations
type TaskStore interface {
	Get(id string) (*Task, error)
	GetAll() ([]*Task, error)
	GetAllByUserID(userID string) ([]*Task, error)
	GetByStatus(status TaskStatus) ([]*Task, error)
	GetByStatusAndUserID(status TaskStatus, userID string) ([]*Task, error)
	Search(query string) ([]*Task, error)
	SearchByUserID(query string, userID string) ([]*Task, error)
	Save(task *Task) error
	Delete(id string) error
}

// MemoryTaskStore implements TaskStore interface with in-memory storage
// This is a simple implementation for development - in production you'd use a database
type MemoryTaskStore struct {
	tasks map[string]*Task
	mutex sync.RWMutex
}

// NewMemoryTaskStore creates a new in-memory task store
func NewMemoryTaskStore() *MemoryTaskStore {
	return &MemoryTaskStore{
		tasks: make(map[string]*Task),
	}
}

// Get retrieves a task by ID
func (s *MemoryTaskStore) Get(id string) (*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}

	// Don't return soft-deleted tasks
	if task.IsDeleted() {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// GetAll returns all non-deleted tasks
func (s *MemoryTaskStore) GetAll() ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if !task.IsDeleted() {
			result = append(result, task)
		}
	}

	return result, nil
}

// GetAllByUserID returns all non-deleted tasks for a specific user
func (s *MemoryTaskStore) GetAllByUserID(userID string) ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if !task.IsDeleted() && task.UserID == userID {
			result = append(result, task)
		}
	}

	return result, nil
}

// GetByStatus returns all tasks with the specified status
func (s *MemoryTaskStore) GetByStatus(status TaskStatus) ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if task.Status == status && !task.IsDeleted() {
			result = append(result, task)
		}
	}

	return result, nil
}

// GetByStatusAndUserID returns all tasks with the specified status for a specific user
func (s *MemoryTaskStore) GetByStatusAndUserID(status TaskStatus, userID string) ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if task.Status == status && task.UserID == userID && !task.IsDeleted() {
			result = append(result, task)
		}
	}

	return result, nil
}

// Save creates or updates a task
func (s *MemoryTaskStore) Save(task *Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.tasks[task.ID] = task
	return nil
}

// Delete soft-deletes a task
func (s *MemoryTaskStore) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	task.Delete()
	return nil
}

// Search finds tasks that match the given query in title or description
func (s *MemoryTaskStore) Search(query string) ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task

	// Convert query to lowercase for case-insensitive search
	lowerQuery := strings.ToLower(query)

	for _, task := range s.tasks {
		// Skip deleted tasks
		if task.IsDeleted() {
			continue
		}

		// Check if query is in title or description (case insensitive)
		if strings.Contains(strings.ToLower(task.Title), lowerQuery) ||
			strings.Contains(strings.ToLower(task.Description), lowerQuery) {
			result = append(result, task)
		}

		// Also search in contexts and tags
		for _, context := range task.Contexts {
			if strings.Contains(strings.ToLower(string(context)), lowerQuery) {
				result = append(result, task)
				break // No need to check other contexts
			}
		}

		// Check tags
		for _, tag := range task.Tags {
			if strings.Contains(strings.ToLower(tag), lowerQuery) {
				result = append(result, task)
				break // No need to check other tags
			}
		}
	}

	return result, nil
}

// SearchByUserID finds tasks for a specific user that match the query in title, description, contexts, or tags
func (s *MemoryTaskStore) SearchByUserID(query string, userID string) ([]*Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Task

	// Convert query to lowercase for case-insensitive search
	lowerQuery := strings.ToLower(query)

	for _, task := range s.tasks {
		// Skip deleted tasks and tasks not belonging to the user
		if task.IsDeleted() || task.UserID != userID {
			continue
		}

		// Check if query is in title or description (case insensitive)
		if strings.Contains(strings.ToLower(task.Title), lowerQuery) ||
			strings.Contains(strings.ToLower(task.Description), lowerQuery) {
			result = append(result, task)
		}

		// Also search in contexts and tags
		for _, context := range task.Contexts {
			if strings.Contains(strings.ToLower(string(context)), lowerQuery) {
				result = append(result, task)
				break // No need to check other contexts
			}
		}

		// Check tags
		for _, tag := range task.Tags {
			if strings.Contains(strings.ToLower(tag), lowerQuery) {
				result = append(result, task)
				break // No need to check other tags
			}
		}
	}

	return result, nil
}
