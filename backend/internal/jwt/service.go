package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) WriteTokensPair(request AccessRequest) (TokensPair, error) {
	accessToken, err := s.WriteAccessToken(request)
	if err != nil {
		return TokensPair{}, err
	}

	refreshToken, err := s.WriteRefreshToken(RefreshRequest{
		Id:    request.Id,
		MaxId: request.MaxId,
	})
	if err != nil {
		return TokensPair{}, err
	}

	return TokensPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) ValidateToken(tokenString string, tokenType string) (*jwt.Token, error) {
	var secret []byte

	switch tokenType {
	case "access":
		secret = []byte(os.Getenv("JWT_ACCESS_SECRET_KEY"))
	case "refresh":
		secret = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))
	default:
		return nil, fmt.Errorf("unsupported token type %q", tokenType)
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("%s is not set", secretEnvironmentKey(tokenType))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return token, nil
}

func (s *Service) WriteAccessToken(request AccessRequest) (string, error) {
	secret := []byte(os.Getenv("JWT_ACCESS_SECRET_KEY"))
	if len(secret) == 0 {
		return "", fmt.Errorf("JWT_ACCESS_SECRET_KEY is not set")
	}
	expiration, err := tokenExpiration("JWT_ACCESS_EXP", 10*time.Minute)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"maxId":             request.MaxId,
			"id":                request.Id,
			"role":              request.Role,
			"applicationStatus": request.ApplicationStatus,
			"organizationId":    request.OrganizationId,
			"exp":               time.Now().Add(expiration).Unix(),
		})

	return token.SignedString(secret)
}

func (s *Service) WriteRefreshToken(request RefreshRequest) (string, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))
	if len(secret) == 0 {
		return "", fmt.Errorf("JWT_REFRESH_SECRET_KEY is not set")
	}
	expiration, err := tokenExpiration("JWT_REFRESH_EXP", 5*24*time.Hour)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"maxId": request.MaxId,
			"id":    request.Id,
			"exp":   time.Now().Add(expiration).Unix(),
		})

	return token.SignedString(secret)
}

func tokenExpiration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	expiration, err := time.ParseDuration(value)
	if err != nil || expiration <= 0 {
		return 0, fmt.Errorf("invalid %s: expected a positive Go duration", key)
	}
	return expiration, nil
}

func secretEnvironmentKey(tokenType string) string {
	if tokenType == "access" {
		return "JWT_ACCESS_SECRET_KEY"
	}
	return "JWT_REFRESH_SECRET_KEY"
}
