package maxauth

import (
	"github.com/exexeeeex/fiks-api/internal/domain"
	jwt "github.com/exexeeeex/fiks-api/internal/jwt"
)

type MaxUser struct {
	Id           int     `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Username     string  `json:"username"`
	LanguageCode *string `json:"language_code"`
	PhotoURL     *string `json:"photo_url"`
}

type AuthenticatedUser struct {
	Id                string                   `json:"id"`
	MaxId             string                   `json:"max_id"`
	FirstName         string                   `json:"first_name"`
	LastName          string                   `json:"last_name"`
	Username          string                   `json:"username"`
	LanguageCode      *string                  `json:"language_code"`
	PhotoURL          *string                  `json:"photo_url"`
	ApplicationStatus domain.ApplicationStatus `json:"applicationStatus"`
}

type AuthResponse struct {
	Success bool              `json:"success"`
	User    AuthenticatedUser `json:"user"`
	Tokens  jwt.TokensPair    `json:"tokens"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type AuthRequest struct {
	InitData string `json:"initData"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
