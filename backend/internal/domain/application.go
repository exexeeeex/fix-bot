package domain

import "time"

type ApplicationType string
type ApplicationStatus string

const (
	ApplicationTypeStudent ApplicationType = "student"
	ApplicationTypeTeacher ApplicationType = "teacher"
)

const (
	ApplicationStatusNone     ApplicationStatus = "none"
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusApproved ApplicationStatus = "approved"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

type Application struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Status         ApplicationStatus `json:"status"`
	Type           ApplicationType   `json:"type"`
	UserId         string            `json:"user_id"`
	OrganizationId string            `json:"organization_id"`
	GroupId        *string           `json:"group_id"`
	Group          *Group            `json:"group,omitempty" gorm:"foreignKey:GroupId;references:Id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
