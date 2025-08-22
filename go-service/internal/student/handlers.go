package student

import (
	"fmt"
	"goservice/internal/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	student, err := h.service.GetStudent(r.Context(), id)
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

	pdf, err := h.service.GenerateReport(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=student_%d_report.pdf", id))
	w.WriteHeader(http.StatusOK)
	pdf.Output(w)
}
