package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/todanni/auth/models"
	"github.com/todanni/auth/test"
	"github.com/todanni/auth/token"
)

type AuthenticationCheckTestSuite struct {
	suite.Suite
	privateKey jwk.Key
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func (s *AuthenticationCheckTestSuite) SetupSuite() {
	s.privateKey = test.ServePublicKey()
}

func (s *AuthenticationCheckTestSuite) Test_AuthenticationCheck_Bad_NoCookie401() {
	handler := NewAuthenticationMiddleware(dummyHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 401)
	require.Equal(s.T(), rw.Body.String(), "unauthorised\n")
}

func (s *AuthenticationCheckTestSuite) Test_AuthenticationCheck_Bad_InvalidToken403() {
	handler := NewAuthenticationMiddleware(dummyHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:   "todanni-access-token",
		Value:  "token",
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 403)
	require.Equal(s.T(), rw.Body.String(), "invalid token\n")
}

func (s *AuthenticationCheckTestSuite) Test_AuthenticationCheck_Good() {
	user := models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:      "test@test.com",
		ProfilePic: "",
	}

	dashboards := make([]models.Dashboard, 0)
	projects := make([]models.Project, 0)
	userInfo := models.UserInfo{
		Email:      user.Email,
		ProfilePic: user.ProfilePic,
		UserID:     user.ID,
	}

	todanniToken := token.NewAccessToken()
	todanniToken.SetUserInfo(userInfo).
		SetProjectsPermissions(projects).
		SetDashboardPermissions(dashboards)

	signedToken, err := todanniToken.SignedToken(s.privateKey)
	require.NoError(s.T(), err)

	handler := NewAuthenticationMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(AccessTokenContextKey).(models.UserInfo)
		if userInfo.UserID != user.ID {
			s.T().Errorf("user info incorrect, expected %v but got %v", user.ID, userInfo.UserID)
		}
	})

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:   "todanni-access-token",
		Value:  signedToken,
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	handler.ServeHTTP(rw, req)

	require.Equal(s.T(), rw.Code, 200)
}

func TestAuthenticationCheckTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationCheckTestSuite))
}
