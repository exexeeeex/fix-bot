package qr

import (
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
	middlewares "github.com/exexeeeex/fiks-api/internal/config/middlewares"
)

type ProtectedRoutes interface {
	ProtectedHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	RoleProtectedHandler(
		w http.ResponseWriter,
		r *http.Request,
		next http.HandlerFunc,
		roles ...string,
	) http.HandlerFunc
}

type Handler struct {
	service         *Service
	protectedRoutes ProtectedRoutes
}

func NewHandler(
	service *Service,
	protectedRoutes ProtectedRoutes,
) *Handler {
	return &Handler{
		service:         service,
		protectedRoutes: protectedRoutes,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/qr-code/generate", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.generateQrCode, "student")
	}))
	mux.HandleFunc("POST /api/v1/qr-code/read", h.readQrCode)
}

func (h *Handler) readQrCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Value string `json:"value"`
	}

	if err := config.DecodeJSON(w, r, &req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON: %w"+err.Error())
		return
	}

	if req.Value == "" {
		config.WriteError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}

	result, err := h.service.ReadQrCode(r.Context(), req.Value)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, "read error: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) generateQrCode(w http.ResponseWriter, r *http.Request) {
	userId := middlewares.AuthenticatedUserID(r.Context())
	if userId == "" {
		config.WriteError(w, http.StatusUnauthorized, "authenticated user id is missing")
		return
	}

	result, err := h.service.GenerateQrCodeForUser(r.Context(), userId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", `inline; filename="student-qr.png"`)
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(result); err != nil {
		return
	}
}
