package storage

import (
	"encoding/json"
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/database"
)

func Test_dashboardStorage_Create(t *testing.T) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = "https://vault.todanni.com"

	vaultClient, err := vault.NewClient(vaultConfig)
	require.NoError(t, err)
	vaultClient.SetToken(os.Getenv("VAULT_SIGNING_KEYS_TOKEN"))

	// Read oauthConfig
	cfg, err := config.NewFromEnv(vaultClient)
	require.NoError(t, err)

	// Open database connection
	db, err := database.Open(cfg)
	require.NoError(t, err)

	s := dashboardStorage{
		db: db,
	}

	result, err := s.Create(4, 3)
	require.NoError(t, err)
	require.NotNil(t, result)

	marshalled, err := json.Marshal(result)
	require.NotNil(t, marshalled)
}

func Test_dashboardStorage_List(t *testing.T) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = "https://vault.todanni.com"

	vaultClient, err := vault.NewClient(vaultConfig)
	require.NoError(t, err)
	vaultClient.SetToken(os.Getenv("VAULT_SIGNING_KEYS_TOKEN"))

	// Read oauthConfig
	cfg, err := config.NewFromEnv(vaultClient)
	require.NoError(t, err)

	// Open database connection
	db, err := database.Open(cfg)
	require.NoError(t, err)

	s := dashboardStorage{
		db: db,
	}

	result, err := s.List(4)
	require.NoError(t, err)
	require.NotNil(t, result)

	marshalled, err := json.Marshal(result)
	require.NoError(t, err)
	require.NotNil(t, marshalled)
}
