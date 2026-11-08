package visit

import (
	"context"
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
	configmiddleware "github.com/exexeeeex/fiks-api/internal/config/middlewares"
	"github.com/exexeeeex/fiks-api/internal/domain"
)

type GroupService interface {
	GetById(ctx context.Context, id string) (domain.Group, error)
}

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
	groupService    GroupService
}

func NewHandler(
	service *Service,
	protectedRoutes ProtectedRoutes,
	groupService GroupService,
) *Handler {
	return &Handler{
		service:         service,
		protectedRoutes: protectedRoutes,
		groupService:    groupService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/visits/by-group/{groupId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getVisitsByGroupId, "teacher", "admin")
	}))
	mux.Handle("GET /api/v1/visits/by-organization/{organizationId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getVisitsByOrganizationId, "teacher", "admin")
	}))
}

func (h *Handler) getVisitsByOrganizationId(
	w http.ResponseWriter,
	r *http.Request,
) {
	organizationId := r.PathValue("organizationId")

	if organizationId == "" {
		config.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}
	authenticatedOrganizationID := configmiddleware.AuthenticatedOrganizationID(r.Context())
	if authenticatedOrganizationID == "" {
		config.WriteError(w, http.StatusUnauthorized, "organization scope is missing")
		return
	}
	if organizationId != authenticatedOrganizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	result, err := h.service.ListByOrganizationId(r.Context(), organizationId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getVisitsByGroupId(
	w http.ResponseWriter,
	r *http.Request,
) {
	groupId := r.PathValue("groupId")

	if groupId == "" {
		config.WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}
	organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context())
	if organizationID == "" {
		config.WriteError(w, http.StatusUnauthorized, "organization scope is missing")
		return
	}
	group, err := h.groupService.GetById(r.Context(), groupId)
	if err != nil {
		config.WriteError(w, http.StatusNotFound, "group not found")
		return
	}
	if group.OrganizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	result, err := h.service.ListByGroupId(r.Context(), groupId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
