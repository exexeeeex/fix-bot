package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	configmiddleware "github.com/exexeeeex/fiks-api/internal/config/middlewares"
	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type testProtectedRoutes struct{}

func (testProtectedRoutes) ProtectedHandler(http.ResponseWriter, *http.Request, http.HandlerFunc) {
}

func (testProtectedRoutes) RoleProtectedHandler(
	_ http.ResponseWriter,
	_ *http.Request,
	next http.HandlerFunc,
	_ ...string,
) http.HandlerFunc {
	return next
}

func TestRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewHandler(nil, testProtectedRoutes{})
	handler.RegisterRoutes(mux)

	tests := []struct {
		path        string
		wantPattern string
	}{
		{
			path:        "/api/v1/application/by-group/group-123",
			wantPattern: "GET /api/v1/application/by-group/{groupId}",
		},
		{
			path:        "/api/v1/application/app-123",
			wantPattern: "GET /api/v1/application/{id}",
		},
	}

	for _, tt := range tests {
		request, err := http.NewRequest(http.MethodGet, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		_, pattern := mux.Handler(request)
		if pattern != tt.wantPattern {
			t.Errorf("pattern for %q = %q, want %q", tt.path, pattern, tt.wantPattern)
		}
	}

}

type notFoundApplicationRepository struct{}

func (notFoundApplicationRepository) GetById(context.Context, string) (domain.Application, error) {
	return domain.Application{}, nil
}

func (notFoundApplicationRepository) GetByGroupId(context.Context, string) ([]*domain.Application, error) {
	return nil, nil
}

func (notFoundApplicationRepository) GetByMaxId(context.Context, string) (domain.Application, error) {
	return domain.Application{}, gorm.ErrRecordNotFound
}

func (notFoundApplicationRepository) GetByUserId(context.Context, string) (domain.Application, error) {
	return domain.Application{}, nil
}

func (notFoundApplicationRepository) List(context.Context, string) ([]domain.Application, error) {
	return nil, nil
}

func (notFoundApplicationRepository) Create(context.Context, *domain.Application) error {
	return nil
}

func (notFoundApplicationRepository) Update(context.Context, *domain.Application) error {
	return nil
}

func (notFoundApplicationRepository) GetByOrganizationId(context.Context, string) ([]*domain.Application, error) {
	return []*domain.Application{}, nil
}

type claimsJwtService struct {
	claims jwt.MapClaims
}

func (s claimsJwtService) ValidateToken(string, string) (*jwt.Token, error) {
	return &jwt.Token{Claims: s.claims, Valid: true}, nil
}

func TestGetByOrganizationIdRestrictsAccessToTokenOrganization(t *testing.T) {
	service := NewService(notFoundApplicationRepository{}, nil, nil, nil, nil, nil, nil, nil)
	protectedRoutes := configmiddleware.NewProtectedRoutes(claimsJwtService{
		claims: jwt.MapClaims{
			"role":              "teacher",
			"applicationStatus": "approved",
			"id":                "teacher-id",
			"organizationId":    "organization-123",
		},
	})
	handler := NewHandler(service, protectedRoutes)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{name: "own organization", path: "/api/v1/application/by-organization/organization-123", wantStatus: http.StatusOK},
		{name: "different organization", path: "/api/v1/application/by-organization/other-organization", wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			request.Header.Set("Authorization", "Bearer access-token")
			response := httptest.NewRecorder()

			mux.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetByMaxIdReturnsNotFoundWhenApplicationDoesNotExist(t *testing.T) {
	service := NewService(notFoundApplicationRepository{}, nil, nil, nil, nil, nil, nil, nil)
	protectedRoutes := configmiddleware.NewProtectedRoutes(claimsJwtService{
		claims: jwt.MapClaims{"maxId": "756186104"},
	})
	handler := NewHandler(service, protectedRoutes)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/application/by-max/756186104", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
