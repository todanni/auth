package storage

import (
	"errors"

	"gorm.io/gorm"

	"github.com/todanni/auth/models"
)

type ProjectStorage interface {
	Create(owner, invited uint) (models.Project, error)
	List(userid uint) ([]models.Project, error)
	GetById(id string) (models.Project, error)
	UpdateStatus(id string, status models.Status) (models.Project, error)
	Delete(id string) error
}

func NewProjectStorage(db *gorm.DB) ProjectStorage {
	return &projectStorage{
		db: db,
	}
}

type projectStorage struct {
	db *gorm.DB
}

func (s *projectStorage) List(userid uint) ([]models.Project, error) {
	var Projects []models.Project
	var user models.User

	result := s.db.Model(&models.User{}).Preload("Projects.Members").First(&user, userid)
	if result.Error != nil {
		return Projects, errors.New("couldn't find Projects")
	}

	return user.Projects, nil
}

func (s *projectStorage) Create(owner, invited uint) (models.Project, error) {
	Project := models.Project{
		Owner: owner,
	}

	result := s.db.Create(&Project)
	if result.Error != nil {
		return models.Project{}, errors.New("couldn't create the Project")
	}

	users := []models.User{{
		Model: gorm.Model{
			ID: owner,
		},
	},
		{
			Model: gorm.Model{
				ID: invited,
			},
		},
	}

	err := s.db.Model(&Project).Association("Members").Append(users)
	if err != nil {
		return Project, err
	}

	return Project, nil
}

func (s *projectStorage) GetById(id string) (models.Project, error) {
	var Project models.Project
	result := s.db.First(&Project, id)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		return Project, errors.New("this Project doesn't exist")
	case nil:
		return Project, nil
	default:
		return models.Project{}, errors.New("couldn't get Project: " + result.Error.Error())
	}
}

func (s *projectStorage) UpdateStatus(id string, status models.Status) (models.Project, error) {
	var Project models.Project

	result := s.db.First(&Project, id).Update("status", status)
	if result.Error != nil {
		return Project, result.Error
	}

	return Project, nil
}

func (s *projectStorage) Delete(id string) error {
	result := s.db.Delete(&models.Project{}, id)
	return result.Error
}
