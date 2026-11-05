package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("JWT_ACCESS_SECRET_KEY"))
var refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))

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
		secret = accessSecret
	case "refresh":
		secret = refreshSecret
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"maxId":             request.MaxId,
			"id":                request.Id,
			"role":              request.Role,
			"applicationStatus": request.ApplicationStatus,
			"exp":               time.Now().Add(time.Duration(10) * time.Minute).Unix(),
		})

	return token.SignedString(accessSecret)
}

func (s *Service) WriteRefreshToken(request RefreshRequest) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"maxId": request.MaxId,
			"id":    request.Id,
			"exp":   time.Now().Add(5 * 24 * time.Hour).Unix(),
		})

	return token.SignedString(refreshSecret)
}
