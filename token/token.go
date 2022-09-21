package token

import (
	"context"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/thanhpk/randstr"

	"github.com/todanni/auth/models"
)

const (
	GoogleCertsUrl = "https://www.googleapis.com/oauth2/v3/certs"
	//TODO: why is this localhost?
	ToDanniCertsUrl            = "http://localhost:8083/auth/public-key"
	ToDanniTokenIssuer         = "todanni.com"
	RefreshTokenExpirationTime = time.Hour * 60 * 30
)

var (
	MissingFieldError = errors.New("missing field in token")
)

// ValidateGoogleToken follows the OAuth 2.0 spec to validate token
// and returns the email of the user the token belongs to
func ValidateGoogleToken(ctx context.Context, tkn string) (string, error) {
	autoRefresh := jwk.NewAutoRefresh(ctx)
	autoRefresh.Configure(GoogleCertsUrl, jwk.WithMinRefreshInterval(time.Hour*1))

	keySet, err := autoRefresh.Fetch(ctx, GoogleCertsUrl)
	if err != nil {
		return "", err
	}

	parsed, err := jwt.Parse([]byte(tkn), jwt.WithKeySet(keySet), jwt.WithValidate(true))
	if err != nil {
		return "", err
	}

	email, ok := parsed.Get("email")
	if !ok {
		return "", errors.New("couldn't find email in token")
	}

	return email.(string), nil
}

func IssueToDanniToken(user models.User, privateKey jwk.Key, dashboards []models.Dashboard, projects []models.Project) (string, error) {
	t, err := jwt.NewBuilder().Issuer(ToDanniTokenIssuer).IssuedAt(time.Now()).Build()
	if err != nil {
		return "", err
	}

	// Set the custom claims
	t.Set("email", user.Email)
	t.Set("userID", user.ID)
	t.Set("profilePic", user.ProfilePic)

	SetDashboardsPermissions(dashboards, t)
	SetProjectPermissions(projects, t)

	signedJWT, err := jwt.Sign(t, jwa.RS256, privateKey)
	if err != nil {
		return "", err
	}

	return string(signedJWT), nil
}

func IssueToDanniRefreshToken(userID int) (models.RefreshToken, error) {
	refreshToken := models.RefreshToken{
		Value:     randstr.Hex(10),
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(RefreshTokenExpirationTime),
	}
	return refreshToken, nil
}

// ValidateToDanniToken checks the token provided is issued by todanni
// and returns the email of the user it belongs to
func ValidateToDanniToken(token string) (models.UserInfo, error) {
	var userInfo models.UserInfo

	autoRefresh := jwk.NewAutoRefresh(context.Background())
	autoRefresh.Configure(ToDanniCertsUrl, jwk.WithMinRefreshInterval(time.Hour*1))
	keySet, err := autoRefresh.Fetch(context.Background(), ToDanniCertsUrl)
	if err != nil {
		return userInfo, err
	}

	parsed, err := jwt.Parse([]byte(token), jwt.WithKeySet(keySet), jwt.WithValidate(true), jwt.WithTypedClaim("userID", uint(1)))
	if err != nil {
		return userInfo, err
	}

	email, ok := parsed.Get("email")
	if !ok {
		return userInfo, MissingFieldError
	}
	userInfo.Email = email.(string)

	userid, ok := parsed.Get("userID")
	if !ok {
		return userInfo, MissingFieldError
	}

	userInfo.UserID = userid.(uint)

	profilePic, ok := parsed.Get("profilePic")
	if !ok {
		return userInfo, MissingFieldError
	}
	userInfo.ProfilePic = profilePic.(string)

	return userInfo, err
}
