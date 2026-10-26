package student

import (
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/exexeeeex/fiks-api/internal/config"
	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/exexeeeex/fiks-api/internal/encrypt"
)

type generateQr struct {
	StudentId string `json:"student_id"`
}

type decryptRequest struct {
	Value string `json:"value"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/students/{id}", h.getStudentById)
	mux.HandleFunc("GET /api/v1/students/by-group/{groupId}", h.getStudentsByGroupId)
	mux.HandleFunc("POST /api/v1/students/validate-key", h.validateStudentKey)
	mux.HandleFunc("POST /api/v1/students", h.createStudent)
	mux.HandleFunc("POST /api/v1/students/decrypt", h.decryptStudent)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	config.WriteJSON(w, status, errorResponse{Error: message})
}

func (h *Handler) decryptStudent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req decryptRequest

	if err := config.DecodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Value == "" {
		writeError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}

	encrypted, err := base64.StdEncoding.DecodeString(req.Value)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base64 data")
		return
	}

	result, err := encrypt.Decrypt(encrypted, os.Getenv("QR_CODE_KEY"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "decryption failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}

func (h *Handler) validateStudentKey(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Key string `json:"key"`
	}

	if err := config.DecodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if request.Key == "" {
		writeError(w, http.StatusBadRequest, "value cannot be empty")
		return
	}

	result, err := h.service.ValidateStudentKey(r.Context(), request.Key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	model := struct {
		Success      bool      `json:"success"`
		ValidateTime time.Time `json:"validate_time"`
	}{
		Success:      result,
		ValidateTime: time.Now(),
	}

	config.WriteJSON(w, http.StatusOK, model)
}

func (h *Handler) getStudentsByGroupId(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")

	result, err := h.service.GetByGroupId(r.Context(), groupId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getStudentById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) createStudent(w http.ResponseWriter, r *http.Request) {
	var req createStudentRequest

	if err := config.DecodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	student := domain.Student{
		Name:     req.Name,
		GroupId:  req.GroupId,
		IsActive: req.IsActive,
	}

	result, err := h.service.Create(r.Context(), student)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusCreated, result)
}
