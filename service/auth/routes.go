package auth

import "net/http"

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
	s.router.HandleFunc(UserInfoHandler, s.UserInfoHandler).Methods(http.MethodGet)

}
