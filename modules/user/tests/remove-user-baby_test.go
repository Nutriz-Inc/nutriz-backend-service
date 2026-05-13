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

func TestRemoveUserBaby(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	headers := &utils.TestHeaders

	t.Run("success", func(t *testing.T) {
        userBabyId := "usb_01JTG8K8N4P2R6T9V1X3Y5Z7DEL"
        endpoint := fmt.Sprintf("/internal/user/baby/%s", userBabyId)
        status, resp := fluxgo.RunTestRequest(app, "DELETE", endpoint, nil, headers)
        assert.Equal(t, http.StatusOK, status)
        assert.Equal(t, true, resp["success"])
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Baby not found", func(t *testing.T) {
            route := "/internal/user/baby/usb_2veL1FPpuXxUaZcFaEC57BfpcZZ"
            status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, headers)
            assert.Equal(t, http.StatusNotFound, status)
            assert.Equal(t, "User baby not found", resp["message"])
        })

        t.Run("No permission", func(t *testing.T) {
            userBabyId := "usb_01JTG8K8N4P2R6T9V1X3Y5Z7MIP"
            route := fmt.Sprintf("/internal/user/baby/%s", userBabyId)
            status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, headers)
            assert.Equal(t, http.StatusForbidden, status)
            assert.Equal(t, "You don't have permission to access this resource", resp["message"])
        })
	})
}