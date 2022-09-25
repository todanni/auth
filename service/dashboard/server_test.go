package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/todanni/token"
	"gorm.io/gorm"

	"github.com/todanni/auth/models"
	"github.com/todanni/auth/storage/mocks"
)

func Test_dashboardService_ListDashboardsHandler(t *testing.T) {
	r := mux.NewRouter()

	dashboardStorageMock := mocks.NewDashboardStorage(t)
	dashboardStorageMock.On("List", mock.Anything).Return([]models.Dashboard{{
		Owner:  1,
		Status: models.PendingStatus,
		Members: []models.User{{
			Model: gorm.Model{
				ID: 1,
			},
			Email:      "test1@mail.com",
			ProfilePic: "profile-pic1.jpeg",
		},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Email:      "test2@mail.com",
				ProfilePic: "profile-pic2.jpeg",
			}},
	}}, nil)

	userStorageMock := mocks.NewUserStorage(t)

	s := NewDashboardService(dashboardStorageMock, userStorageMock, r)

	w := httptest.NewRecorder()
	r.Handle(ListAndCreateDashboardHandler, token.NewAuthenticationMiddleware(s.ListDashboardsHandler)).Methods(http.MethodGet)

	request := httptest.NewRequest(http.MethodGet, ListAndCreateDashboardHandler, nil)

	accessToken := token.NewAccessToken()
	accessToken.SetUserInfo(models.UserInfo{
		UserID: 1,
	})

	ctx := context.WithValue(context.Background(), token.AccessTokenContextKey, accessToken)
	request.WithContext(ctx)

	r.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)
}
