package project

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

func (s *projectService) routes() {
	s.router.Handle("/api/projects/", middleware.NewAuthenticationCheck(s.CreateProjectHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/projects/", middleware.NewAuthenticationCheck(s.ListProjectsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/projects/{id}", middleware.NewAuthenticationCheck(s.DeleteProjectHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/projects/{id}", middleware.NewAuthenticationCheck(s.GetProjectHandler)).Methods(http.MethodGet)
}
