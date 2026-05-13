package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAddress(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		addressID := "adr_01JTX0H1V8N5Q3W7E2R4T6Y8ZXE"
		endpoint := fmt.Sprintf("/internal/user/address/%s", addressID)

		status, resp := fluxgo.RunTestRequest(app, "DELETE", endpoint, nil, headers)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, true, resp["success"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Address not found", func(t *testing.T) {
			route := "/internal/user/address/adr_01JTX0H1V8N5Q3W7E2R4T6Y8ZZZ"

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Address not found", resp["message"])
		})

		t.Run("You don't have permission to access this resource", func(t *testing.T) {
			route := "/internal/user/address/adr_01JTX0H1V8N5Q3W7E2R4T6Y8MBL"

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", resp["message"])
		})
	})
}
