package token

import (
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/todanni/auth/models"
)

type Permissions struct {
	Dashboards []string `json:"dashboards"`
	Projects   []string `json:"projects"`
}

func SetDashboardsPermissions(dashboards []models.Dashboard, token jwt.Token) {
	userDashboardIDs := make([]uint, 0)

	for _, dashboard := range dashboards {
		userDashboardIDs = append(userDashboardIDs, dashboard.ID)
	}

	token.Set("dashboards", userDashboardIDs)
}

func SetProjectPermissions(projects []models.Project, token jwt.Token) {
	userDashboardIDs := make([]uint, 0)

	for _, project := range projects {
		userDashboardIDs = append(userDashboardIDs, project.ID)
	}

	token.Set("projects", userDashboardIDs)
}
