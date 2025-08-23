package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(APIResponse{Data: data}); err != nil {
		http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(APIResponse{Error: err.Error()}); err != nil {
		http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
	}
}
