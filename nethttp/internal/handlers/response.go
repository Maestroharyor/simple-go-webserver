package handlers

import (
	"encoding/json"
	"net/http"
)

type response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitzero"`
}

func WriteSuccess[T any](w http.ResponseWriter, status int, message string, data T) {
	WriteJSON(w, status, response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, response[*struct{}]{
		Success: false,
		Message: message,
	})
}

func WriteJSON[T any](w http.ResponseWriter, status int, body response[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
