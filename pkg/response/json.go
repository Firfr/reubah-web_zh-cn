package response

import (
	"encoding/json"
	_ "fmt"
	"net/http"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// Common content type header
const (
	jsonContentType = "application/json; charset=utf-8"
	maxAge = 31536000 // 1 year in seconds
)

// JSON sends a success response with data efficiently
func JSON(w http.ResponseWriter, code int, data interface{}) {
	// Set headers once
	h := w.Header()
	h.Set("Content-Type", jsonContentType)
	h.Set("X-Content-Type-Options", "nosniff")

	// Write status code
	w.WriteHeader(code)

	// Encode response
	response := Response{
		Success: code >= 200 && code < 300,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"success":false,"error":"Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}

// File sends a file download response
// func File(w http.ResponseWriter, filename string, contentType string) {
// 	h := w.Header()
// 	h.Set("Content-Type", contentType)
// 	h.Set("Content-Disposition", "attachment; filename="+filename)
// 	h.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
// 	h.Set("X-Content-Type-Options", "nosniff")
// }

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created sends a 201 Created response with optional location
func Created(w http.ResponseWriter, location string, data interface{}) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	JSON(w, http.StatusCreated, data)
} 