package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name    string `json:"name"`
	Owner   uint   `json:"owner"`
	Members []User `json:"members" gorm:"many2many:user_projects;"`
}
