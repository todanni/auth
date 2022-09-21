package project

import "net/http"

func (s *projectService) routes() {
	s.router.HandleFunc("/api/Projects/", s.CreateProjectHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/Projects/", s.ListProjectsHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/Projects/{id}", s.DeleteProjectHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/Projects/{id}", s.GetProjectHandler).Methods(http.MethodGet)
}
