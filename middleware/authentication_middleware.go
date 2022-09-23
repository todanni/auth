package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/todanni/token"
)

type ContextKey string

var (
	MissingFieldError = errors.New("missing field in token")
)

const (
	UserInfoContextKey    ContextKey = "userInfo"
	AccessTokenContextKey ContextKey = "accessToken"
)

type AuthenticationMiddleware struct {
	handler http.HandlerFunc
}

func (m *AuthenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if cookie is set
	accessTokenCookie, err := r.Cookie(token.AccessTokenCookieName)
	if err != nil {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}

	accessToken := token.NewAccessToken()
	err = accessToken.Validate(accessTokenCookie.Value)

	switch err {
	case MissingFieldError:
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	case nil:
		break
	default:
		http.Error(w, "invalid keys", http.StatusForbidden)
		return
	}

	ctx := context.WithValue(r.Context(), AccessTokenContextKey, accessToken)
	m.handler(w, r.WithContext(ctx))
}

func NewAuthenticationMiddleware(handlerToWrap http.HandlerFunc) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{handlerToWrap}
}
