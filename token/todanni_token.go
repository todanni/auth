package token

import (
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/exp/slices"

	"github.com/todanni/auth/models"
)

type ToDanniToken struct {
	token jwt.Token
}

func NewToDanniToken() (*ToDanniToken, error) {
	t, err := jwt.NewBuilder().Issuer(ToDanniTokenIssuer).IssuedAt(time.Now()).Build()
	if err != nil {
		return nil, err
	}
	return &ToDanniToken{token: t}, nil
}

func (t *ToDanniToken) Validate(token string, keySet jwk.Set) error {
	parsed, err := jwt.Parse([]byte(token),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.WithTypedClaim("userID", uint(1)))

	if err != nil {
		return err
	}

	t.token = parsed
	return nil
}

func (t *ToDanniToken) SignToken(privateKey jwk.Key) (string, error) {
	signedJWT, err := jwt.Sign(t.token, jwa.RS256, privateKey)
	if err != nil {
		return "", err
	}

	return string(signedJWT), nil
}

func (t *ToDanniToken) HasDashboardPermission(dashboard string) bool {
	dashboardsPermissions, ok := t.token.Get("dashboards")
	if !ok {
		return false
	}

	dashboardsPermissionsArray := dashboardsPermissions.([]string)
	return slices.Contains(dashboardsPermissionsArray, dashboard)
}

func (t *ToDanniToken) HasProjectPermission(project string) bool {
	projectPermissions, ok := t.token.Get("projects")
	if !ok {
		return false
	}

	dashboardsPermissionsArray := projectPermissions.([]string)
	return slices.Contains(dashboardsPermissionsArray, project)
}

func (t *ToDanniToken) GetUserInfo() (models.UserInfo, error) {
	userInfo, ok := t.token.Get("user-info")
	if !ok {
		return models.UserInfo{}, errors.New("token doesn't contain user info")
	}

	return userInfo.(models.UserInfo), nil
}

func (t *ToDanniToken) SetUserInfo(userInfo models.UserInfo) {
	t.token.Set("user-info", userInfo)
}

func (t *ToDanniToken) SetDashboardPermissions(dashboards []models.Dashboard) {
	userDashboardIDs := make([]uint, 0)

	for _, dashboard := range dashboards {
		userDashboardIDs = append(userDashboardIDs, dashboard.ID)
	}

	t.token.Set("dashboards", userDashboardIDs)
}

func (t *ToDanniToken) SetProjectsPermissions(projects []models.Project) {
	userProjectIDs := make([]uint, 0)

	for _, project := range projects {
		userProjectIDs = append(userProjectIDs, project.ID)
	}

	t.token.Set("projects", userProjectIDs)
}
