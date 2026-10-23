package domain

import "time"

type Visit struct {
	Id        string    `json:"id"`
	StudentId string    `json:"student_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
