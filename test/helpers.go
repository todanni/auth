package test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"github.com/thanhpk/randstr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ServePublicKey() jwk.Key {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyID := randstr.Hex(10)
	privateJWK, _ := jwk.New(privateKey)
	privateJWK.Set(jwk.KeyIDKey, keyID)
	privateJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	publicJWK, _ := jwk.New(privateKey.PublicKey)
	publicJWK.Set(jwk.KeyIDKey, keyID)
	publicJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyset := jwk.NewSet()
		keyset.Add(publicJWK)

		buf, err := json.Marshal(keyset)
		if err != nil {
			http.Error(w, "Failed to marshal key", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(buf)
	})

	http.Handle("/auth/public-key", handler)
	go http.ListenAndServe("localhost:8083", nil)

	return privateJWK
}

type DbSuite struct {
	suite.Suite
	Db          *gorm.DB
	CleanupFunc func()
}

const (
	dbName = "test"
	passwd = "test"
)

// Creates a temporary Postgres docker container to run tests against.
//
// Returns connected gorm DB and function to cleanup container after testing is complete.
func SetupGormWithDocker() (*gorm.DB, func()) {
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
