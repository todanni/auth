package project

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

func (s *projectService) routes() {
	s.router.Handle("/api/projects/", middleware.NewAuthenticationMiddleware(s.CreateProjectHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/projects/", middleware.NewAuthenticationMiddleware(s.ListProjectsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/projects/{id}", middleware.NewAuthenticationMiddleware(s.DeleteProjectHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/projects/{id}", middleware.NewAuthenticationMiddleware(s.GetProjectHandler)).Methods(http.MethodGet)
}
