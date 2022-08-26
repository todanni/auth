package config

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/caarlos0/env/v6"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/thanhpk/randstr"
)

// Config contains the env variables needed to run the servers
type Config struct {
	DBHost     string `env:"POSTGRES_HOST,required"`
	DBPort     int    `env:"POSTGRES_PORT,required"`
	DBUser     string `env:"POSTGRES_USER,required"`
	DBPassword string `env:"POSTGRES_PASSWORD,required"`
	DBName     string `env:"POSTGRES_NAME,required"`
	PrivateJWK jwk.Key
	PublicJWK  jwk.Key
}

func NewFromEnv() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Config{}, err
	}

	keyID := randstr.Hex(10)
	config.PrivateJWK, _ = jwk.New(privateKey)
	config.PrivateJWK.Set(jwk.KeyIDKey, keyID)
	config.PrivateJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	config.PublicJWK, _ = jwk.New(privateKey.PublicKey)
	config.PublicJWK.Set(jwk.KeyIDKey, keyID)
	config.PublicJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	return config, nil
}
