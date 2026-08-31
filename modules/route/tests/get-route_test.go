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

func TestGetRoute(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	driverHeaders := &utils.TestHeadersDriver
	nurseHeaders := &utils.TestHeadersNurse
	commonHeaders := &utils.TestHeaders

	const (
		idRoute     = "rot_3Idl55gnHyIrunoK64VSS5XX6Rg"
		idRouteFake = "rot_2veL1FPpuXxUaZcFaEC57Bfpd53"
		driverName  = "Carlos Motorista"
		idStopOne   = "rds_3Idl54SAt7SBYcutzdS6MYaKViF"
		idStopTwo   = "rds_3Idl53TDdo3E18YXp3pAacU89zi"
	)

	endpointOf := func(id string) string {
		return fmt.Sprintf("/internal/route/%s", id)
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Returns the route with driver name and ordered stops", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "GET", endpointOf(idRoute), nil, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, idRoute, resp["id_route"])
			assert.Equal(t, driverName, resp["driver_name"])

			stops, ok := resp["stops"].([]interface{})
			assert.True(t, ok)
			assert.Len(t, stops, 2)

			first := stops[0].(map[string]interface{})
			second := stops[1].(map[string]interface{})

			assert.Equal(t, idStopOne, first["id_route_donation_step"])
			assert.Equal(t, idStopTwo, second["id_route_donation_step"])
			assert.Equal(t, float64(0), first["stop_order"])
			assert.Equal(t, float64(1), second["stop_order"])

			address, ok := first["address"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, "Avenida Paulista", address["street"])
			assert.Equal(t, "01310200", address["zipcode"])
			assert.Equal(t, "Bela Vista", address["neighborhood"])
		})

		t.Run("Drivers and nurses can also see a route", func(t *testing.T) {
			for _, headers := range []*fluxgo.Headers{driverHeaders, nurseHeaders} {
				status, resp := fluxgo.RunTestRequest(app, "GET", endpointOf(idRoute), nil, headers)

				assert.Equal(t, http.StatusOK, status)
				assert.Equal(t, idRoute, resp["id_route"])
			}
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "GET", endpointOf(idRoute), nil, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to get route", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Route not found", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "GET", endpointOf(idRouteFake), nil, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route not found", resp["message"])
		})

		t.Run("Invalid route id", func(t *testing.T) {
			status, _ := fluxgo.RunTestRequest(app, "GET", endpointOf("invalid-id"), nil, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
		})
	})
}
