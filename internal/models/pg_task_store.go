package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgTaskStore implements TaskStore interface with PostgreSQL storage
type PgTaskStore struct {
	db *pgxpool.Pool
}

// NewPgTaskStore creates a new PostgreSQL task store
func NewPgTaskStore(connString string) (*PgTaskStore, error) {
	// Create a connection pool
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Ping database to verify connection
	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	// Create store instance
	store := &PgTaskStore{
		db: db,
	}

	// Initialize database schema
	if err = store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return store, nil
}

// initSchema creates the necessary database tables if they don't exist
func (s *PgTaskStore) initSchema() error {
	// Create tasks table
	_, err := s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			user_id TEXT,
			project_id TEXT,
			parent_id TEXT,
			contexts JSONB,
			tags JSONB,
			due_date TIMESTAMP WITH TIME ZONE,
			scheduled_date TIMESTAMP WITH TIME ZONE,
			time_estimate INTEGER,
			energy_required TEXT,
			priority INTEGER,
			timeframe TEXT,
			is_recurring BOOLEAN DEFAULT FALSE,
			recurring_rule TEXT,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			completed_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
		CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
		CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
	`)

	return err
}

// Close closes the database connection
func (s *PgTaskStore) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Get retrieves a task by ID
func (s *PgTaskStore) Get(id string) (*Task, error) {
	query := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	var task Task
	var contextsJSON, tagsJSON []byte
	var projectID, parentID, energyRequired, timeframe sql.NullString
	var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
	var timeEstimate, priority sql.NullInt32
	var isRecurring sql.NullBool
	var recurringRule sql.NullString

	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
		&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
		&energyRequired, &priority, &timeframe, &isRecurring,
		&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
	)
	
	// Handle nullable string fields
	if projectID.Valid {
		task.ProjectID = projectID.String
	}
	if parentID.Valid {
		task.ParentID = parentID.String
	}
	if energyRequired.Valid {
		task.EnergyRequired = energyRequired.String
	}
	if timeframe.Valid {
		task.Timeframe = Timeframe(timeframe.String)
	}
	if timeEstimate.Valid {
		task.TimeEstimate = int(timeEstimate.Int32)
	}
	if priority.Valid {
		task.Priority = int(priority.Int32)
	}
	if isRecurring.Valid {
		task.IsRecurring = isRecurring.Bool
	}
	if recurringRule.Valid {
		task.RecurringRule = recurringRule.String
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	// Convert JSON fields back to Go structures
	if contextsJSON != nil {
		var contexts []string
		if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
			for _, c := range contexts {
				task.Contexts = append(task.Contexts, Context(c))
			}
		}
	}

	if tagsJSON != nil {
		if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
			// Log the error but continue
			fmt.Printf("Error unmarshaling tags: %v\n", err)
		}
	}

	// Handle nullable time.Time fields
	if dueDate.Valid {
		t := dueDate.Time.Local()
		task.DueDate = &t
	}
	if scheduledDate.Valid {
		t := scheduledDate.Time.Local()
		task.ScheduledDate = &t
	}
	if completedAt.Valid {
		t := completedAt.Time.Local()
		task.CompletedAt = &t
	}
	if deletedAt.Valid {
		t := deletedAt.Time.Local()
		task.DeletedAt = &t
	}

	return &task, nil
}

// GetAll returns all non-deleted tasks
func (s *PgTaskStore) GetAll() ([]*Task, error) {
	query := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				// Log the error but continue
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetAllByUserID returns all non-deleted tasks for a specific user
func (s *PgTaskStore) GetAllByUserID(userID string) ([]*Task, error) {
	query := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				// Log the error but continue
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByStatus returns all tasks with the specified status
func (s *PgTaskStore) GetByStatus(status TaskStatus) ([]*Task, error) {
	query := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(context.Background(), query, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByStatusAndUserID returns all tasks with the specified status for a specific user
func (s *PgTaskStore) GetByStatusAndUserID(status TaskStatus, userID string) ([]*Task, error) {
	query := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE status = $1 AND user_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(context.Background(), query, string(status), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Save creates or updates a task
func (s *PgTaskStore) Save(task *Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	// Prepare contexts and tags for JSON storage
	var contextsSlice []string
	for _, c := range task.Contexts {
		contextsSlice = append(contextsSlice, string(c))
	}

	contextsJSON, err := json.Marshal(contextsSlice)
	if err != nil {
		return err
	}

	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return err
	}

	// Ensure task has an updated timestamp
	task.UpdatedAt = time.Now()

	query := `
		INSERT INTO tasks (
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		) ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			user_id = EXCLUDED.user_id,
			project_id = EXCLUDED.project_id,
			parent_id = EXCLUDED.parent_id,
			contexts = EXCLUDED.contexts,
			tags = EXCLUDED.tags,
			due_date = EXCLUDED.due_date,
			scheduled_date = EXCLUDED.scheduled_date,
			time_estimate = EXCLUDED.time_estimate,
			energy_required = EXCLUDED.energy_required,
			priority = EXCLUDED.priority,
			timeframe = EXCLUDED.timeframe,
			is_recurring = EXCLUDED.is_recurring,
			recurring_rule = EXCLUDED.recurring_rule,
			updated_at = EXCLUDED.updated_at,
			completed_at = EXCLUDED.completed_at,
			deleted_at = EXCLUDED.deleted_at
	`

	_, err = s.db.Exec(context.Background(), query,
		task.ID, task.Title, task.Description, string(task.Status), task.UserID, task.ProjectID, task.ParentID,
		contextsJSON, tagsJSON, task.DueDate, task.ScheduledDate, task.TimeEstimate,
		task.EnergyRequired, task.Priority, string(task.Timeframe), task.IsRecurring,
		task.RecurringRule, task.CreatedAt, task.UpdatedAt, task.CompletedAt, task.DeletedAt,
	)

	return err
}

// Delete soft-deletes a task
func (s *PgTaskStore) Delete(id string) error {
	// First check if task exists
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	// Soft delete the task
	task.Delete()
	return s.Save(task)
}

// Search finds tasks that match the query in title, description, contexts, or tags
func (s *PgTaskStore) Search(query string) ([]*Task, error) {
	// Build a query that searches in multiple columns with case-insensitive matching
	sqlQuery := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND (
			LOWER(title) LIKE LOWER($1) OR 
			LOWER(description) LIKE LOWER($1) OR
			contexts::text ILIKE $1 OR
			tags::text ILIKE $1
		)
		ORDER BY created_at DESC
	`

	// Add wildcards for partial matching
	searchPattern := "%" + query + "%"

	rows, err := s.db.Query(context.Background(), sqlQuery, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SearchByUserID finds tasks for a specific user that match the query in title, description, contexts, or tags
func (s *PgTaskStore) SearchByUserID(query string, userID string) ([]*Task, error) {
	// Build a query that searches in multiple columns with case-insensitive matching
	sqlQuery := `
		SELECT 
			id, title, description, status, user_id, project_id, parent_id, 
			contexts, tags, due_date, scheduled_date, time_estimate, 
			energy_required, priority, timeframe, is_recurring, 
			recurring_rule, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND user_id = $2 AND (
			LOWER(title) LIKE LOWER($1) OR 
			LOWER(description) LIKE LOWER($1) OR
			contexts::text ILIKE $1 OR
			tags::text ILIKE $1
		)
		ORDER BY created_at DESC
	`

	// Add wildcards for partial matching
	searchPattern := "%" + query + "%"

	rows, err := s.db.Query(context.Background(), sqlQuery, searchPattern, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var contextsJSON, tagsJSON []byte
		var projectID, parentID, energyRequired, timeframe sql.NullString
		var dueDate, scheduledDate, completedAt, deletedAt pgtype.Timestamptz
		var timeEstimate, priority sql.NullInt32
		var isRecurring sql.NullBool
		var recurringRule sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &projectID, &parentID,
			&contextsJSON, &tagsJSON, &dueDate, &scheduledDate, &timeEstimate,
			&energyRequired, &priority, &timeframe, &isRecurring,
			&recurringRule, &task.CreatedAt, &task.UpdatedAt, &completedAt, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable string fields
		if projectID.Valid {
			task.ProjectID = projectID.String
		}
		if parentID.Valid {
			task.ParentID = parentID.String
		}
		if energyRequired.Valid {
			task.EnergyRequired = energyRequired.String
		}
		if timeframe.Valid {
			task.Timeframe = Timeframe(timeframe.String)
		}
		if timeEstimate.Valid {
			task.TimeEstimate = int(timeEstimate.Int32)
		}
		if priority.Valid {
			task.Priority = int(priority.Int32)
		}
		if isRecurring.Valid {
			task.IsRecurring = isRecurring.Bool
		}
		if recurringRule.Valid {
			task.RecurringRule = recurringRule.String
		}

		// Convert JSON fields back to Go structures
		if contextsJSON != nil {
			var contexts []string
			if err := json.Unmarshal(contextsJSON, &contexts); err == nil {
				for _, c := range contexts {
					task.Contexts = append(task.Contexts, Context(c))
				}
			}
		}

		if tagsJSON != nil {
			if err := json.Unmarshal(tagsJSON, &task.Tags); err != nil {
				fmt.Printf("Error unmarshaling tags: %v\n", err)
			}
		}

		// Handle nullable time.Time fields
		if dueDate.Valid {
			t := dueDate.Time.Local()
			task.DueDate = &t
		}
		if scheduledDate.Valid {
			t := scheduledDate.Time.Local()
			task.ScheduledDate = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.Local()
			task.CompletedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time.Local()
			task.DeletedAt = &t
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
