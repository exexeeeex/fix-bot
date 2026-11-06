package organization

import (
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
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
	mux.HandleFunc("GET /api/v1/organizations", h.list)
	mux.HandleFunc("POST /api/v1/organizations/by-name", h.getByName)
}

func (h *Handler) getByName(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req requestByName
	if err := config.DecodeJSON(w, r, &req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Name == "" {
		config.WriteError(w, http.StatusBadRequest, "invalid name")
		return
	}

	result, err := h.service.GetByName(r.Context(), req.Name)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) list(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.List(r.Context())
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
