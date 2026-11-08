package application

import (
	domain "github.com/exexeeeex/fiks-api/internal/domain"
)

type createApplicationRequest struct {
	Name           string                 `json:"name"`
	OrganizationId string                 `json:"organization_id"`
	GroupId        *string                `json:"group_id"`
	UserId         string                 `json:"user_id"`
	Type           domain.ApplicationType `json:"type"`
}

type updateApplicationRequest struct {
	Id     string                   `json:"id"`
	Status domain.ApplicationStatus `json:"status"`
}
