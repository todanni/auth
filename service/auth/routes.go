package auth

const (
	AuthoriseHandler = "/authorize"
	TokenHandler     = "/token"
)

func (s *authService) routes() {
	// GET an auth
	s.router.HandleFunc(TokenHandler, s.TokenHandler)
	s.router.HandleFunc(AuthoriseHandler, s.AuthoriseHandler)
}
