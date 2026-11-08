package jwt

import (
	application "github.com/exexeeeex/fiks-api/internal/domain"
)

type AccessRequest struct {
	Id                string                        `json:"id"`
	MaxId             string                        `json:"max_id"`
	Role              string                        `json:"role"`
	ApplicationStatus application.ApplicationStatus `json:"applicationStatus"`
	OrganizationId    string                        `json:"organizationId"`
}

type RefreshRequest struct {
	Id    string `json:"id"`
	MaxId string `json:"max_id"`
}

type TokensPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
