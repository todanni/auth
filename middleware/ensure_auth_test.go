package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thanhpk/randstr"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/token"
	"gorm.io/gorm"
)

type EnsureAuthTestSuite struct {
	suite.Suite
	key		jwk.Key
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func (s *EnsureAuthTestSuite) SetupSuite() {
	s.key = servePublicKey()
}

func (s *EnsureAuthTestSuite) Test_NoCookie_401() {
	handler := NewEnsureAuth(dummyHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 401)
	require.Equal(s.T(), rw.Body.String(), "unauthorised\n")
}

func (s *EnsureAuthTestSuite)Test_BadToken_403() {
	handler := NewEnsureAuth(dummyHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name: "todanni-access-token",
		Value: "token",
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 403)
	require.Equal(s.T(), rw.Body.String(), "invalid token\n")
}

func (s *EnsureAuthTestSuite)Test_ValidToken_200() {
	user := models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email: "test@test.com",
		ProfilePic: "",
	}

	token, _ := token.IssueToDanniToken(user, s.key)
	handler := NewEnsureAuth(func(w http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(UserInfoContextKey).(models.UserInfo)
		if userInfo.UserID != user.ID {
			s.T().Errorf("user info incorrect, expected %v but got %v", user.ID, userInfo.UserID)
		}
	})

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name: "todanni-access-token",
		Value: token,
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 200)
}

func TestEnsureAuthTestSuite(t *testing.T){
	suite.Run(t, new(EnsureAuthTestSuite))
}

func servePublicKey() jwk.Key {
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
		_, err = w.Write(buf)
	})

	http.Handle("/auth/public-key", handler)
	go http.ListenAndServe("localhost:8083", nil)

	

	return privateJWK
}