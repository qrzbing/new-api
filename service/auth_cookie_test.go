package service

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshCookieUsesAppBasePath(t *testing.T) {
	previousBasePath := common.AppBasePath
	common.AppBasePath = "/tools/new-api"
	t.Cleanup(func() { common.AppBasePath = previousBasePath })

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	WriteRefreshCookie(context, "invalid-test-token")

	cookies := response.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "/tools/new-api/api/user/auth", cookies[0].Path)
}
