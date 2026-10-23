package domain

import "time"

type Group struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationId string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
