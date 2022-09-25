package main

import (
	b64 "encoding/base64"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	scopes "google.golang.org/api/oauth2/v2"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/database"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/service/auth"
	"github.com/todanni/auth/service/dashboard"
	"github.com/todanni/auth/service/project"
	"github.com/todanni/auth/storage"
)

const (
	VaultAddress = "https://vault.todanni.com"
)

func main() {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = VaultAddress

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("unable to initialize a Vault client: %v", err)
	}
	vaultClient.SetToken(os.Getenv("VAULT_SIGNING_KEYS_TOKEN"))

	// Read oauthConfig
	cfg, err := config.NewFromEnv(vaultClient)
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
	err = db.AutoMigrate(&models.User{}, &models.Dashboard{}, &models.Project{})
	if err != nil {
		log.Fatalf("couldn't auto migrate: %v", err)
	}

	// Create storage with the DB connection
	userStorage := storage.NewUserStorage(db)
	dashboardStorage := storage.NewDashboardStorage(db)
	projectStorage := storage.NewProjectStorage(db)

	// Initialise router
	r := mux.NewRouter()

	// Create OAuth oauthConfig
	googleCredentials := os.Getenv("GOOGLE_CREDENTIALS")
	decodedCredentials, err := b64.StdEncoding.DecodeString(googleCredentials)

	oauthConfig, err := google.ConfigFromJSON(decodedCredentials, scopes.OpenIDScope, scopes.UserinfoEmailScope, scopes.UserinfoProfileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to oauthConfig: %v", err)
	}

	// Create auth service
	auth.NewAuthService(r, cfg, userStorage, dashboardStorage, projectStorage, oauthConfig)

	// Create dashboard service
	dashboard.NewDashboardService(storage.NewDashboardStorage(db), userStorage, r)

	// Create project service
	project.NewProjectService(r, projectStorage)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
