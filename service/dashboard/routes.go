package dashboard

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

func (s *dashboardService) routes() {
	s.router.Handle("/api/dashboards/", middleware.NewAuthenticationMiddleware(s.CreateDashboardHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/dashboards/", middleware.NewAuthenticationMiddleware(s.ListDashboardsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}", middleware.NewAuthenticationMiddleware(s.DeleteDashboardHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/dashboards/{id}", middleware.NewAuthenticationMiddleware(s.GetDashboardHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}/accept", middleware.NewAuthenticationMiddleware(s.AcceptDashboardInviteHandler)).Methods(http.MethodPut)
	s.router.Handle("/api/dashboards/{id}/reject", middleware.NewAuthenticationMiddleware(s.RejectDashboardInviteHandler)).Methods(http.MethodPut)
}
