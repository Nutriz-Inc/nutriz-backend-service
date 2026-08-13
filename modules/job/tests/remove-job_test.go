package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestRemoveJob(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/v1/internal/job"
	adminHeaders := &utils.TestHeadersAdmin
	nurseHeaders := &utils.TestHeadersNurse
	commonHeaders := &utils.TestHeaders

	_, jobResp := fluxgo.RunTestRequest(
		app,
		"POST",
		endpoint,
		map[string]any{
			"id_user":     "usr_2veL1FPpuXxUaZcFaEC57BfplNV",
			"id_step":     "dst_2veL1FPpuXxUaZcFaEC57BfslJG",
			"name":        "Job para remoção",
			"description": "Job criado para teste de remoção",
			"date_set":    time.Now().UTC().AddDate(0, 0, 5).Format(time.RFC3339),
		},
		adminHeaders,
	)
	jobId, _ := jobResp["id_job"].(string)

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission to remove job", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, jobId)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only admins can remove jobs", resp["message"])
			assert.Equal(t, "job.forbidden", resp["code"])
		})

		t.Run("Nurse does not have permission to remove job", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, jobId)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, nurseHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "Only admins can remove jobs", resp["message"])
			assert.Equal(t, "job.forbidden", resp["code"])
		})

		t.Run("Job not found", func(t *testing.T) {
			route := fmt.Sprintf("%s/job_2veL1FPpuXxUaZcFaEC57BfpcZZ", endpoint)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Job not found", resp["message"])
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("Admin removes job", func(t *testing.T) {
			route := fmt.Sprintf("%s/%s", endpoint, jobId)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", route, nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, resp["success"])
		})
	})
}
