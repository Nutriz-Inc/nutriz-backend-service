package tests

import (
	"fmt"
	"net/http"
	"testing"

	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestGetDonation(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		id := "don_2veL1FPpuXxUaZcFaEC57BfpcKE"

		t.Run("Normal", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_donation"])
			assert.NotNil(t, body["steps"])
			assert.Len(t, body["steps"], 1)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Not found", func(t *testing.T) {
			id := "don_2veL1FPpuXxUaZcFaEC57Bfpd53"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation not found", body["message"])
		})

		t.Run("No permission", func(t *testing.T) {
			id := "don_2veL1FPpuXxUaZcFaEC57BfpcKF"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", body["message"])
		})
	})
}
