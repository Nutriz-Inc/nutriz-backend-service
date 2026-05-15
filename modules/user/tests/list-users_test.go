package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/user?page=1&page_size=25"
	headers := &utils.TestHeadersAdmin

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_user"])

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("name filter", func(t *testing.T) {
			name := "Maria Silva"
			route := fmt.Sprintf("%s&name=%s", endpoint, url.QueryEscape(name))

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, name, item["name"])
		})

		t.Run("type filter", func(t *testing.T) {
			userType := "common"
			route := fmt.Sprintf("%s&type=%s", endpoint, userType)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, userType, item["type"])
		})
		t.Run("internal_identifier filter", func(t *testing.T) {
			identifier := "234567898765435"
			route := fmt.Sprintf("%s&internal_identifier=%s", endpoint, identifier)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, headers)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			item := fluxgo.ConvertToMap(data[0])

			assert.Equal(t, identifier, item["internal_identifier"])
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("No permission", func(t *testing.T) {
			commonHeaders := &utils.TestHeaders
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to list users", body["message"])
		})
	})
}
