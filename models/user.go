package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email      string `json:"email"`
	LoginType  string `json:"loginType"`
	ProfilePic string `json:"profilePic"`
}

type UserInfo struct {
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	UserID     uint   `json:"userID"`
}
