package middleware

import (
	"context"
	"net/http"

	"github.com/todanni/auth/token"
)

type ContextKey string

const (
	UserInfoContextKey    ContextKey = "userInfo"
	AccessTokenContextKey ContextKey = "accessToken"
)

type AuthenticationCheck struct {
	handler http.HandlerFunc
}

func (ea *AuthenticationCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if cookie is set
	accessTokenCookie, err := r.Cookie(token.AccessTokenCookieName)
	if err != nil {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}

	accessToken := token.NewAccessToken()
	err = accessToken.Validate(accessTokenCookie.Value)

	switch err {
	case token.MissingFieldError:
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	case nil:
		break
	default:
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	ctx := context.WithValue(r.Context(), AccessTokenContextKey, accessToken)
	ea.handler(w, r.WithContext(ctx))
}

func NewAuthenticationCheck(handlerToWrap http.HandlerFunc) *AuthenticationCheck {
	return &AuthenticationCheck{handlerToWrap}
}
