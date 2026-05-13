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
	headers := &utils.TestHeaders

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
				assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", row["id_user"])
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
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", item["id_user"])
		})
	})
}
