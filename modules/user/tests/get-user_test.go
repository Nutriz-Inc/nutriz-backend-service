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
		})
		t.Run("With address", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_address=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])
			assert.NotNil(t, body["addresses"])
			assert.Len(t, body["addresses"], 1)
		})
		t.Run("With baby", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s?show_baby=true", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_user"])
			assert.NotNil(t, body["babies"])
			assert.Len(t, body["babies"], 2)
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
