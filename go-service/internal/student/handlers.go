package student

import (
	"errors"
	"fmt"
	"goservice/internal/client"
	"goservice/internal/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var (
	ErrAuthTokens = errors.New("missing or invalid required tokens")
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", h.GetStudent)
	r.Get("/{id}/report", h.GenerateReport)
	return r
}

func checkRequiredCookie(r *http.Request) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	accessToken, err := r.Cookie(client.AccesTokenName)
	if err != nil || accessToken.Value == "" {
		return nil, ErrAuthTokens
	}
	refreshToken, err := r.Cookie(client.RefreshTokenName)
	if err != nil || refreshToken.Value == "" {
		return nil, ErrAuthTokens
	}
	csrfToken, err := r.Cookie(client.CSFRTokenName)
	if err != nil || csrfToken.Value == "" {
		return nil, ErrAuthTokens
	}
	cookies = append(cookies, accessToken, refreshToken, csrfToken)
	// Here you would typically validate the session ID against your session store
	return cookies, nil
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	cookies, err := checkRequiredCookie(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	student, err := h.service.GetStudent(r.Context(), id, cookies)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, student)
}

func (h *Handler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	cookies, err := checkRequiredCookie(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	pdf, err := h.service.GenerateReport(r.Context(), id, cookies)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=student_%d_report.pdf", id))
	w.WriteHeader(http.StatusOK)
	if err := pdf.Output(w); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
}
