package storage

import (
	"gorm.io/gorm"

	"github.com/todanni/auth/models"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return UserStorage{
		db: db,
	}
}

func (storage *UserStorage) CreateUser(email, loginType string) {
	storage.db.Create(models.User{
		Email:     email,
		LoginType: loginType,
	})
}

func (storage *UserStorage) GetUser(email string) (models.User, error) {
	var result models.User
	storage.db.Where("email = ?", email).First(&result)
	return result, nil
}
