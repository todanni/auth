package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/storage"
	"github.com/todanni/auth/token"
)

type AuthService interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	RefreshTokenHandler(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	router  *mux.Router
	server  *osin.Server
	storage *storage.UserStorage
	config  config.Config
}

const (
	AccessTokenCookieName  = "todanni-access-token"
	RefreshTokenCookieName = "todanni-refresh-token"
)

func NewAuthService(router *mux.Router, conf config.Config, strg *storage.UserStorage) AuthService {
	server := &authService{
		router:  router,
		config:  conf,
		storage: strg,
	}
	server.routes()
	return server
}

func (s *authService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *authService) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	config := oauth2.Config{}
	ctx := context.Background()

	code := r.URL.Query().Get("code")

	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	idToken := tok.Extra("id_token").(string)
	email, err := token.ValidateGoogleToken(ctx, idToken)
	if err != nil {
		http.Error(w, "invalid Google token", http.StatusBadRequest)
	}
	// TODO: Check if user exists

	// TODO: translate google token to ToDanni token
	accessToken, err := token.IssueToDanniToken(email, s.config.PrivateJWK)
	if err != nil {
		http.Error(w, "couldn't create the ToDanni token", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	// Set access and refresh token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	})

	refreshToken, err := token.IssueToDanniRefreshToken(1)
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken.Value,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	})
	http.Redirect(w, r, "https://todanni.com/", http.StatusFound)
}

func (s *authService) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Check that the provided refresh token is valid by querying the database

	accessToken, err := token.IssueToDanniToken("email", s.config.PrivateJWK)
	if err != nil {
		http.Error(w, "couldn't issue refresh token", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: 2,
	})
}

func (s *authService) ServePublicKey(w http.ResponseWriter, r *http.Request) {
	keyset := jwk.NewSet()
	keyset.Add(s.config.PublicJWK)

	buf, err := json.Marshal(keyset)
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed to marshal key", http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(buf)
}
