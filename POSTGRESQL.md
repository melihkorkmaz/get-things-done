# PostgreSQL Database Setup for GTD App

This document provides information about setting up and using PostgreSQL with the GTD application.

## Database Setup

### Option 1: Docker (Recommended)

The simplest way to set up PostgreSQL is using the included Docker Compose file:

1. **Start PostgreSQL container**:
   ```bash
   docker compose up -d
   ```
   
   This will start a PostgreSQL 16 container with the following configuration:
   - Username: `postgres`
   - Password: `postgres`
   - Database: `gtd`
   - Port: `5432` (mapped to your localhost)
   - Data is persisted in a Docker volume named `pg_data`

2. **Run the application**:
   The application is already configured to connect to this PostgreSQL container via the `.env` file.
   ```bash
   go run cmd/server/main.go
   ```

### Option 2: Manual PostgreSQL Setup

1. **Install PostgreSQL**:
   - Mac: `brew install postgresql` and `brew services start postgresql`
   - Ubuntu: `sudo apt-get install postgresql postgresql-contrib`
   - Windows: Download the installer from the [PostgreSQL website](https://www.postgresql.org/download/windows/)

2. **Create a database**:
   ```bash
   createdb gtd
   ```

3. **Configure the application**:
   Edit the `.env` file in the project root:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres  # Change to your PostgreSQL username
   DB_PASSWORD=postgres  # Change to your PostgreSQL password
   DB_NAME=gtd
   DB_SSL_MODE=disable
   USE_POSTGRES=true
   ```

   Alternatively, set the DATABASE_URL environment variable:
   ```bash
   export DATABASE_URL=postgres://user:password@localhost:5432/gtd?sslmode=disable
   ```

## Schema Information

The application creates the following schema in PostgreSQL:

### Tasks Table

```sql
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
```

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT,
    provider TEXT NOT NULL,
    provider_id TEXT,
    picture TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_login_at TIMESTAMP WITH TIME ZONE,
    activated_at TIMESTAMP WITH TIME ZONE,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_provider_provider_id ON users(provider, provider_id) 
WHERE provider <> 'email' AND provider_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);
```

## Database Operations

### Querying Tasks

```sql
-- Get all active tasks
SELECT * FROM tasks WHERE deleted_at IS NULL;

-- Get all tasks for a specific user
SELECT * FROM tasks WHERE user_id = '12345' AND deleted_at IS NULL;

-- Get all tasks in the inbox for a specific user
SELECT * FROM tasks WHERE status = 'inbox' AND user_id = '12345' AND deleted_at IS NULL;

-- Get all next actions for a specific user
SELECT * FROM tasks WHERE status = 'next' AND user_id = '12345' AND deleted_at IS NULL;

-- Get completed tasks for a specific user
SELECT * FROM tasks WHERE status = 'done' AND user_id = '12345' AND deleted_at IS NULL;
```

### Working with Contexts and Tags

The `contexts` and `tags` fields are stored as JSONB arrays. You can query them like this:

```sql
-- Find tasks with a specific context for a user
SELECT * FROM tasks WHERE contexts @> '["home"]'::jsonb AND user_id = '12345' AND deleted_at IS NULL;

-- Find tasks with a specific tag for a user
SELECT * FROM tasks WHERE tags @> '["important"]'::jsonb AND user_id = '12345' AND deleted_at IS NULL;

-- Find tasks with specific text in title or description for a user
SELECT * FROM tasks 
WHERE (title ILIKE '%search term%' OR description ILIKE '%search term%') 
  AND user_id = '12345' AND deleted_at IS NULL;
```

### Working with Users

```sql
-- Get all active users
SELECT * FROM users WHERE active = TRUE;

-- Find a user by email
SELECT * FROM users WHERE LOWER(email) = LOWER('user@example.com');

-- Find a user by social provider
SELECT * FROM users WHERE provider = 'google' AND provider_id = '12345';

-- Create a test user with password 'password123'
INSERT INTO users (
  id, first_name, last_name, email, password, 
  provider, active, created_at, updated_at, activated_at
) VALUES (
  '202503010001', 'Test', 'User', 'test@example.com', 
  '$2a$10$Q7K5UHkvhvECFHqLNGD4GuZOySKIoYV5MmUvZ3FNOoT9ZqmWrO0p6', 
  'email', true, NOW(), NOW(), NOW()
);

-- Get tasks for a specific user
SELECT * FROM tasks WHERE user_id = '202503010001' AND deleted_at IS NULL;
```

### Backup and Restore

To backup the database:
```bash
pg_dump -U postgres gtd > gtd_backup.sql
```

To restore from a backup:
```bash
psql -U postgres -d gtd -f gtd_backup.sql
```