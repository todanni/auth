package auth

const (
	LoginHandler        = "/login"
	CallbackHandler     = "/auth/callback"
	PublicKeyHandler    = "/public-key"
	RefreshTokenHandler = "/refresh"
)

func (s *authService) routes() {
	// GET an auth
	s.router.HandleFunc(LoginHandler, s.LoginHandler)
	s.router.HandleFunc(CallbackHandler, s.CallbackHandler)
	s.router.HandleFunc(RefreshTokenHandler, s.RefreshTokenHandler)
}
