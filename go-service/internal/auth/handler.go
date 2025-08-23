package auth

import (
	"encoding/json"
	"goservice/internal/client"
	"goservice/internal/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	client client.IBackend
}

func NewHandler(s client.IBackend) *Handler {
	return &Handler{client: s}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	return r
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Call service to authenticate and get cookies
	cookies, err := h.client.Login(r.Context(), creds.Username, creds.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	// Set cookies in response
	for _, c := range cookies {
		http.SetCookie(w, c)
	}

	// Send JSON response back
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Login successful",
	})
}
