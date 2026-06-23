package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetStatusDoesNotExposeCustomLogo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	GetStatus(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var body map[string]any
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &body))
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	require.NotContains(t, data, "logo")
}

func TestGetOptionsDoesNotExposeRemovedLogoOption(t *testing.T) {
	gin.SetMode(gin.TestMode)

	common.OptionMapRWMutex.Lock()
	originalOptionMap := common.OptionMap
	common.OptionMap = map[string]string{
		"Logo":       "https://example.com/custom-logo.png",
		"SystemName": "New API",
	}
	common.OptionMapRWMutex.Unlock()
	t.Cleanup(func() {
		common.OptionMapRWMutex.Lock()
		common.OptionMap = originalOptionMap
		common.OptionMapRWMutex.Unlock()
	})

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	GetOptions(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var body struct {
		Data []model.Option `json:"data"`
	}
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &body))
	for _, option := range body.Data {
		require.NotEqual(t, "Logo", option.Key)
	}
}

func TestUpdateOptionRejectsRemovedLogoOption(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload, err := common.Marshal(OptionUpdateRequest{
		Key:   "Logo",
		Value: "https://example.com/custom-logo.png",
	})
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/option/", bytes.NewReader(payload))

	UpdateOption(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var body struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &body))
	require.False(t, body.Success)
	require.Contains(t, body.Message, "Logo")
}
