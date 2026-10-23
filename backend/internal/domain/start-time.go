package domain

import "time"

type StartTime struct {
	Id        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
