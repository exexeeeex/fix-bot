package user

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
	mux.HandleFunc("GET /api/v1/users/{id}", h.getUserById)
}

func (h *Handler) getUserById(
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
		config.WriteError(w, http.StatusInternalServerError, "failed to get user by id: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}
