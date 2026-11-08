package maxauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/exexeeeex/fiks-api/internal/config"
)

type Handler struct {
	service     *Service
	testEnabled bool
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
	mux.HandleFunc("POST /api/v1/auth/max/refresh", h.refresh)
	mux.HandleFunc("POST /api/v1/auth/demo", h.demoAuthenticate)
	if h.testEnabled {
		mux.HandleFunc("POST /api/v1/auth/max/test", h.testAuthenticate)
		mux.HandleFunc("GET /api/v1/auth/max/test-data", h.testData)
	}
}

func (h *Handler) demoAuthenticate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}
	if err := config.DecodeJSON(w, r, &req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	response, err := h.service.AuthenticateDemo(r.Context(), req.Role)
	if errors.Is(err, ErrInvalidDemoRole) {
		config.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, "failed to authenticate demo user: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) || errors.Is(err, ErrRefreshTokenRequired) {
			config.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		config.WriteError(w, http.StatusInternalServerError, "failed to refresh tokens: "+err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, tokens)
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
		User: AuthenticatedUser{
			Id:           "test-user-id",
			MaxId:        "12345",
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Username:     user.Username,
			LanguageCode: user.LanguageCode,
			PhotoURL:     user.PhotoURL,
		},
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
		config.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if user.Id <= 0 {
		config.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	data, err := h.service.GenerateTestInitialData(user, time.Now())
	if err != nil {
		config.WriteError(w, http.StatusInternalServerError, err.Error())
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
		config.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	response, err := h.service.Authenticate(r.Context(), req.InitData)
	if err != nil {
		config.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	config.WriteJSON(w, http.StatusOK, response)
}
