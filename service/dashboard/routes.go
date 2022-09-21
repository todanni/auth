package dashboard

import "net/http"

func (s *dashboardService) routes() {
	s.router.HandleFunc("/api/dashboards/", s.CreateDashboardHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/dashboards/", s.ListDashboardsHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/dashboards/{id}", s.DeleteDashboardHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/dashboards/{id}", s.GetDashboardHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/dashboards/{id}/accept", s.AcceptDashboardInviteHandler).Methods(http.MethodPut)
	s.router.HandleFunc("/api/dashboards/{id}/reject", s.RejectDashboardInviteHandler).Methods(http.MethodPut)
}
