package test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/thanhpk/randstr"
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
