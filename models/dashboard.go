package models

import (
	"gorm.io/gorm"
)

type Dashboard struct {
	Owner   uint   `json:"owner"`
	Status  Status `json:"status"`
	Members []User `json:"members" gorm:"many2many:user_dashboards;"`
	gorm.Model
}

type Status string

const (
	PendingStatus  Status = "PENDING"
	AcceptedStatus Status = "ACCEPTED"
	RejectedStatus Status = "REJECTED"
)

type DashboardCreateRequest struct {
	Email string `json:"email"`
}
