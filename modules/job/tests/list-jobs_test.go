package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestListJobs(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/job?page=1&page_size=25"
	adminHeaders := &utils.TestHeadersAdmin
	headers := &utils.TestHeadersNurse

	dateSetInTwoDays := time.Now().UTC().AddDate(0, 0, 2).Format(time.RFC3339)
	fluxgo.RunTestRequest(
		app,
		"PUT",
		"/internal/job/job_3EEMMlZS3VexkBHkGHPjq0Qt86V",
		dto.UpdateJobReq{
			IdUser:  utils.StringPtr("usr_2veL1FPpuXxUaZcFaEC57BfplNV"),
			DateSet: utils.StringPtr(dateSetInTwoDays),
		},
		adminHeaders,
	)

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_job"])
			assert.NotEmpty(t, first["id_donation"])

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfplNV", row["id_user"])
			}

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("No filters (admin)", func(t *testing.T) {
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
			dateSet := time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02")
			route := fmt.Sprintf("%s&date_set=%s", endpoint, dateSet)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Contains(t, item["date_set"], dateSet)
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfplNV", item["id_user"])
		})

		t.Run("id_user_nurse filter", func(t *testing.T) {
			idNurse := "usr_2veL1FPpuXxUaZcFaEC57BfplNV"
			route := fmt.Sprintf("%s&id_user_nurse=%s", endpoint, idNurse)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, idNurse, row["id_user"])
				assert.Equal(t, "Paula Fernandes", row["user_nurse_name"])
			}
		})

		t.Run("id_user_common filter", func(t *testing.T) {
			idDonor := "usr_2veL1FPpuXxUaZcFaEC57BfpcKE"
			route := fmt.Sprintf("%s&id_user_common=%s", endpoint, idDonor)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, "Maria Silva", row["user_common_name"])
			}
		})

		t.Run("user_nurse_name filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&user_nurse_name=Paula", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Contains(t, row["user_nurse_name"], "Paula")
			}
		})

		t.Run("user_common_name filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&user_common_name=Maria", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Contains(t, row["user_common_name"], "Maria")
			}
		})

		t.Run("status filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&status=done", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, "done", row["status"])
			}
		})
	})
}
