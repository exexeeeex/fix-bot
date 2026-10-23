package domain

import "time"

type Teacher struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationId string    `json:"organization_id"`
	RoleId         string    `json:"role_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
