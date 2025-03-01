# GTD App

A fullstack Go application built with:

- Go (backend)
- Chi router (HTTP routing)
- HTMX (dynamic HTML without writing JavaScript)
- Alpine.js (lightweight JavaScript framework)
- DaisyUI (Tailwind CSS component library)

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Templ cli (for template generation)

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
   # Database configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=gtd
   DB_SSL_MODE=disable
   USE_POSTGRES=true

   # Authentication
   JWT_SECRET=your-256-bit-secret-key-change-this-in-production
   COOKIE_DOMAIN=localhost
   COOKIE_SECURE=false

   # Google OAuth
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   GOOGLE_REDIRECT_URL=http://localhost:3000/auth/google/callback

   # Server configuration
   PORT=3000
   ```
   
   You should copy .env.example to .env and customize the values for your environment.

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
│   ├── config        # Application configuration
│   ├── handlers      # HTTP request handlers
│   ├── models        # Domain models
│   ├── templates     # HTML templates (legacy, using template.html)
│   └── views         # Templ templates (new templating system)
│       ├── layouts   # Base layout templates
│       ├── pages     # Page templates
│       └── partials  # Reusable component templates
└── static            # Static assets
    ├── css           # CSS files
    └── js            # JavaScript files
```

## Features

- Complete GTD (Getting Things Done) methodology implementation:
  - Capture: Quick capture forms accessible from anywhere
  - Clarify: Process inbox items into actionable tasks
  - Organize: Projects, contexts, tags, and status organization
  - Reflect: Weekly review features and dashboards
  - Engage: Context-based filtering and prioritization
- User authentication system:
  - Traditional email/password authentication
  - Social login with Google OAuth
  - JWT-based authentication with secure cookies
  - User profile management
  - Password reset functionality
  - Data isolation: Users can only see and manage their own tasks and projects
- Project management with task relationships and progress tracking
- Advanced task filtering by status, context, and tags
- Modern UI with DaisyUI Bumblebee theme and Tailwind CSS
- Interactive UI with minimal JavaScript using HTMX and Alpine.js
- PostgreSQL database integration for persistence with proper handling of NULL values

## Built With

- [Go](https://golang.org/)
- [Chi Router](https://github.com/go-chi/chi)
- [Templ](https://templ.guide/) - Go HTML templating language
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Token implementation
- [Bcrypt](https://golang.org/x/crypto/bcrypt) - Password hashing
- [OAuth2](https://golang.org/x/oauth2) - OAuth2 client implementation
- [HTMX](https://htmx.org/) - HTML-based AJAX for modern web apps
- [Alpine.js](https://alpinejs.dev/) - Minimal JavaScript framework
- [DaisyUI](https://daisyui.com/) - Component library for Tailwind CSS, using Bumblebee theme
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [PostgreSQL](https://www.postgresql.org/) - Relational database

