package models

import "gorm.io/gorm"

type Dashboard struct {
	Owner   uint   `json:"owner"`
	Members []uint `json:"members"`
	Status  Status `json:"status"`
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
