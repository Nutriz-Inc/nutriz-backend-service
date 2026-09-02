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

func TestListDonationSteps(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders
	nurseHeaders := &utils.TestHeadersNurse
	driverHeaders := &utils.TestHeadersDriver

	const (
		basePath   = "/internal/donation/step"
		endpoint   = basePath + "?page=1&page_size=25"
		idDonation = "don_2veL1FPpuXxUaZcFaEC57BfpcR2"
	)

	withQuery := func(params map[string]string) string {
		q := url.Values{}
		q.Set("page", "1")
		q.Set("page_size", "25")
		for k, v := range params {
			q.Set(k, v)
		}
		return fmt.Sprintf("%s?%s", basePath, q.Encode())
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Lists donation steps with address and pagination", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "GET", endpoint, nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, float64(1), resp["page"])
			assert.Equal(t, float64(25), resp["page_size"])
			assert.GreaterOrEqual(t, resp["total"], float64(1))

			data, ok := resp["data"].([]interface{})
			assert.True(t, ok)
			assert.NotEmpty(t, data)

			first := data[0].(map[string]interface{})
			assert.NotNil(t, first["id_donation_step"])
			assert.NotNil(t, first["status"])
			assert.Contains(t, first, "address")
		})

		t.Run("Filters by id_donation", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"id_donation": idDonation}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.Len(t, data, 3)
			for _, item := range data {
				step := item.(map[string]interface{})
				assert.Equal(t, idDonation, step["id_donation"])
				assert.Contains(t, step, "address")
			}
		})

		t.Run("Filters by name", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"id_donation": idDonation, "name": "Entregar kit de ordenha"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.Len(t, data, 1)
			assert.Equal(t, "Entregar kit de ordenha", data[0].(map[string]interface{})["name"])
		})

		t.Run("Filters by city and neighborhood", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{
					"id_donation":  idDonation,
					"city":         "Sao Paulo",
					"neighborhood": "Vila Mariana",
				}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.Len(t, data, 1)

			address := data[0].(map[string]interface{})["address"].(map[string]interface{})
			assert.Equal(t, "Vila Mariana", address["neighborhood"])
			assert.Equal(t, "Rua Domingos de Morais", address["street"])
		})

		t.Run("Filters by has_address = true", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"has_address": "true"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.NotEmpty(t, data)
			for _, item := range data {
				step := item.(map[string]interface{})
				assert.NotNil(t, step["address"], "every step must carry an address")
			}
		})

		t.Run("Filters by has_address = false", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"has_address": "false"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.NotEmpty(t, data)
			for _, item := range data {
				step := item.(map[string]interface{})
				assert.Nil(t, step["address"], "no step should carry an address")
			}
		})

		t.Run("Filters by available_for_route = false (already routed)", func(t *testing.T) {
			seedRouted := []string{
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd54",
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd55",
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd56",
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"available_for_route": "false", "page_size": "50"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})

			got := map[string]bool{}
			for _, item := range data {
				got[item.(map[string]interface{})["id_donation_step"].(string)] = true
			}

			for _, id := range seedRouted {
				assert.True(t, got[id], "routed seed step %s should be listed", id)
			}
			// a step with no route_donation_step at all must never appear here
			assert.False(t, got["dst_2veL1FPpuXxUaZcFaEC57BfpcAA"])
		})

		t.Run("Filters by available_for_route = true (free steps)", func(t *testing.T) {
			routed := map[string]bool{
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd54": true,
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd55": true,
				"dst_2veL1FPpuXxUaZcFaEC57Bfpd56": true,
			}

			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"available_for_route": "true"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			data := resp["data"].([]interface{})
			assert.NotEmpty(t, data)
			for _, item := range data {
				step := item.(map[string]interface{})
				assert.False(t, routed[step["id_donation_step"].(string)], "steps tied to an active route stop must be hidden")
			}
		})

		t.Run("Filters by status with no matches", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"GET",
				withQuery(map[string]string{"id_donation": idDonation, "status": "done"}),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Empty(t, resp["data"])
			assert.Equal(t, float64(0), resp["total"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Non-admin users cannot list donation steps", func(t *testing.T) {
			for _, headers := range []*fluxgo.Headers{commonHeaders, nurseHeaders, driverHeaders} {
				status, resp := fluxgo.RunTestRequest(app, "GET", endpoint, nil, headers)

				assert.Equal(t, http.StatusForbidden, status)
				assert.Equal(t, "User does not have permission to list donation steps", resp["message"])
				assert.Equal(t, "user.forbidden", resp["code"])
			}
		})
	})
}
