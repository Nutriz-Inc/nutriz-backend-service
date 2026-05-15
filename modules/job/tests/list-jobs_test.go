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

func TestListJobs(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/job?page=1&page_size=25"
	headers := &utils.TestHeadersNurse

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_job"])

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfplNV", row["id_user"])
			}

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("No filters (admin)", func(t *testing.T) {
			adminHeaders := &utils.TestHeadersAdmin
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_job"])

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				idUser := row["id_user"]
				assert.Contains(
					t,
					[]string{
						"usr_2veL1FPpuXxUaZcFaEC57BfplNV",
						"usr_2veL1FPpuXxUaZcFaEC57BfpcKL",
						"usr_2veL1FPpuXxUaZcFaEC57BfpcNF",
					},
					idUser,
				)
			}

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("date_set filter", func(t *testing.T) {
			dateSet := utils.GetTodayFormattedDate(2)
			route := fmt.Sprintf("%s&date_set=%s", endpoint, dateSet)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Contains(t, item["date_set"], dateSet)
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfplNV", item["id_user"])
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			commonHeaders := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to list jobs", body["message"])
		})
	})
}
