# GTD App

A fullstack Go application built with:

- Go (backend)
- Chi router (HTTP routing)
- HTMX (dynamic HTML without writing JavaScript)
- Alpine.js (lightweight JavaScript framework)
- DaisyUI (Tailwind CSS component library)

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone this repository
```bash
git clone https://github.com/melihkorkmaz/gtd.git
cd gtd
```

2. Install dependencies
```bash
go mod tidy
```

3. Choose a storage backend:

   **Option 1: In-memory storage**
   ```bash
   # Edit .env file and set USE_POSTGRES=false
   # Then run:
   go run cmd/server/main.go
   ```

   **Option 2: PostgreSQL with Docker (Recommended)**
   ```bash
   # Start PostgreSQL container
   docker compose up -d
   
   # Run the application (it will use settings from .env file)
   go run cmd/server/main.go
   ```

   **Option 3: PostgreSQL (Manual setup)**
   ```bash
   # First, make sure PostgreSQL is running
   # Create a database for the application
   createdb gtd

   # Configure .env file with your database settings
   # Then run:
   go run cmd/server/main.go
   ```

   **Configuration**
   The application uses a `.env` file for configuration with these default settings:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=gtd
   DB_SSL_MODE=disable
   PORT=3000
   USE_POSTGRES=true
   ```

   You can override these settings with environment variables:
   ```bash
   DB_HOST=custom-host DB_PASSWORD=custom-password go run cmd/server/main.go
   ```

   Alternatively, you can set the complete DATABASE_URL:
   ```bash
   DATABASE_URL=postgres://user:password@localhost:5432/gtd?sslmode=disable go run cmd/server/main.go
   ```

4. Open your browser and navigate to `http://localhost:3000`

## Project Structure

```
.
├── cmd
│   └── server        # Main application entry point
├── internal
│   ├── handlers      # HTTP request handlers
│   ├── models        # Domain models
│   └── templates     # HTML templates
└── static            # Static assets
    ├── css           # CSS files
    └── js            # JavaScript files
```

## Features

- Task management (GTD methodology)
- Modern UI with DaisyUI and Tailwind CSS
- Interactive UI with minimal JavaScript using HTMX and Alpine.js

## Built With

- [Go](https://golang.org/)
- [Chi Router](https://github.com/go-chi/chi)
- [HTMX](https://htmx.org/)
- [Alpine.js](https://alpinejs.dev/)
- [DaisyUI](https://daisyui.com/)
- [Tailwind CSS](https://tailwindcss.com/)