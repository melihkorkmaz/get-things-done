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

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
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
```

## Database Operations

### Querying Tasks

```sql
-- Get all active tasks
SELECT * FROM tasks WHERE deleted_at IS NULL;

-- Get all tasks in the inbox
SELECT * FROM tasks WHERE status = 'inbox' AND deleted_at IS NULL;

-- Get all next actions
SELECT * FROM tasks WHERE status = 'next' AND deleted_at IS NULL;

-- Get completed tasks
SELECT * FROM tasks WHERE status = 'done' AND deleted_at IS NULL;
```

### Working with Contexts and Tags

The `contexts` and `tags` fields are stored as JSONB arrays. You can query them like this:

```sql
-- Find tasks with a specific context
SELECT * FROM tasks WHERE contexts @> '["home"]'::jsonb;

-- Find tasks with a specific tag
SELECT * FROM tasks WHERE tags @> '["important"]'::jsonb;
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