package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateToDanniToken(t *testing.T) {
	tokenString := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhNDM2MmM0YmY4MmZkZTJiMDU5IiwidHlwIjoiSldUIn0.eyJlbWFpbCI6ImRhbm5pQHRvZGFubmkuY29tIiwiaWF0IjoxNjYzMjU3MTU4LCJpc3MiOiJ0b2Rhbm5pLmNvbSIsInByb2ZpbGVQaWMiOiIiLCJ1c2VySUQiOjN9.FUq5KRCMmrPTbiDMZXhoSMngJVsJMi99yi2nSNz4QRKriE8agfFkw8_rEwBiIX3_Trc4BPMNN91h2wZW1Ni32GCjAiLgeNAKOF4AM6Px9SaBpkvGu3YWv7pQsA9CaQo-6RzrxM3e7a_-tMMJXBXhnyfy79h3JLiJFjkA7owEvUBOXfEM1nIsXd0-BXClpHbCfcgCVpx0Sx5pZMBVI-cVS60o1pp7_YXIIcvedsnHrsnNrgiU0xZSl4lvj8YVKr8qfxPCP2wa-jQVNN2-1Km6HtTPzyfDEAtfuGCkaBegrOf7zFwFlT-y8TgO1RXeWaw7kbHczZdmszTjEWXQcoAtqg"

	result, err := ValidateToDanniToken(tokenString)
	require.NoError(t, err)
	require.Equal(t, "danni@todanni.com", result.Email)
}
