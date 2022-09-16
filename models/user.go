package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email      string      `json:"email"`
	LoginType  string      `json:"loginType"`
	ProfilePic string      `json:"profilePic"`
	Dashboards []Dashboard `json:"-" gorm:"many2many:user_dashboards;"`
}

type UserInfo struct {
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	UserID     uint   `json:"userID"`
}
