package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/omid3699/god_says/internal"
)

const (
	// RequestTimeout is the maximum time for a request
	RequestTimeout = 30 * time.Second
	// ShutdownTimeout is the maximum time to wait for graceful shutdown
	ShutdownTimeout = 30 * time.Second
)

// GodResponse represents the JSON response structure
type GodResponse struct {
	GodSays string `json:"god_says"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status     string `json:"status"`
	WordsCount int    `json:"words_count"`
	Uptime     string `json:"uptime"`
}

// Server holds the server state and dependencies
type Server struct {
	god       *internal.God
	startTime time.Time
}

// NewServer creates a new server instance
func NewServer() (*Server, error) {
	god, err := internal.NewGod(internal.DefaultAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to create god instance: %w", err)
	}

	return &Server{
		god:       god,
		startTime: time.Now(),
	}, nil
}

// writeErrorResponse writes a JSON error response
func (s *Server) writeErrorResponse(w http.ResponseWriter, statusCode int, errorMsg, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   errorMsg,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// parseAmount parses and validates the amount parameter from request
func (s *Server) parseAmount(r *http.Request) (int, error) {
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		return internal.DefaultAmount, nil
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return 0, fmt.Errorf("invalid amount parameter: must be a number")
	}

	if amount < internal.MinAmount || amount > internal.MaxAmount {
		return 0, fmt.Errorf("amount must be between %d and %d", internal.MinAmount, internal.MaxAmount)
	}
	return amount, nil
}

// RunServer starts the HTTP server
func RunServer(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	server, err := NewServer()
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.handleRoot).Methods("GET", "OPTIONS")
	r.HandleFunc("/json", server.handleJSON).Methods("GET", "OPTIONS")
	r.HandleFunc("/health", server.handleHealth).Methods("GET", "OPTIONS")

	// Add middlewares
	r.Use(loggingMiddleware)
	r.Use(securityMiddleware)
	r.Use(timeoutMiddleware(RequestTimeout))

	// Create HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start Server in a goroutin
	go func() {
		log.Printf("Endpoints:")
		log.Printf("  GET /        - Plain text response")
		log.Printf("  GET /json    - JSON response")
		log.Printf("  GET /health  - Health check")

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
	return nil
}
