package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/shared/module"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestListDonation(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/public/donation?page=1&page_size=25"

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, nil)
			fmt.Println("body", body)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			fluxgo.ConvertToMap(data[0])
			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_donation"])

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("is_active filter", func(t *testing.T) {
			is_active := true
			route := fmt.Sprintf("%s&is_active=%t", endpoint, is_active)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, nil)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, is_active, item["is_active"])
		})
	})
}