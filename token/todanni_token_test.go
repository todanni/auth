package token

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/todanni/auth/models"
	"github.com/todanni/auth/test"
)

func TestToDanniToken_SignAndValidate(t *testing.T) {
	privateKey := test.ServePublicKey()

	token := NewAccessToken()
	signedToken, err := token.SignedToken(privateKey)
	require.NoError(t, err)

	err = token.Validate(signedToken)
	require.NoError(t, err)
}

func TestToDanniToken_GetAndSetUserInfo(t *testing.T) {
	testEmail, testProfilePic := "test@mail.com", "profile-pic.jpg"

	token := NewAccessToken()
	token.SetUserInfo(models.UserInfo{
		Email:      testEmail,
		ProfilePic: testProfilePic,
		UserID:     1,
	})

	userInfo, err := token.GetUserInfo()
	require.NoError(t, err)
	require.Equal(t, testEmail, userInfo.Email)
	require.Equal(t, testProfilePic, userInfo.ProfilePic)
	require.Equal(t, 1, userInfo.UserID)
}

func TestToDanniToken_ProjectPermission(t *testing.T) {
	token := NewAccessToken()
	token.SetProjectsPermissions([]models.Project{
		{
			Model: gorm.Model{
				ID: 1,
			},
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
		},
	})

	require.True(t, token.HasProjectPermission(1))
	require.True(t, token.HasProjectPermission(2))
	require.False(t, token.HasProjectPermission(3))
}

func TestToDanniToken_DashboardPermission(t *testing.T) {
	token := NewAccessToken()
	token.SetDashboardPermissions([]models.Dashboard{
		{
			Model: gorm.Model{
				ID: 1,
			},
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
		},
	})

	require.True(t, token.HasDashboardPermission(1))
	require.True(t, token.HasDashboardPermission(2))
	require.False(t, token.HasDashboardPermission(3))
}
