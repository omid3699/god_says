package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/omid3699/god_says/internal"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if server == nil {
		t.Fatal("Server is nil")
	}

	if server.god == nil {
		t.Fatal("Server god instance is nil")
	}
}

func TestServerHandleRoot(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleRoot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	if rr.Body.String() == "" {
		t.Error("Expected non-empty response body")
	}

	// Check content type
	expectedContentType := "text/plain; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
}

func TestServerHandleRootWithAmount(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req, err := http.NewRequest("GET", "/?amount=5", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleRoot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Check that we got some content (Happy.TXT contains phrases, not single words)
	response := rr.Body.String()
	if len(strings.TrimSpace(response)) == 0 {
		t.Error("Expected non-empty response")
	}
}

func TestServerHandleRootInvalidAmount(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	testCases := []string{"0", "-1", "1001", "abc"}
	for _, amount := range testCases {
		req, err := http.NewRequest("GET", fmt.Sprintf("/?amount=%s", amount), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.handleRoot)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d for amount %s, got %d", http.StatusBadRequest, amount, status)
		}
	}
}

func TestServerHandleJSON(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req, err := http.NewRequest("GET", "/json", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleJSON)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Parse JSON response
	var response GodResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response.GodSays == "" {
		t.Error("Expected non-empty god_says field")
	}
}

func TestServerHandleJSONWithAmount(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req, err := http.NewRequest("GET", "/json?amount=3", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleJSON)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var response GodResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that we got some content (Happy.TXT contains phrases, not single words)
	if len(strings.TrimSpace(response.GodSays)) == 0 {
		t.Error("Expected non-empty god_says field")
	}
}

func TestServerHandleHealth(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleHealth)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.WordsCount == 0 {
		t.Error("Expected non-zero words count")
	}

	if response.Uptime == "" {
		t.Error("Expected non-empty uptime")
	}
}

func TestServerMiddleware(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.handleRoot).Methods("GET", "OPTIONS")
	r.Use(securityMiddleware)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check security headers
	securityHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}

	for header, expectedValue := range securityHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("Expected %s header to be '%s', got '%s'", header, expectedValue, value)
		}
	}
}

func TestServerCORS(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.handleRoot).Methods("GET", "OPTIONS")
	r.Use(securityMiddleware)

	// Test OPTIONS request
	req, err := http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusOK, status)
	}

	// Check CORS headers
	corsHeaders := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
	}

	for _, header := range corsHeaders {
		if value := rr.Header().Get(header); value == "" {
			t.Errorf("Expected %s header to be present", header)
		}
	}
}

func TestParseAmount(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test valid amounts
	validCases := map[string]int{
		"":     internal.DefaultAmount,
		"1":    1,
		"10":   10,
		"100":  100,
		"1000": 1000,
	}

	for amountStr, expected := range validCases {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/?amount=%s", amountStr), nil)
		amount, err := server.parseAmount(req)
		if err != nil {
			t.Errorf("Unexpected error for amount '%s': %v", amountStr, err)
		}
		if amount != expected {
			t.Errorf("Expected amount %d for '%s', got %d", expected, amountStr, amount)
		}
	}

	// Test invalid amounts
	invalidCases := []string{"0", "-1", "1001", "abc", "10.5"}
	for _, amountStr := range invalidCases {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/?amount=%s", amountStr), nil)
		_, err := server.parseAmount(req)
		if err == nil {
			t.Errorf("Expected error for invalid amount '%s', got none", amountStr)
		}
	}
}

func TestServerIntegration(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.handleRoot).Methods("GET", "OPTIONS")
	r.HandleFunc("/json", server.handleJSON).Methods("GET", "OPTIONS")
	r.HandleFunc("/health", server.handleHealth).Methods("GET")
	r.Use(loggingMiddleware)
	r.Use(securityMiddleware)

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Test all endpoints
	endpoints := []string{"/", "/json", "/health"}
	for _, endpoint := range endpoints {
		resp, err := http.Get(testServer.URL + endpoint)
		if err != nil {
			t.Fatalf("Failed to make request to %s: %v", endpoint, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d for %s, got %d", http.StatusOK, endpoint, resp.StatusCode)
		}
	}
}

func TestServerTimeout(t *testing.T) {
	// Test that timeout middleware works
	handler := timeoutMiddleware(100 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow operation
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Should timeout and return 503
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d for timeout, got %d", http.StatusServiceUnavailable, status)
	}
}
