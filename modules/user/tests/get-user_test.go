package tests

import (
	"fmt"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	"net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/user"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		id := "usr_2veL1FPpuXxUaZcFaEC57BfpcKE"
		internalIdentifier := "234567898765435"

		t.Run("Normal", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])
			assert.Equal(t, internalIdentifier, body["internal_identifier"])
			assert.Nil(t, body["password"])
			assert.Nil(t, body["donations_completed"])
			assert.Nil(t, body["current_donation"])
			assert.Nil(t, body["addresses"])
			assert.Nil(t, body["babies"])
		})
		t.Run("With address", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_address=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])

			addresses, ok := body["addresses"].([]interface{})
			assert.True(t, ok)

			assert.NotNil(t, addresses)
			assert.Equal(t, len(addresses), 1)
		})
		t.Run("With baby", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_baby=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])

			babies, ok := body["babies"].([]interface{})
			assert.True(t, ok)

			assert.NotNil(t, babies)
			assert.GreaterOrEqual(t, len(babies), 2)
		})
		t.Run("With donations completed", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_donations_completed=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])

			donationsCompleted, ok := body["donations_completed"].(float64)
			assert.True(t, ok)
			assert.Equal(t, float64(1), donationsCompleted)
		})
		t.Run("With current donation", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_current_donation=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])

			currentDonation, ok := body["current_donation"].(map[string]interface{})
			assert.True(t, ok)
			assert.NotNil(t, currentDonation)
			assert.Equal(t, "don_2veL1FPpuXxUaZcFaEC57BfpcC3", currentDonation["id_donation"])
			assert.Equal(t, true, currentDonation["is_active"])
		})
		t.Run("Donation filters are ignored for non-common users", func(t *testing.T) {
			adminId := "usr_2veL1FPpuXxUaZcFaEC57BfpxWS"
			route := fmt.Sprintf("%s/%s?show_address=true&show_baby=true&show_donations_completed=true&show_current_donation=true", endpoint, adminId)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, &utils.TestHeadersAdmin)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, adminId, body["id_user"])
			assert.Nil(t, body["addresses"])
			assert.Nil(t, body["babies"])
			assert.Nil(t, body["donations_completed"])
			assert.Nil(t, body["current_donation"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Not found", func(t *testing.T) {
			id := "usr_2veL1FPpuXxUaZcAbEC57BfpcMM"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "User not found", body["message"])
		})
		t.Run("No permission", func(t *testing.T) {
			id := "usr_2veL1FPpuXxUaZcFaEC57BfpcKL"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", body["message"])
		})
	})
}
