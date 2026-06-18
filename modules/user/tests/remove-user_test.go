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

func TestRemoveUser(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	nurseHeaders := &utils.TestHeadersNurse

	t.Run("Success", func(t *testing.T) {
		t.Run("Remove common user by nurse", func(t *testing.T) {
			userID := "usr_2veL1FPpuXxUaZcFaEC57BfpcKF"
			endpoint := fmt.Sprintf("/internal/user/%s", userID)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", endpoint, nil, nurseHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, resp["success"])
		})

		t.Run("Remove nurse user by admin", func(t *testing.T) {
			userID := "usr_2veL1FPpuXxUaZcFaEC57BfpcKH"
			endpoint := fmt.Sprintf("/internal/user/%s", userID)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", endpoint, nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, resp["success"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User not found", func(t *testing.T) {
			route := "/internal/user/usr_2veL1FPpuXxUaZcAbEC57BfpcMM"

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "User not found", resp["message"])
		})

		t.Run("No permission", func(t *testing.T) {
			userID := "usr_2veL1FPpuXxUaZcFaEC57BfpcKG"
			endpoint := fmt.Sprintf("/internal/user/%s", userID)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", endpoint, nil, adminHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", resp["message"])
		})
	})
}
