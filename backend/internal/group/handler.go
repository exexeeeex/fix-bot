package group

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
	mux.HandleFunc("GET /api/v1/groups/by-org/{organizationId}", h.listGroups)
	mux.HandleFunc("GET /api/v1/groups/{id}", h.getGroupById)
	mux.HandleFunc("GET /api/v1/groups/by-name/{name}/{organizationId}", h.getGroupByName)
}

func (h *Handler) listGroups(w http.ResponseWriter, r *http.Request) {
	organizationId := r.PathValue("organizationId")
	result, err := h.service.List(r.Context(), organizationId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getGroupById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getGroupByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	organizationId := r.PathValue("organizationId")

	result, err := h.service.GetByName(r.Context(), name, organizationId)
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
