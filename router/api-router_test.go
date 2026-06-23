package router

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLegacySynchronousLogCleanupRouteRemoved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	SetApiRouter(engine)

	hasLegacyDeleteRoute := false
	hasAsyncCleanupRoute := false
	for _, route := range engine.Routes() {
		if route.Method == http.MethodDelete && route.Path == "/api/log/" {
			hasLegacyDeleteRoute = true
		}
		if route.Method == http.MethodPost && route.Path == "/api/system-task/log-cleanup" {
			hasAsyncCleanupRoute = true
		}
	}

	require.False(t, hasLegacyDeleteRoute)
	require.True(t, hasAsyncCleanupRoute)
}
