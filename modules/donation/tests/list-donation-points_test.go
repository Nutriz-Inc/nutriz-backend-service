package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"nutriz-backend-service/shared/module"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestListDonationPoint(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/public/donation/point?page=1&page_size=25"

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, nil)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			fluxgo.ConvertToMap(data[0])
			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_donation_point"])

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})
		t.Run("has_home filter", func(t *testing.T) {
			has_home := true
			route := fmt.Sprintf("%s&has_home=%t", endpoint, has_home)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, nil)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, has_home, item["has_home"])
		})
		t.Run("name filter", func(t *testing.T) {
			name := "Ponto Solidário Centro"
			route := fmt.Sprintf("%s&name=%s", endpoint, url.QueryEscape(name))

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, nil)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, name, item["name"])
		})
	})
}
