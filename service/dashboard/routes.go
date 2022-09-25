package dashboard

import (
	"net/http"

	"github.com/todanni/token"
)

const (
	ApiPath                       = "/api/dashboards"
	ListAndCreateDashboardHandler = "/api/dashboards/"
)

func (s *dashboardService) routes() {
	r := s.router.PathPrefix(ApiPath).Subrouter()
	r.Use(token.AuthMiddleware)

	r.HandleFunc("/", s.CreateDashboardHandler).Methods(http.MethodPost)
	r.HandleFunc("/", s.ListDashboardsHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", s.DeleteDashboardHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", s.GetDashboardHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}/accept", s.AcceptDashboardInviteHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}/reject", s.RejectDashboardInviteHandler).Methods(http.MethodPut)
}
