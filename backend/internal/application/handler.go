package application

import (
	"errors"
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
	configmiddleware "github.com/exexeeeex/fiks-api/internal/config/middlewares"
	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
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

func NewHandler(service *Service, protectedRoutes ProtectedRoutes) *Handler {
	return &Handler{
		service:         service,
		protectedRoutes: protectedRoutes,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/application/by-group/{groupId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getByGroupId, "teacher")
	}))
	mux.HandleFunc("GET /api/v1/application/{id}", h.getById)
	mux.Handle("GET /api/v1/application/by-max/{maxId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.ProtectedHandler(w, r, h.getByMaxId)
	}))
	mux.Handle("GET /api/v1/application/by-organization/{organizationId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.getByOrganizationId, "teacher", "admin")
	}))
	mux.Handle("POST /api/v1/application/send", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.ProtectedHandler(w, r, h.create)

	}))
	mux.Handle("PUT /api/v1/application/update-status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.protectedRoutes.RoleProtectedHandler(w, r, h.update, "teacher", "admin")
	}))
}

func (h *Handler) getByMaxId(
	w http.ResponseWriter,
	r *http.Request,
) {
	maxId := r.PathValue("maxId")

	if maxId == "" {
		config.WriteError(w, http.StatusBadRequest, "max id cannot be empty")
		return
	}
	if maxId != configmiddleware.AuthenticatedMaxID(r.Context()) {
		config.WriteError(w, http.StatusForbidden, "forbidden: cannot view another user's application")
		return
	}

	result, err := h.service.GetByMaxId(r.Context(), maxId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.WriteError(w, http.StatusNotFound, "application not found")
			return
		}
		config.WriteError(w, http.StatusInternalServerError, "failed to get application by max id: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getById(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	if id == "" {
		config.WriteError(w, http.StatusBadRequest, "id cannot be empty")
		return
	}

	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, "failed to get application by id: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req createApplicationRequest

	if err := config.DecodeJSON(w, r, &req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Name == "" ||
		req.OrganizationId == "" ||
		req.UserId == "" {
		config.WriteError(w, http.StatusBadRequest, "values cannot be null or empty")
		return
	}
	if req.UserId != configmiddleware.AuthenticatedUserID(r.Context()) {
		config.WriteError(w, http.StatusForbidden, "forbidden: cannot create an application for another user")
		return
	}
	if organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context()); organizationID != "" &&
		req.OrganizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	switch req.Type {
	case domain.ApplicationTypeStudent:
		if req.GroupId == nil || *req.GroupId == "" {
			config.WriteError(w, http.StatusBadRequest, "group_id is required for student applications")
			return
		}
	case domain.ApplicationTypeTeacher:
		req.GroupId = nil
	default:
		config.WriteError(w, http.StatusBadRequest, "unsupported application type")
		return
	}

	application := domain.Application{
		Name:           req.Name,
		Status:         domain.ApplicationStatusPending,
		Type:           req.Type,
		UserId:         req.UserId,
		OrganizationId: req.OrganizationId,
		GroupId:        req.GroupId,
	}

	result, err := h.service.Create(r.Context(), application)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) update(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req updateApplicationRequest
	if err := config.DecodeJSON(w, r, &req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Id == "" {
		config.WriteError(w, http.StatusBadRequest, "id cannot be empty")
		return
	}

	current, err := h.service.GetById(r.Context(), req.Id)
	if err != nil {
		config.WriteError(w, http.StatusNotFound, "application not found")
		return
	}
	organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context())
	if organizationID == "" {
		config.WriteError(w, http.StatusUnauthorized, "organization scope is missing")
		return
	}
	if current.OrganizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	if req.Status != "" {
		current.Status = req.Status
	}

	result, err := h.service.UpdateApplicationStatus(r.Context(), current, req.Status)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) getByOrganizationId(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("organizationId")

	if id == "" {
		config.WriteError(w, http.StatusBadRequest, "organization id cannot be empty")
		return
	}
	organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context())
	if organizationID == "" {
		config.WriteError(w, http.StatusUnauthorized, "organization scope is missing")
		return
	}
	if id != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	result, err := h.service.GetByOrganizationId(r.Context(), id)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getByGroupId(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("groupId")

	if id == "" {
		config.WriteError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}
	organizationID := configmiddleware.AuthenticatedOrganizationID(r.Context())
	if organizationID == "" {
		config.WriteError(w, http.StatusUnauthorized, "organization scope is missing")
		return
	}
	group, err := h.service.groupService.GetById(r.Context(), id)
	if err != nil {
		config.WriteError(w, http.StatusNotFound, "group not found")
		return
	}
	if group.OrganizationId != organizationID {
		config.WriteError(w, http.StatusForbidden, "forbidden: organization access denied")
		return
	}

	result, err := h.service.GetByGroupId(r.Context(), id)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
