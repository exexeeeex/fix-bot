package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTSecretsAreReadAfterEnvironmentSetup(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET_KEY", "development-access-secret")
	t.Setenv("JWT_REFRESH_SECRET_KEY", "development-refresh-secret")
	t.Setenv("JWT_ACCESS_EXP", "1h")
	t.Setenv("JWT_REFRESH_EXP", "24h")

	service := NewService()
	accessToken, err := service.WriteAccessToken(AccessRequest{Id: "user-1"})
	if err != nil {
		t.Fatalf("WriteAccessToken() error = %v", err)
	}

	token, err := service.ValidateToken(accessToken, "access")
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		t.Fatalf("claims type = %T, want jwt.MapClaims", token.Claims)
	}
}

func TestJWTRejectsInvalidExpiration(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET_KEY", "development-access-secret")
	t.Setenv("JWT_ACCESS_EXP", "not-a-duration")

	if _, err := NewService().WriteAccessToken(AccessRequest{}); err == nil {
		t.Fatal("WriteAccessToken() error = nil, want invalid expiration error")
	}
}
