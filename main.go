package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/database"
	"github.com/todanni/auth/models"
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

	// Perform migrations
	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		log.Fatalf("couldn't auto migrate: %v", err)
	}

	strg := storage.NewUserStorage(db)

	// Initialise router
	r := mux.NewRouter()

	// Create HTTP service
	auth.NewAuthService(r, cfg, &strg)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
