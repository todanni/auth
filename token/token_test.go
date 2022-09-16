package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateToDanniToken(t *testing.T) {
	tokenString := ""

	result, err := ValidateToDanniToken(tokenString)
	require.NoError(t, err)
	require.Equal(t, "danni@todanni.com", result.Email)
}
