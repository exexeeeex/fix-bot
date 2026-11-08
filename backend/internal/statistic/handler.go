package statistic

import (
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
	configmiddleware "github.com/exexeeeex/fiks-api/internal/config/middlewares"
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
	service *Service
	protectedRoutes ProtectedRoutes
}

func NewHandler(
	service *Service,
	protectedRoutes ProtectedRoutes,
) *Handler {
	return &Handler{
		service: service,
		protectedRoutes: protectedRoutes,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/statistic/organization/{organizationId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getStatistics, "teacher", "admin")
	}))
	mux.Handle("GET /api/v1/statistic/organization-abcenses/{organizationId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getAbcenses, "teacher", "admin")
	}))
}

func (h *Handler) getStatistics(
	w http.ResponseWriter,
	r *http.Request,
) {
	organizationId := r.PathValue("organizationId")
	if organizationId == "" {
		config.WriteError(w, http.StatusBadRequest, "organization id cannot be empty")
		return
	}

	result, err := h.service.GetOrganizationStatistics(r.Context(), organizationId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context()); organizationID != "" &&
		organizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getAbcenses(
	w http.ResponseWriter,
	r *http.Request,
) {
	organizationId := r.PathValue("organizationId")
	if organizationId == "" {
		config.WriteError(w, http.StatusBadRequest, "organization id cannot be empty")
		return
	}
	
	result, err := h.service.GetOrganizationAbsences(r.Context(), organizationId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context()); organizationID != "" &&
		organizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}