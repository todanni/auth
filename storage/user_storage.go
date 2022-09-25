package storage

import (
	"gorm.io/gorm"

	"github.com/todanni/auth/models"
)

type userStorage struct {
	db *gorm.DB
}

type UserStorage interface {
	CreateUser(email, loginType, profilePic string) (models.User, error)
	GetUser(email string) (models.User, error)
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &userStorage{
		db: db,
	}
}

func (s *userStorage) CreateUser(email, loginType, profilePic string) (models.User, error) {
	user := models.User{
		Email:      email,
		LoginType:  loginType,
		ProfilePic: profilePic,
	}
	result := s.db.Create(&user)
	return user, result.Error
}

func (s *userStorage) GetUser(email string) (models.User, error) {
	var result models.User
	s.db.Where("email = ?", email).First(&result)
	return result, nil
}
