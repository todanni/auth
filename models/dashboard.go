package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Dashboard struct {
	Owner   uint          `json:"owner"`
	Members pq.Int64Array `json:"members" gorm:"type:integer[]"`
	Status  Status        `json:"status"`
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
