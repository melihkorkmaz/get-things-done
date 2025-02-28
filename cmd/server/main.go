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

	// Initialize the task store (PostgreSQL or in-memory)
	var taskStore models.TaskStore
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
		
		pgStore, err := models.NewPgTaskStore(dbConnString)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		defer pgStore.Close()
		
		taskStore = pgStore
		log.Println("Using PostgreSQL database for task storage")
	} else {
		// Use in-memory store
		taskStore = models.NewMemoryTaskStore()
		log.Println("Using in-memory storage for tasks (data will be lost when server stops)")
		
		// Create some sample tasks for testing (only for in-memory store)
		createSampleTasks(taskStore)
	}

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

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
	
	// Register task routes
	taskHandler.RegisterRoutes(r)
	taskHandler.RegisterTaskStatusRoutes(r)
	
	// Initialize project handler
	projectHandler, err := handlers.NewProjectHandler(taskStore, templatesDir)
	if err != nil {
		log.Fatalf("Failed to create project handler: %v", err)
	}
	
	// Register project routes
	projectHandler.RegisterRoutes(r)

	// Initialize index handler
	indexHandler, err := handlers.NewIndexHandler(taskStore, templatesDir)
	if err != nil {
		log.Fatalf("Failed to create index handler: %v", err)
	}
	
	// Home page
	r.Get("/", indexHandler.HomePage)
	
	// Weekly review page
	r.Get("/weekly-review", indexHandler.WeeklyReviewPage)
	
	// API routes for hello example (from original setup)
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", handlers.HelloHandler)
	})

	// Debug: Print registered routes
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("Route registered: %s %s\n", method, route)
		return nil
	}
	
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Error walking routes: %v\n", err)
	}

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
	// Inbox item
	task1 := models.NewTask("Capture all open loops", "Gather all tasks, ideas, and commitments into the inbox")
	store.Save(task1)
	
	// Next action
	task2 := models.NewTask("Process inbox items", "Go through inbox and decide what to do with each item")
	task2.MarkAsNext()
	store.Save(task2)
	
	// Waiting for
	task3 := models.NewTask("Response from email", "Waiting for reply from team about project timeline")
	task3.MarkAsWaiting()
	store.Save(task3)
	
	// Someday/maybe
	task4 := models.NewTask("Learn a new language", "Consider learning Spanish or German")
	task4.MarkAsSomeday()
	store.Save(task4)
	
	// Project
	task5 := models.NewTask("Redesign personal website", "Project to update and refresh my personal website")
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