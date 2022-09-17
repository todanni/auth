package auth

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

const (
	LoginHandler        = "/auth/login"
	CallbackHandler     = "/auth/callback"
	PublicKeyHandler    = "/auth/public-key"
	RefreshTokenHandler = "/auth/refresh"
	UserInfoHandler     = "/auth/user-info"
)

func (s *authService) routes() {
	// GET an auth
	s.router.HandleFunc(LoginHandler, s.LoginHandler)
	s.router.HandleFunc(CallbackHandler, s.CallbackHandler)
	s.router.HandleFunc(RefreshTokenHandler, s.RefreshTokenHandler)
	s.router.HandleFunc(PublicKeyHandler, s.ServePublicKey).Methods(http.MethodGet)
	
	// only UserInfoHandler requires auth
	s.router.Handle(UserInfoHandler, middleware.NewEnsureAuth(s.UserInfoHandler)).Methods(http.MethodGet)
	
}
