package application

import (
	"net/http"

	"github.com/exexeeeex/fiks-api/internal/config"
	"github.com/exexeeeex/fiks-api/internal/domain"
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
	mux.HandleFunc("GET /api/v1/application/{groupId}", h.getByGroupId)
	mux.HandleFunc("POST /api/v1/application/send", h.create)
}

type errorResponse struct {
	Error string `json:"errror"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	config.WriteJSON(w, status, errorResponse{Error: message})
}

func (h *Handler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req createApplicationRequest

	if err := config.DecodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.GroupId == "" ||
		req.Name == "" ||
		req.OrganizationId == "" ||
		req.UserId == "" {
		writeError(w, http.StatusBadRequest, "values cannot be null or empty")
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) getByGroupId(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("groupId")

	if id == "" {
		writeError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}

	result, err := h.service.GetByGroupId(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
