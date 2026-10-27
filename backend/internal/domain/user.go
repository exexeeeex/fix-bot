package domain

import "time"

type User struct {
	Id        string    `json:"id"`
	MaxId     string    `json:"max_id"`
	RoleId    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
