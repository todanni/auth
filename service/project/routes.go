package project

import (
	"net/http"

	"github.com/todanni/token"
)

func (s *projectService) routes() {
	s.router.Handle("/api/projects/", token.NewAuthenticationMiddleware(s.CreateProjectHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/projects/", token.NewAuthenticationMiddleware(s.ListProjectsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/projects/{id}", token.NewAuthenticationMiddleware(s.DeleteProjectHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/projects/{id}", token.NewAuthenticationMiddleware(s.GetProjectHandler)).Methods(http.MethodGet)
}
