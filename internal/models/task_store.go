package models

import (
	"errors"
	"sync"
)

// TaskStore defines the interface for task storage operations
type TaskStore interface {
	Get(id string) (*Task, error)
	GetAll() ([]*Task, error)
	GetByStatus(status TaskStatus) ([]*Task, error)
	Save(task *Task) error
	Delete(id string) error
}

// MemoryTaskStore implements TaskStore interface with in-memory storage
// This is a simple implementation for development - in production you'd use a database
type MemoryTaskStore struct {
	tasks  map[string]*Task
	mutex  sync.RWMutex
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