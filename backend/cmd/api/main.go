package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum-diapason/infrastructure/database"
	"forum-diapason/infrastructure/http/handlers"
	"forum-diapason/internal/domain/repositories"
	"forum-diapason/pkg/config"
	"forum-diapason/pkg/logger"
	"forum-diapason/usecases"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()
	l := logger.New()

	l.Info("Starting Forum Diapason API server...")

	// Initialize database
	db, err := database.NewSQLiteConnection(cfg)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to initialize database: %v", err))
		os.Exit(1)
	}
	defer db.Close()

	l.Info("Database connected and migrations run successfully")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo)
	postUseCase := usecases.NewPostUseCase(postRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userUseCase, cfg)
	postHandler := handlers.NewPostHandler(postUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Initialize router
	router := mux.NewRouter()

	// Redirect root to frontend dev server
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
	})

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected routes (would need JWT middleware)
	api.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts/{id}", postHandler.GetPost).Methods("GET")
	api.HandleFunc("/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	api.HandleFunc("/posts/{id}", postHandler.DeletePost).Methods("DELETE")
	api.HandleFunc("/posts/{id}/like", postHandler.LikePost).Methods("POST")
	api.HandleFunc("/posts/{id}/like", postHandler.UnlikePost).Methods("DELETE")

	api.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	api.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")
	api.HandleFunc("/users/change-password", userHandler.ChangePassword).Methods("PUT")
	api.HandleFunc("/users/{userId}/posts", userHandler.GetUserPosts).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Server configuration
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      c.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		l.Info(fmt.Sprintf("Server starting on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error(fmt.Sprintf("Server failed to start: %v", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
		os.Exit(1)
	}

	l.Info("Server exited")
}
