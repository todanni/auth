package token

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/todanni/auth/models"
	"github.com/todanni/auth/test"
	"gorm.io/gorm"
)

func Test_ValidateToDanniToken(t *testing.T) {
	key := test.ServePublicKey()
	user := models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:      "danni@todanni.com",
		ProfilePic: "",
	}

	dashboards := make([]models.Dashboard, 0)
	projects := make([]models.Project, 0)

	tokenString, err := IssueToDanniToken(user, key, dashboards, projects)
	require.NoError(t, err)

	result, err := ValidateToDanniToken(tokenString)
	require.NoError(t, err)
	require.Equal(t, "danni@todanni.com", result.Email)
}
