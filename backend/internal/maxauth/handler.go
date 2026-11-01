package maxauth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/exexeeeex/fiks-api/internal/config"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Handler struct {
	service     *Service
	testEnabled bool
}

func writeError(w http.ResponseWriter, status int, message string) {
	config.WriteJSON(w, status, errorResponse{Error: message})
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service:     service,
		testEnabled: os.Getenv("APP_ENV") == "development",
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/max", h.authenticate)
	if h.testEnabled {
		mux.HandleFunc("POST /api/v1/auth/max/test", h.testAuthenticate)
		mux.HandleFunc("GET /api/v1/auth/max/test-data", h.testData)
	}
}

func (h *Handler) testAuthenticate(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !h.testEnabled {
		http.NotFound(w, r)
		return
	}

	user := MaxUser{
		Id:        12345,
		FirstName: "Test",
		LastName:  "User",
		Username:  "test_user",
	}

	config.WriteJSON(w, http.StatusOK, AuthResponse{
		Success: true,
		User:    user,
	})
}

func (h *Handler) testData(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !h.testEnabled {
		http.NotFound(w, r)
		return
	}

	var user MaxUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if user.Id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	data, err := h.service.GenerateTestInitialData(user, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, map[string]string{
		"initData": data,
	})
}

func (h *Handler) authenticate(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.service.Authenticate(req.InitData)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, AuthResponse{
		Success: true,
		User:    user,
	})
}
