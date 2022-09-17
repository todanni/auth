package dashboard

import (
	"net/http"

	"github.com/todanni/auth/middleware"
)

func (s *dashboardService) routes() {
	s.router.Handle("/api/dashboards", middleware.NewEnsureAuth(s.CreateDashboardHandler)).Methods(http.MethodPost)
	s.router.Handle("/api/dashboards", middleware.NewEnsureAuth(s.ListDashboardsHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}", middleware.NewEnsureAuth(s.DeleteDashboardHandler)).Methods(http.MethodDelete)
	s.router.Handle("/api/dashboards/{id}", middleware.NewEnsureAuth(s.GetDashboardHandler)).Methods(http.MethodGet)
	s.router.Handle("/api/dashboards/{id}/accept", middleware.NewEnsureAuth(s.AcceptDashboardInviteHandler)).Methods(http.MethodPut)
	s.router.Handle("/api/dashboards/{id}/reject", middleware.NewEnsureAuth(s.RejectDashboardInviteHandler)).Methods(http.MethodPut)
}
