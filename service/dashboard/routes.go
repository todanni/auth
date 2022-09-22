package dashboard

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

func (s *dashboardService) routes() {
	s.router.Handle("/api/dashboards/", middleware.NewAuthenticationCheck(s.CreateDashboardHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/dashboards/", middleware.NewAuthenticationCheck(s.ListDashboardsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}", middleware.NewAuthenticationCheck(s.DeleteDashboardHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/dashboards/{id}", middleware.NewAuthenticationCheck(s.GetDashboardHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}/accept", middleware.NewAuthenticationCheck(s.AcceptDashboardInviteHandler)).Methods(http.MethodPut)
	s.router.Handle("/api/dashboards/{id}/reject", middleware.NewAuthenticationCheck(s.RejectDashboardInviteHandler)).Methods(http.MethodPut)
}
