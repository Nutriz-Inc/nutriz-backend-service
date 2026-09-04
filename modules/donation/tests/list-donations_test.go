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

func TestListDonation(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/donation?page=1&page_size=25"
	headers := &utils.TestHeaders

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters (common user)", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_donation"])

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)

				assert.Equal(
					t,
					"usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
					row["created_by"],
				)
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
			assert.NotEmpty(t, first["id_donation"])

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				createdBy := row["created_by"]
				assert.Contains(
					t,
					[]string{
						"usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
						"usr_2veL1FPpuXxUaZcFaEC57BfpcKL",
						"usr_2veL1FPpuXxUaZcFaEC57BfpcKF",
						"usr_2veL1FPpuXxUaZcFaEC57BfpcBE",
					},
					createdBy,
				)
			}

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("is_active filter", func(t *testing.T) {
			is_active := true
			route := fmt.Sprintf("%s&is_active=%t", endpoint, is_active)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, is_active, item["is_active"])
		})
		t.Run("user_document filter as a admin", func(t *testing.T) {
			user_document := "47046117012"
			route := fmt.Sprintf("%s&user_document=%s", endpoint, user_document)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, &utils.TestHeadersAdmin)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpcKE", item["created_by"])
		})

		t.Run("id_user_common filter as admin", func(t *testing.T) {
			idUserCommon := "usr_2veL1FPpuXxUaZcFaEC57BfpcKE"
			route := fmt.Sprintf("%s&id_user_common=%s", endpoint, idUserCommon)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, &utils.TestHeadersAdmin)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, idUserCommon, row["created_by"])
			}
		})

		t.Run("id_user_common filter is ignored for common users", func(t *testing.T) {
			otherUserId := "usr_2veL1FPpuXxUaZcFaEC57BfpcKL"
			route := fmt.Sprintf("%s&id_user_common=%s", endpoint, otherUserId)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(
					t,
					"usr_2veL1FPpuXxUaZcFaEC57BfpcKE",
					row["created_by"],
				)
			}
		})
	})
}
