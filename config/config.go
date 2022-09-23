package config

import (
	"context"

	"github.com/caarlos0/env/v6"
	vault "github.com/hashicorp/vault/api"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"

	"github.com/todanni/auth/keys"
)

// Config contains the env variables needed to run the servers
type Config struct {
	DBHost       string `env:"POSTGRES_HOST,required"`
	DBPort       int    `env:"POSTGRES_PORT,required"`
	DBUser       string `env:"POSTGRES_USER,required"`
	DBPassword   string `env:"POSTGRES_PASSWORD,required"`
	DBName       string `env:"POSTGRES_NAME,required"`
	PrivateKeyID string `env:"PRIVATE_KEY_ID,required"`

	vault      *vault.Client
	PrivateJWK jwk.Key
	PublicJWK  jwk.Key
}

func NewFromEnv(vault *vault.Client) (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	privateKey, err := keys.GetSigningKey(context.Background(), vault, config.PrivateKeyID)
	if err != nil {
		return Config{}, err
	}

	config.PrivateJWK, _ = jwk.New(privateKey)
	config.PrivateJWK.Set(jwk.KeyIDKey, config.PrivateKeyID)
	config.PrivateJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	config.PublicJWK, _ = jwk.New(privateKey.PublicKey)
	config.PublicJWK.Set(jwk.KeyIDKey, config.PrivateKeyID)
	config.PublicJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	return config, nil
}
