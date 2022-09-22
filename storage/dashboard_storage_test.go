package storage

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/todanni/auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DashboardStorageTestSuite struct {
	suite.Suite
	db          *gorm.DB
	cleanupFunc func()
}

func (s *DashboardStorageTestSuite) SetupSuite() {
	s.db, s.cleanupFunc = setupGormWithDocker()
	s.db.AutoMigrate(&models.Dashboard{}, &models.User{})
}

func (s *DashboardStorageTestSuite) TearDownSuite() {
	s.cleanupFunc()
}

func (s *DashboardStorageTestSuite) Test_dashboardStorage_Create() {
	storage := dashboardStorage{
		db: s.db,
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
		db: s.db,
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

const (
	dbName = "test"
	passwd = "test"
)

func setupGormWithDocker() (*gorm.DB, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "latest",   // version
		Env:        []string{"POSTGRES_PASSWORD=" + passwd, "POSTGRES_DB=" + dbName},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	conStr := fmt.Sprintf("host=localhost port=%s user=postgres dbname=%s password=%s sslmode=disable",
		resource.GetPort("5432/tcp"), // get port of localhost
		dbName,
		passwd,
	)

	var gdb *gorm.DB
	// retry until db server is ready
	err = pool.Retry(func() error {
		gdb, err = gorm.Open(postgres.Open(conStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	})
	chk(err)

	// container is ready, return *gorm.Db for testing
	return gdb, fnCleanup
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
