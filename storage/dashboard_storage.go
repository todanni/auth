package storage

import (
	"errors"

	"gorm.io/gorm"

	"github.com/todanni/auth/models"
)

type DashboardStorage interface {
	Create(owner, invited uint) (models.Dashboard, error)
	List(userid uint) ([]models.Dashboard, error)
	GetById(id string) (models.Dashboard, error)
	UpdateStatus(id string, status models.Status) (models.Dashboard, error)
	Delete(id string) error
}

func NewDashboardStorage(db *gorm.DB) DashboardStorage {
	return &dashboardStorage{
		db: db,
	}
}

type dashboardStorage struct {
	db *gorm.DB
}

func (s *dashboardStorage) List(userid uint) ([]models.Dashboard, error) {
	var dashboards []models.Dashboard
	var user models.User

	result := s.db.Model(&models.User{}).Preload("Dashboards.Members").First(&user, userid)
	if result.Error != nil {
		return dashboards, errors.New("couldn't find dashboards")
	}

	return user.Dashboards, nil
}

func (s *dashboardStorage) Create(owner, invited uint) (models.Dashboard, error) {
	dashboard := models.Dashboard{
		Owner:  owner,
		Status: models.PendingStatus,
	}

	result := s.db.Create(&dashboard)
	if result.Error != nil {
		return models.Dashboard{}, errors.New("couldn't create the dashboard")
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

	err := s.db.Model(&dashboard).Association("Members").Append(users)
	if err != nil {
		return dashboard, err
	}

	return dashboard, nil
}

func (s *dashboardStorage) GetById(id string) (models.Dashboard, error) {
	var dashboard models.Dashboard
	result := s.db.First(&dashboard, id)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		return dashboard, errors.New("this dashboard doesn't exist")
	case nil:
		return dashboard, nil
	default:
		return models.Dashboard{}, errors.New("couldn't get dashboard: " + result.Error.Error())
	}
}

func (s *dashboardStorage) UpdateStatus(id string, status models.Status) (models.Dashboard, error) {
	var dashboard models.Dashboard

	result := s.db.First(&dashboard, id).Update("status", status)
	if result.Error != nil {
		return dashboard, result.Error
	}

	return dashboard, nil
}

func (s *dashboardStorage) Delete(id string) error {
	result := s.db.Delete(&models.Dashboard{}, id)
	return result.Error
}
