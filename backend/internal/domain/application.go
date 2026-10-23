package domain

import "time"

type Application struct {
	Id        	string 		`json:"id"`
	TeacherId 	string 		`json:"teacher_id"`
	CreatedAt 	time.Time	`json:"created_at"`
	Updatedat	time.Time	`json:"updated_at"`
}