package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/test"
)

type DashboardStorageTestSuite struct {
	test.DbSuite
}

func (s *DashboardStorageTestSuite) SetupSuite() {
	s.Db, s.CleanupFunc = test.SetupGormWithDocker()
	s.Db.AutoMigrate(&models.Dashboard{}, &models.User{})
}

func (s *DashboardStorageTestSuite) TearDownSuite() {
	s.CleanupFunc()
}

func (s *DashboardStorageTestSuite) Test_dashboardStorage_Create() {
	storage := dashboardStorage{
		db: s.Db,
	}

	result, err := storage.Create(4, 3)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	marshalled, err := json.Marshal(result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), marshalled)
}

func (s *DashboardStorageTestSuite) Test_dashboardStorage_List() {
	storage := dashboardStorage{
		db: s.Db,
	}

	result, err := storage.List(4)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	marshalled, err := json.Marshal(result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), marshalled)
}

func TestDashboardStorageTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardStorageTestSuite))
}
