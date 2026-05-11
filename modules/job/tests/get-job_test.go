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

func TestGetJob(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/job"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		id := "job_2veL1FPpuXxUaZcFaEC57Bfpc01"

		t.Run("Normal", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, id, body["id_job"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Not found", func(t *testing.T) {
			id := "job_not_exists"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Job not found", body["message"])
		})

		t.Run("No permission", func(t *testing.T) {
			id := "job_2veL1FPpuXxUaZcFaEC57Bfpc02"
			route := fmt.Sprintf("%s/%s", endpoint, id)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "You don't have permission to access this resource", body["message"])
		})
	})
}
