package token

import (
	"context"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

const (
	GoogleCertsUrl = "https://www.googleapis.com/oauth2/v3/certs"
	// TODO: why is this localhost?
	ToDanniCertsUrl            = "http://localhost:8083/auth/public-key"
	ToDanniTokenIssuer         = "todanni.com"
	RefreshTokenExpirationTime = time.Hour * 60 * 30
	AccessTokenCookieName      = "todanni-access-token"
	RefreshTokenCookieName     = "todanni-refresh-token"
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
