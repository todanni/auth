package project

import "net/http"

func (s *projectService) routes() {
	s.router.HandleFunc("/api/projects/", s.CreateProjectHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/projects/", s.ListProjectsHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/projects/{id}", s.DeleteProjectHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/projects/{id}", s.GetProjectHandler).Methods(http.MethodGet)
}
