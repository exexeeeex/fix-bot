package qr

import (
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/qr-code/generate", h.generateQrCode)
	mux.HandleFunc("POST /api/v1/qr-code/read", h.readQrCode)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	config.WriteJSON(w, status, errorResponse{Error: message})
}

func (h *Handler) readQrCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Value string `json:"value"`
	}

	if err := config.DecodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: %w"+err.Error())
		return
	}

	if req.Value == "" {
		writeError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}

	result, err := h.service.ReadQrCode(r.Context(), req.Value)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read error: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) generateQrCode(w http.ResponseWriter, r *http.Request) {
	var req generateQrRequest

	if err := config.DecodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON:"+err.Error())
		return
	}

	if req.StudentId == "" {
		writeError(w, http.StatusBadRequest, "student_id cannot be empty")
		return
	}

	result, err := h.service.GenerateQrCode(r.Context(), req.StudentId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", `inline; filename="student-qr.png"`)
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(result); err != nil {
		return
	}
}
