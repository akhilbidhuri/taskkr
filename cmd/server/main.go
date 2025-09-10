package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/akhilbidhuri/taskkr/internal/config"
	"github.com/akhilbidhuri/taskkr/internal/handler"
	appMiddleware "github.com/akhilbidhuri/taskkr/internal/middleware"
	"github.com/akhilbidhuri/taskkr/internal/repository"
	"github.com/akhilbidhuri/taskkr/internal/repository/postgres"
	"github.com/akhilbidhuri/taskkr/internal/service"
)

var cfg *config.Config
var db *gorm.DB
var taskRepo repository.TaskRepository
var taskService *service.TaskService
var taskHandler *handler.TaskHandler

func initialize() {
	log.Println("init method run")
	// Load configuration
	cfg = config.Load()

	// Initialize database
	db = postgres.NewPostgresDB(cfg)

	// Initialize repository
	taskRepo = postgres.NewTaskRepository(db)

	// Initialize service
	taskService = service.NewTaskService(taskRepo)

	// Initialize handler
	taskHandler = handler.NewTaskHandler(taskService)

}

func main() {
	initialize()
	// Setup router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(appMiddleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(appMiddleware.CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/tasks", taskHandler.Routes())
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server is running on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
