package dashboard

import (
	"net/http"

	"github.com/todanni/token"
)

const (
	ListAndCreateDashboardHandler = "/api/dashboards/"
	GetAndDeleteDashboardHandler  = "/api/dashboards/{id}"
	AcceptDashboardHandler        = "/api/dashboards/{id}/accept"
	RejectDashboardHandler        = "/api/dashboards/{id}/reject"
)

func (s *dashboardService) routes() {
	s.router.Handle(ListAndCreateDashboardHandler, token.NewAuthenticationMiddleware(s.CreateDashboardHandler)).Methods(http.MethodPost)
	s.router.Handle(ListAndCreateDashboardHandler, token.NewAuthenticationMiddleware(s.ListDashboardsHandler)).Methods(http.MethodGet)
	s.router.Handle(GetAndDeleteDashboardHandler, token.NewAuthenticationMiddleware(s.DeleteDashboardHandler)).Methods(http.MethodDelete)
	s.router.Handle(GetAndDeleteDashboardHandler, token.NewAuthenticationMiddleware(s.GetDashboardHandler)).Methods(http.MethodGet)
	s.router.Handle(AcceptDashboardHandler, token.NewAuthenticationMiddleware(s.AcceptDashboardInviteHandler)).Methods(http.MethodPut)
	s.router.Handle(RejectDashboardHandler, token.NewAuthenticationMiddleware(s.RejectDashboardInviteHandler)).Methods(http.MethodPut)
}
