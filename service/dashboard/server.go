package dashboard

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/todanni/token"

	"github.com/todanni/auth/models"
	"github.com/todanni/auth/storage"
)

type DashboardService interface {
	CreateDashboardHandler(w http.ResponseWriter, r *http.Request)
	GetDashboardHandler(w http.ResponseWriter, r *http.Request)
	ListDashboardsHandler(w http.ResponseWriter, r *http.Request)
	AcceptDashboardInviteHandler(w http.ResponseWriter, r *http.Request)
	RejectDashboardInviteHandler(w http.ResponseWriter, r *http.Request)
	DeleteDashboardHandler(w http.ResponseWriter, r *http.Request)
}

func NewDashboardService(dashboardStorage storage.DashboardStorage, userStorage storage.UserStorage, router *mux.Router) DashboardService {
	server := &dashboardService{
		dashboardStorage: dashboardStorage,
		userStorage:      userStorage,
		router:           router,
	}
	if router != nil {
		server.routes()
	}
	return server
}

type dashboardService struct {
	router           *mux.Router
	dashboardStorage storage.DashboardStorage
	userStorage      storage.UserStorage
}

func (s *dashboardService) ListDashboardsHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userInfo, err := accessToken.GetUserInfo()
	if err != nil {
		http.Error(w, "token didn't contain user info", http.StatusBadRequest)
	}

	dashboards, err := s.dashboardStorage.List(userInfo.UserID)
	if err != nil {
		log.Error(err)
		return
	}

	response := make([]models.ListDashboardsResponse, 0)
	for _, dashboard := range dashboards {
		response = append(response, models.ListDashboardsResponse{
			Owner:  dashboard.Owner,
			Status: dashboard.Status,
			Members: []models.Member{
				{
					ID:         dashboard.Members[0].ID,
					Email:      dashboard.Members[0].Email,
					ProfilePic: dashboard.Members[0].ProfilePic,
				},
				{
					ID:         dashboard.Members[1].ID,
					Email:      dashboard.Members[1].Email,
					ProfilePic: dashboard.Members[1].ProfilePic,
				},
			},
		})
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't marshal body", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *dashboardService) CreateDashboardHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userInfo, err := accessToken.GetUserInfo()
	if err != nil {
		http.Error(w, "token didn't contain user info", http.StatusBadRequest)
	}

	// Parse the body of the request to get the email of the invited user
	var requestBody models.DashboardCreateRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userStorage.GetUser(requestBody.Email)
	if err != nil {
		http.Error(w, "couldn't find user", http.StatusBadRequest)
	}

	dashboard, err := s.dashboardStorage.Create(userInfo.UserID, user.ID)
	if err != nil {
		http.Error(w, "couldn't persist the new dashboard", http.StatusInternalServerError)
	}

	response := models.DashboardCreateResponse{
		Owner:   dashboard.Owner,
		Status:  dashboard.Status,
		Members: []uint{dashboard.Members[0].ID, dashboard.Members[1].ID},
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't marshal body", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *dashboardService) DeleteDashboardHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dashboardID := params["id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userInfo, err := accessToken.GetUserInfo()
	if err != nil {
		http.Error(w, "token didn't contain user info", http.StatusBadRequest)
	}

	dashboard, err := s.dashboardStorage.GetById(dashboardID)
	if err != nil {
		http.Error(w, "couldn't find dashboard", http.StatusNotFound)
	}

	if dashboard.Owner != (userInfo.UserID) {
		http.Error(w, "only owners can delete dashboards", http.StatusForbidden)
	}

	err = s.dashboardStorage.Delete(dashboardID)
	if err != nil {
		http.Error(w, "couldn't delete dashboard from database", http.StatusInternalServerError)
	}
}

func (s *dashboardService) GetDashboardHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dashboardID := params["id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	dashIdUint, err := strconv.ParseUint(dashboardID, 10, 32)
	if err != nil {
		http.Error(w, "invalid dashboard ID", http.StatusBadRequest)
	}

	if !accessToken.HasDashboardPermission(uint(dashIdUint)) {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
	}

	dashboard, err := s.dashboardStorage.GetById(dashboardID)
	if err != nil {
		http.Error(w, "couldn't find dashboard", http.StatusNotFound)
	}

	responseBody, err := json.Marshal(dashboard)
	if err != nil {
		http.Error(w, "couldn't marshal body", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *dashboardService) AcceptDashboardInviteHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *dashboardService) RejectDashboardInviteHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
