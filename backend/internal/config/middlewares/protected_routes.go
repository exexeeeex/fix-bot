package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userIDContextKey struct{}
type maxIDContextKey struct{}
type organizationIDContextKey struct{}

func AuthenticatedUserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDContextKey{}).(string)
	return userID
}

func AuthenticatedMaxID(ctx context.Context) string {
	maxID, _ := ctx.Value(maxIDContextKey{}).(string)
	return maxID
}

func AuthenticatedOrganizationID(ctx context.Context) string {
	organizationID, _ := ctx.Value(organizationIDContextKey{}).(string)
	return organizationID
}

type JwtService interface {
	ValidateToken(tokenString string, tokenType string) (*jwt.Token, error)
}

type ProtectedRoutes struct {
	jwtService JwtService
}

func NewProtectedRoutes(jwtService JwtService) *ProtectedRoutes {
	return &ProtectedRoutes{
		jwtService: jwtService,
	}
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenString, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header")
	}
	tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
	if tokenString == "" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return tokenString, nil
}

func (s *ProtectedRoutes) RoleProtectedHandler(
	w http.ResponseWriter,
	r *http.Request,
	next http.HandlerFunc,
	roles ...string,
) http.HandlerFunc {
	tokenString, err := GetTokenFromHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "missing or invalid authorization header")
		return nil
	}

	token, err := s.jwtService.ValidateToken(tokenString, "access")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return nil
	}

	userRole, ok := claims["role"].(string)
	if !ok || !s.checkRoles(w, userRole, roles) {
		return nil
	}

	applicationStatus, ok := claims["applicationStatus"].(string)
	if !ok || applicationStatus != "approved" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "forbidden: application not approved")
		return nil
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return nil
	}

	r = r.WithContext(context.WithValue(r.Context(), userIDContextKey{}, userID))
	if maxID, ok := claims["maxId"].(string); ok {
		r = r.WithContext(context.WithValue(r.Context(), maxIDContextKey{}, maxID))
	}
	if organizationID, ok := claims["organizationId"].(string); ok {
		r = r.WithContext(context.WithValue(r.Context(), organizationIDContextKey{}, organizationID))
	}
	next(w, r)
	return next
}

func (s *ProtectedRoutes) ApprovedApplicationProtectedHandler(
	w http.ResponseWriter,
	r *http.Request,
	next http.HandlerFunc,
) {
	tokenString, err := GetTokenFromHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "missing or invalid authorization header")
		return
	}

	token, err := s.jwtService.ValidateToken(tokenString, "access")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return
	}

	applicationStatus, ok := claims["applicationStatus"].(string)
	if !ok || applicationStatus != "approved" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "forbidden: application not approved")
		return
	}

	next(w, r)
}

func (s *ProtectedRoutes) ProtectedHandler(
	w http.ResponseWriter,
	r *http.Request,
	next http.HandlerFunc,
) {
	tokenString, err := GetTokenFromHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "missing or invalid authorization header")
		return
	}

	token, err := s.jwtService.ValidateToken(tokenString, "access")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid authorization token")
		return
	}
	if userID, ok := claims["id"].(string); ok && userID != "" {
		r = r.WithContext(context.WithValue(r.Context(), userIDContextKey{}, userID))
	}
	if maxID, ok := claims["maxId"].(string); ok && maxID != "" {
		r = r.WithContext(context.WithValue(r.Context(), maxIDContextKey{}, maxID))
	}
	if organizationID, ok := claims["organizationId"].(string); ok {
		r = r.WithContext(context.WithValue(r.Context(), organizationIDContextKey{}, organizationID))
	}
	next(w, r)
}

func (s *ProtectedRoutes) checkRoles(
	w http.ResponseWriter,
	userRole string,
	roles []string,
) bool {
	for _, role := range roles {
		if userRole == role {
			return true
		}
	}

	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "forbidden: insufficient role")
	return false
}
