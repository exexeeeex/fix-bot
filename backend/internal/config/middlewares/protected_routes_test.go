package config

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type testJwtService struct {
	token *jwt.Token
	err   error
}

func (s testJwtService) ValidateToken(string, string) (*jwt.Token, error) {
	return s.token, s.err
}

func TestApprovedApplicationProtectedHandler(t *testing.T) {
	tests := []struct {
		name           string
		token          *jwt.Token
		tokenErr       error
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name: "approved application",
			token: &jwt.Token{
				Claims: jwt.MapClaims{"applicationStatus": "approved"},
				Valid:  true,
			},
			wantStatus:     http.StatusNoContent,
			wantNextCalled: true,
		},
		{
			name: "pending application",
			token: &jwt.Token{
				Claims: jwt.MapClaims{"applicationStatus": "pending"},
				Valid:  true,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "missing application status",
			token: &jwt.Token{
				Claims: jwt.MapClaims{},
				Valid:  true,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "invalid token",
			tokenErr:   errors.New("invalid token"),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routes := NewProtectedRoutes(testJwtService{token: tt.token, err: tt.tokenErr})
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Authorization", "Bearer test-token")
			response := httptest.NewRecorder()
			nextCalled := false

			routes.ApprovedApplicationProtectedHandler(response, request, func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				nextCalled = true
				w.WriteHeader(http.StatusNoContent)
			})

			if response.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %t, want %t", nextCalled, tt.wantNextCalled)
			}
		})
	}
}

func TestRoleProtectedHandlerAllowsAnyConfiguredRole(t *testing.T) {
	tests := []struct {
		name           string
		userRole       string
		allowedRoles   []string
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name:           "first allowed role",
			userRole:       "teacher",
			allowedRoles:   []string{"teacher", "admin"},
			wantStatus:     http.StatusNoContent,
			wantNextCalled: true,
		},
		{
			name:           "another allowed role",
			userRole:       "admin",
			allowedRoles:   []string{"teacher", "admin"},
			wantStatus:     http.StatusNoContent,
			wantNextCalled: true,
		},
		{
			name:         "role is not allowed",
			userRole:     "student",
			allowedRoles: []string{"teacher", "admin"},
			wantStatus:   http.StatusForbidden,
		},
		{
			name:       "no roles configured",
			userRole:   "teacher",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routes := NewProtectedRoutes(testJwtService{
				token: &jwt.Token{
					Claims: jwt.MapClaims{
						"role":              tt.userRole,
						"applicationStatus": "approved",
						"id":                "user-123",
						"organizationId":    "organization-123",
					},
					Valid: true,
				},
			})
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Authorization", "Bearer token")
			response := httptest.NewRecorder()
			nextCalled := false

			routes.RoleProtectedHandler(response, request, func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				nextCalled = true
				if got := AuthenticatedOrganizationID(r.Context()); got != "organization-123" {
					t.Errorf("authenticated organization id = %q, want %q", got, "organization-123")
				}
				if got := AuthenticatedMaxID(r.Context()); got != "" {
					t.Errorf("authenticated max id = %q, want empty", got)
				}
				w.WriteHeader(http.StatusNoContent)
			}, tt.allowedRoles...)

			if response.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %t, want %t", nextCalled, tt.wantNextCalled)
			}
		})
	}
}

func TestProtectedHandlerAddsUserScopeToRequestContext(t *testing.T) {
	routes := NewProtectedRoutes(testJwtService{
		token: &jwt.Token{
			Claims: jwt.MapClaims{
				"id":             "user-123",
				"maxId":          "max-123",
				"organizationId": "organization-123",
			},
			Valid: true,
		},
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	routes.ProtectedHandler(response, request, func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if got := AuthenticatedUserID(r.Context()); got != "user-123" {
			t.Errorf("authenticated user id = %q, want %q", got, "user-123")
		}
		if got := AuthenticatedMaxID(r.Context()); got != "max-123" {
			t.Errorf("authenticated max id = %q, want %q", got, "max-123")
		}
		if got := AuthenticatedOrganizationID(r.Context()); got != "organization-123" {
			t.Errorf("authenticated organization id = %q, want %q", got, "organization-123")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}
