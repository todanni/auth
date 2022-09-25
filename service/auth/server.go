package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"github.com/todanni/token"
	"golang.org/x/oauth2"

	"github.com/todanni/auth/config"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/storage"
)

type AuthService interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	RefreshTokenHandler(w http.ResponseWriter, r *http.Request)
	UserInfoHandler(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	router           *mux.Router
	server           *osin.Server
	userStorage      storage.UserStorage
	dashboardStorage storage.DashboardStorage
	projectStorage   storage.ProjectStorage
	config           config.Config
	oauthConfig      *oauth2.Config
}

func NewAuthService(router *mux.Router, conf config.Config, userStorage storage.UserStorage, dashboardStorage storage.DashboardStorage, projectStorage storage.ProjectStorage, oauthConfig *oauth2.Config) AuthService {
	server := &authService{
		oauthConfig:      oauthConfig,
		config:           conf,
		router:           router,
		userStorage:      userStorage,
		dashboardStorage: dashboardStorage,
		projectStorage:   projectStorage,
	}
	server.routes()
	return server
}

func (s *authService) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Received callback request")
	ctx := context.Background()

	code := r.URL.Query().Get("code")
	log.Info(s.oauthConfig)
	log.Info(s.oauthConfig.RedirectURL)

	tok, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Errorf("Couldn't exchange keys: %v", err)
		http.Error(w, "couldn't exchange keys for code", http.StatusInternalServerError)
		return
	}

	idToken := tok.Extra("id_token").(string)
	log.Info(idToken)

	email, err := s.validateGoogleToken(ctx, idToken)
	if err != nil {
		log.Errorf("Couldn't validate keys: %v", err)
		http.Error(w, "invalid Google keys", http.StatusBadRequest)
		return
	}

	// Check if user exists
	userRecord, err := s.userStorage.GetUser(email)
	if err != nil {
		log.Errorf("Couldn't check if user exists: %v", err)
		http.Error(w, "some error with user", http.StatusInternalServerError)
		return
	}

	// User doesn't exist, we have to create it
	if userRecord.ID == 0 {
		userRecord, err = s.userStorage.CreateUser(email, "google", "https://www.dictionary.com/e/wp-content/uploads/2018/03/rickrolling-300x300.jpg")
		if err != nil {
			log.Errorf("Couldn't create user: %v", err)
			http.Error(w, "couldn't create new user", http.StatusInternalServerError)
			return
		}
	}

	userInfo := models.UserInfo{
		Email:      userRecord.Email,
		ProfilePic: userRecord.ProfilePic,
		UserID:     userRecord.ID,
	}
	dashboards := make([]models.Dashboard, 0)
	projects := make([]models.Project, 0)

	if userRecord.ID != 0 {
		dashboards, err = s.dashboardStorage.List(userRecord.ID)
		if err != nil {
			log.Error("couldn't look up user dashboards")
		}

		projects, err = s.projectStorage.List(userRecord.ID)
		if err != nil {
			log.Error("couldn't look up user dashboards")
		}
	}

	todanniToken := token.NewAccessToken()
	todanniToken.SetUserInfo(userInfo).
		SetProjectsPermissions(projects).
		SetDashboardPermissions(dashboards)

	signedToken, err := todanniToken.SignedToken(s.config.PrivateJWK)
	if err != nil {
		log.Errorf("Couldn't sign todanni keys: %v", err)
		http.Error(w, "couldn't create the ToDanni keys", http.StatusInternalServerError)
		return
	}

	// Set access and refresh keys cookies
	http.SetCookie(w, &http.Cookie{
		Name:     token.AccessTokenCookieName,
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,
	})

	// TODO: Generate and persist a refresh keys
	http.SetCookie(w, &http.Cookie{
		Name:     token.RefreshTokenCookieName,
		Value:    "madeup-keys-todo",
		Path:     "/",
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (s *authService) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Check that the provided refresh keys is valid by querying the database
	var user models.User

	dashboards, err := s.dashboardStorage.List(user.ID)
	if err != nil {
		log.Error("couldn't look up user dashboards")
	}

	projects, err := s.projectStorage.List(user.ID)
	if err != nil {
		log.Error("couldn't look up user dashboards")
	}

	userInfo := models.UserInfo{
		Email:      user.Email,
		ProfilePic: user.ProfilePic,
		UserID:     user.ID,
	}

	todanniToken := token.NewAccessToken()
	todanniToken.SetUserInfo(userInfo).
		SetProjectsPermissions(projects).
		SetDashboardPermissions(dashboards)

	signedToken, err := todanniToken.SignedToken(s.config.PrivateJWK)
	if err != nil {
		log.Errorf("Couldn't sign todanni keys: %v", err)
		http.Error(w, "couldn't create the ToDanni keys", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     token.AccessTokenCookieName,
		Value:    signedToken,
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
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(buf)
}

func (s *authService) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	log.Info(accessToken)

	userInfo, err := accessToken.GetUserInfo()
	if err != nil {
		http.Error(w, "keys didn't contain user info", http.StatusBadRequest)
	}

	marshalled, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, "couldn't marshal keys", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(marshalled)
}

func (s *authService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unimplemented", http.StatusMethodNotAllowed)
}

func (s *authService) validateGoogleToken(ctx context.Context, tkn string) (string, error) {
	autoRefresh := jwk.NewAutoRefresh(ctx)
	autoRefresh.Configure("https://www.googleapis.com/oauth2/v3/certs", jwk.WithMinRefreshInterval(time.Hour*1))

	keySet, err := autoRefresh.Fetch(ctx, "https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return "", err
	}

	parsed, err := jwt.Parse([]byte(tkn), jwt.WithKeySet(keySet), jwt.WithValidate(true))
	if err != nil {
		return "", err
	}

	email, ok := parsed.Get("email")
	if !ok {
		return "", errors.New("couldn't find email in keys")
	}

	return email.(string), nil
}
