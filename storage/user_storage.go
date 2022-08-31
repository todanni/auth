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

func (storage *UserStorage) CreateUser(email, loginType, profilePic string) (models.User, error) {
	user := models.User{
		Email:     email,
		LoginType: loginType,
		ProfilePic: profilePic,
	}
	result := storage.db.Create(&user)
	return user, result.Error
}

func (storage *UserStorage) GetUser(email string) (models.User, error) {
	var result models.User
	storage.db.Where("email = ?", email).First(&result)
	return result, nil
}
