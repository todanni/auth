package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/database"
	"github.com/todanni/auth/service/auth"
	"github.com/todanni/auth/storage"
)

func main() {
	// Read config
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Open database connection
	db, err := database.Open(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Initialise router
	r := mux.NewRouter()

	osinServer := osin.NewServer(&osin.ServerConfig{
		AllowedAuthorizeTypes:     osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN},
		AllowedAccessTypes:        osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.REFRESH_TOKEN},
		AllowGetAccessRequest:     true,
		AllowClientSecretInParams: true,
	}, storage.New(db))

	auth.NewAuthService(osinServer, r)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
