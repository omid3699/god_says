// Package server provides god says HTTP server functionality.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/omid3699/god_says/internal"
)

// handleRoot handles the root endpoint returning plain text
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	amount, err := s.parseAmount(r)
	if err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	var message string
	if amount == internal.DefaultAmount {
		message = s.god.Speak()
	} else {
		message, err = s.god.SpeakWithAmount(amount)
		if err != nil {
			s.writeErrorResponse(w, http.StatusInternalServerError, "generation_error", err.Error())
			return
		}
	}
	if message == "" {
		s.writeErrorResponse(w, http.StatusInternalServerError, "empty_message", "Failed to generate message")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, message)
}

// handleJSON handles the JSON endpoint
func (s *Server) handleJSON(w http.ResponseWriter, r *http.Request) {
	amount, err := s.parseAmount(r)
	if err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	var message string
	if amount == internal.DefaultAmount {
		message = s.god.Speak()
	} else {
		message, err = s.god.SpeakWithAmount(amount)
		if err != nil {
			s.writeErrorResponse(w, http.StatusBadRequest, "generation_error", err.Error())
			return
		}
	}

	if message == "" {
		s.writeErrorResponse(w, http.StatusInternalServerError, "empty_message", "Failed to generate message")
		return
	}

	response := GodResponse{GodSays: message}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		s.writeErrorResponse(w, http.StatusInternalServerError, "encoding_error", "Failed to encode response")
	}
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.startTime).Round(time.Second).String()
	response := HealthResponse{
		Status:     "ok",
		WordsCount: s.god.GetWordsCount(),
		Uptime:     uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode health response: %v", err)
	}
}
