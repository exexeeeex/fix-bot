package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHandler struct{}

func (testHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCORS(t *testing.T) {
	server := New(testHandler{})

	tests := []struct {
		name             string
		method           string
		origin           string
		wantStatus       int
		wantAllowOrigin  string
		wantCredentials  string
		wantAllowMethods string
		wantAllowHeaders string
	}{
		{
			name:             "allowed preflight",
			method:           http.MethodOptions,
			origin:           allowedOrigin,
			wantStatus:       http.StatusNoContent,
			wantAllowOrigin:  allowedOrigin,
			wantCredentials:  "true",
			wantAllowMethods: "GET, POST, OPTIONS",
			wantAllowHeaders: "Authorization, Content-Type",
		},
		{
			name:            "allowed origin",
			method:          http.MethodGet,
			origin:          allowedOrigin,
			wantStatus:      http.StatusOK,
			wantAllowOrigin: allowedOrigin,
			wantCredentials: "true",
		},
		{
			name:       "other origin",
			method:     http.MethodGet,
			origin:     "http://localhost:3001",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/test", nil)
			request.Header.Set("Origin", tt.origin)
			response := httptest.NewRecorder()

			server.Handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowOrigin)
			}
			if got := response.Header().Get("Access-Control-Allow-Credentials"); got != tt.wantCredentials {
				t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, tt.wantCredentials)
			}
			if got := response.Header().Get("Access-Control-Allow-Methods"); got != tt.wantAllowMethods {
				t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, tt.wantAllowMethods)
			}
			if got := response.Header().Get("Access-Control-Allow-Headers"); got != tt.wantAllowHeaders {
				t.Errorf("Access-Control-Allow-Headers = %q, want %q", got, tt.wantAllowHeaders)
			}
		})
	}
}
