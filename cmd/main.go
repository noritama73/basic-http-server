package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/finatext/academia-basic-http-server/internal/interface/handler"
	"github.com/finatext/academia-basic-http-server/internal/interface/repository"
	"github.com/finatext/academia-basic-http-server/internal/usecase"
	"github.com/finatext/academia-basic-http-server/internal/util"
)

const (
	defaultPort      = "8080"
	defaultJWTSecret = "academia-secret-key"
	jwtExpiry        = 24 * time.Hour // Token valid for 24 hours
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Get JWT secret from environment variable or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
	}

	// Initialize components
	jwtManager := util.NewJWTManager(jwtSecret, jwtExpiry)
	userRepo := repository.NewInMemoryUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepo, jwtManager)
	userHandler := handler.NewUserHandler(userUseCase, jwtManager)

	// Set up HTTP server
	mux := http.NewServeMux()
	
	// Register routes
	userHandler.RegisterRoutes(mux)
	
	// Add middleware for logging
	handler := logMiddleware(mux)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// logMiddleware logs all HTTP requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Response: %s %s - %s", r.Method, r.URL.Path, time.Since(start))
	})
}
