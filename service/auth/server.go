package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openshift/osin"
)

type AuthService interface {
	TokenHandler(w http.ResponseWriter, r *http.Request)
	AuthoriseHandler(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	router *mux.Router
	server *osin.Server
}

func NewAuthService(osinServer *osin.Server, router *mux.Router) AuthService {
	server := &authService{
		router: router,
		server: osinServer,
	}
	server.routes()
	return server
}

func (s *authService) TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp := s.server.NewResponse()
	defer resp.Close()

	if ar := s.server.HandleAccessRequest(resp, r); ar != nil {
		ar.Authorized = true
		s.server.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func (s *authService) AuthoriseHandler(w http.ResponseWriter, r *http.Request) {
	resp := s.server.NewResponse()
	defer resp.Close()

	if ar := s.server.HandleAuthorizeRequest(resp, r); ar != nil {

		// HANDLE LOGIN PAGE HERE

		ar.Authorized = true
		s.server.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}
