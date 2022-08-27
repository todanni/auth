package main

import (
	b64 "encoding/base64"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	scopes "google.golang.org/api/oauth2/v2"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/database"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/service/auth"
	"github.com/todanni/auth/storage"
)

func main() {
	// Read oauthConfig
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

	// Create storage with the DB connection
	strg := storage.NewUserStorage(db)

	// Initialise router
	r := mux.NewRouter()

	// Create OAuth oauthConfig
	googleCredentials := os.Getenv("GOOGLE_CREDENTIALS")
	decodedCredentials, err := b64.StdEncoding.DecodeString(googleCredentials)

	oauthConfig, err := google.ConfigFromJSON(decodedCredentials, scopes.OpenIDScope, scopes.UserinfoEmailScope, scopes.UserinfoProfileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to oauthConfig: %v", err)
	}

	// Create HTTP service
	auth.NewAuthService(r, cfg, &strg, oauthConfig)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
