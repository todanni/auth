package dashboard

import (
	"net/http"

	"github.com/todanni/token"
)

func (s *dashboardService) routes() {
	s.router.Handle("/api/dashboards/", token.NewAuthenticationMiddleware(s.CreateDashboardHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/dashboards/", token.NewAuthenticationMiddleware(s.ListDashboardsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}", token.NewAuthenticationMiddleware(s.DeleteDashboardHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/dashboards/{id}", token.NewAuthenticationMiddleware(s.GetDashboardHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}/accept", token.NewAuthenticationMiddleware(s.AcceptDashboardInviteHandler)).Methods(http.MethodPut)
	s.router.Handle("/api/dashboards/{id}/reject", token.NewAuthenticationMiddleware(s.RejectDashboardInviteHandler)).Methods(http.MethodPut)
}
