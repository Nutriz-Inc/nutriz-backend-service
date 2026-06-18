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

func TestListDonationStepTimeline(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation/step"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("List timeline with data", func(t *testing.T) {
			id := "dst_2veL1FPpuXxUaZcFaEC57BfpcAA"
			route := fmt.Sprintf("%s/%s/timeline", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_donation_step_timeline"])
			assert.Equal(t, id, first["id_donation_step"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Donation step not found", func(t *testing.T) {
			id := "dst_2veL1FPpuXxUaZcFaEC57Bfpd99"
			route := fmt.Sprintf("%s/%s/timeline", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation step not found", body["message"])
		})

		t.Run("User does not have permission to list donation step timeline", func(t *testing.T) {
			id := "dst_2veL1FPpuXxUaZcFaEC57BfpcBB"
			route := fmt.Sprintf("%s/%s/timeline", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to list donation step timeline", body["message"])
		})
	})
}
