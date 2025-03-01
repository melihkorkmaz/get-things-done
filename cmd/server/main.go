package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/melihkorkmaz/gtd/internal/config"
	"github.com/melihkorkmaz/gtd/internal/handlers"
	"github.com/melihkorkmaz/gtd/internal/models"
	"github.com/melihkorkmaz/gtd/internal/views/pages"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	} else {
		log.Println("Loaded configuration from .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize the task store and user store (PostgreSQL or in-memory)
	var taskStore models.TaskStore
	var userStore models.UserStore
	var err error

	// Check if we should use PostgreSQL
	useDB := os.Getenv("USE_POSTGRES") == "true"

	if useDB {
		// Use PostgreSQL store with configuration
		dbConfig := config.NewDatabaseConfigFromEnv()
		dbConnString := os.Getenv("DATABASE_URL")

		// Allow override via DATABASE_URL if set
		if dbConnString == "" {
			dbConnString = dbConfig.ConnectionString()
		}

		log.Printf("Connecting to PostgreSQL database: %s/%s", dbConfig.Host, dbConfig.DBName)

		// Initialize task store
		pgTaskStore, err := models.NewPgTaskStore(dbConnString)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL for tasks: %v", err)
		}
		defer pgTaskStore.Close()
		taskStore = pgTaskStore

		// Initialize user store
		pgUserStore, err := models.NewPgUserStore(dbConnString)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL for users: %v", err)
		}
		userStore = pgUserStore

		log.Println("Using PostgreSQL database for task and user storage")
	} else {
		// Use in-memory store
		taskStore = models.NewMemoryTaskStore()
		userStore = models.NewMemoryUserStore()
		log.Println("Using in-memory storage (data will be lost when server stops)")

		// Create some sample tasks for testing (only for in-memory store)
		createSampleTasks(taskStore)
	}

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	
	// Initialize auth service
	authConfig := config.NewAuthConfigFromEnv()
	authService := models.NewAuthService(userStore, authConfig.JWTSecret, authConfig.TokenExpiry)
	
	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(authService, authConfig)
	
	// Apply authentication middleware to all routes
	r.Use(authHandler.AuthMiddleware)

	// Serve static files
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "static")
	fileServer(r, "/static", http.Dir(staticDir))

	// Templates directory
	templatesDir := filepath.Join(workDir, "internal/templates")

	// Initialize task handler
	taskHandler, err := handlers.NewTaskHandler(taskStore, templatesDir)
	if err != nil {
		log.Fatalf("Failed to create task handler: %v", err)
	}

	// Initialize project handler
	projectHandler, err := handlers.NewProjectHandler(taskStore, templatesDir)
	if err != nil {
		log.Fatalf("Failed to create project handler: %v", err)
	}

	// Initialize index handler
	indexHandler, err := handlers.NewIndexHandler(taskStore, templatesDir)
	if err != nil {
		log.Fatalf("Failed to create index handler: %v", err)
	}

	// Public routes (don't require authentication)
	
	// Authentication API routes
	authHandler.RegisterRoutes(r)

	// Authentication pages
	r.Get("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login().Render(r.Context(), w)
	})

	r.Get("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		pages.Register().Render(r.Context(), w)
	})
	
	// Public API endpoints
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", handlers.HelloHandler)
	})

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		// Apply authentication middleware to all routes in this group
		r.Use(authHandler.RequireAuthMiddleware)
		
		// Home page
		r.Get("/", indexHandler.HomePage)
		
		// Weekly review page
		r.Get("/weekly-review", indexHandler.WeeklyReviewPage)
		
		// Profile page
		r.Get("/profile", authHandler.ProfilePage)
		
		// Register task routes (all task routes require authentication)
		taskHandler.RegisterRoutes(r)
		taskHandler.RegisterTaskStatusRoutes(r)
		
		// Register project routes (all project routes require authentication)
		projectHandler.RegisterRoutes(r)
	})

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Visit http://localhost:%s to get started\n", port)
	fmt.Printf("Task list available at http://localhost:%s/tasks\n", port)
	fmt.Printf("Inbox available at http://localhost:%s/tasks?status=inbox\n", port)
	log.Fatal(server.ListenAndServe())
}

// createSampleTasks creates a few sample tasks for testing
func createSampleTasks(store models.TaskStore) {
	// Create a sample user ID for demonstration purposes
	sampleUserID := "sample-user-123"
	
	// Inbox item
	task1 := models.NewTask("Capture all open loops", "Gather all tasks, ideas, and commitments into the inbox", sampleUserID)
	store.Save(task1)

	// Next action
	task2 := models.NewTask("Process inbox items", "Go through inbox and decide what to do with each item", sampleUserID)
	task2.MarkAsNext()
	store.Save(task2)

	// Waiting for
	task3 := models.NewTask("Response from email", "Waiting for reply from team about project timeline", sampleUserID)
	task3.MarkAsWaiting()
	store.Save(task3)

	// Someday/maybe
	task4 := models.NewTask("Learn a new language", "Consider learning Spanish or German", sampleUserID)
	task4.MarkAsSomeday()
	store.Save(task4)

	// Project
	task5 := models.NewTask("Redesign personal website", "Project to update and refresh my personal website", sampleUserID)
	task5.MarkAsProject()
	store.Save(task5)
}

// fileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	fs := http.FileServer(root)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the file path from the URL
		filePath := r.URL.Path
		extension := strings.ToLower(filepath.Ext(filePath))

		// Set appropriate MIME types
		switch extension {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		}

		// Serve the file
		http.StripPrefix(strings.TrimSuffix(path, "*"), fs).ServeHTTP(w, r)
	})

	r.Get(path, handler)
}
