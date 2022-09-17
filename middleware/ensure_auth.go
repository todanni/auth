package middleware

import (
	"context"
	"net/http"

	"github.com/todanni/auth/token"
)

type ContextKey string

const (
	UserInfoContextKey ContextKey = "userInfo"
)

type EnsureAuth struct {
    handler http.HandlerFunc
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if cookie is set
	accessTokenCookie, err := r.Cookie(token.AccessTokenCookieName)
	if err != nil {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}

	userInfo, err := token.ValidateToDanniToken(accessTokenCookie.Value)
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

	ctx := context.WithValue(r.Context(), UserInfoContextKey, userInfo)

	ea.handler(w, r.WithContext(ctx))
}

func NewEnsureAuth(handlerToWrap http.HandlerFunc) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}
