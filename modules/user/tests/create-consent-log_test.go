package tests

import (
	"nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	"net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

var ConsentLog = dtos.CreateConsentBase{
	IpAddress:    "192.168.1.1",
	TermsVersion: "v1",
	UserAgent:    "Mozilla",
}

func TestCreateConsentLog(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/user/consent"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, ConsentLog, headers)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, resp["id_consent_log"])
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", resp["id_user"])
			assert.Equal(t, ConsentLog.TermsVersion, resp["terms_version"])
			assert.NotNil(t, resp["accepted_at"])
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			invalidHeader := &utils.TestHeadersAdmin

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, ConsentLog, invalidHeader)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create consent log", resp["message"])
		})
	})
}
