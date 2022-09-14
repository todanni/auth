package dashboard

import "net/http"

func (s *dashboardService) routes() {
	// GET an dashboard
	s.router.HandleFunc("/api/dashboard", s.CreateDashboardHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/dashboard/{id}", s.DeleteDashboardHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/dashboard/{id}", s.GetDashboardHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/dashboard/{id}/accept", s.AcceptDashboardInviteHandler).Methods(http.MethodPut)
	s.router.HandleFunc("/api/dashboard/{id}/reject", s.RejectDashboardInviteHandler).Methods(http.MethodPut)
}
