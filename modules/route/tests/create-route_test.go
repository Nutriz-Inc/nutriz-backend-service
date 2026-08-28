package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoute(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/route"
	adminHeaders := &utils.TestHeadersAdmin
	commonHeaders := &utils.TestHeaders

	const (
		idDriver          = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		idCommonUser      = "usr_2veL1FPpuXxUaZcFaEC57BfpcKL"
		idStepOne         = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idStepTwo         = "dst_2veL1FPpuXxUaZcFaEC57BfpcR2"
		idStepThree       = "dst_2veL1FPpuXxUaZcFaEC57BfpcR3"
		idStepInactive    = "dst_2veL1FPpuXxUaZcFaEC57BfpcR4"
		idStepNoCoords    = "dst_2veL1FPpuXxUaZcFaEC57BfpcR5"
		idStepNotFound    = "dst_2veL1FPpuXxUaZcFaEC57Bfpd53"
		futureDateSetDays = 2
	)

	futureDate := time.Now().Add(futureDateSetDays * 24 * time.Hour).UTC().Format(time.RFC3339)

	makeBody := func(stops []string, dateSet string) dto.CreateRouteReq {
		return dto.CreateRouteReq{
			IdDriver:    idDriver,
			DateSet:     dateSet,
			Stops:       stops,
			Name:        "Rota zona sul",
			Description: "Coletas da zona sul",
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("With a single stop", func(t *testing.T) {
			body := makeBody([]string{idStepOne}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotEmpty(t, resp["id_route"])
			assert.Equal(t, idDriver, resp["id_driver"])
			assert.Equal(t, body.Name, resp["name"])
			assert.Equal(t, body.Description, resp["description"])
			assert.Equal(t, string(entities.EnumRouteStatusPending), resp["status"])
			assert.Nil(t, resp["city"])
			assert.Nil(t, resp["neighborhood"])

			stops, ok := resp["stops"].([]interface{})
			assert.True(t, ok)
			assert.Len(t, stops, 1)

			stop := stops[0].(map[string]interface{})
			assert.Equal(t, idStepOne, stop["id_donation_step"])
			assert.Equal(t, float64(0), stop["stop_order"])
		})

		t.Run("With multiple stops ordered by the routing provider", func(t *testing.T) {
			body := makeBody([]string{idStepOne, idStepTwo, idStepThree}, futureDate)
			body.City = utils.StringPtr("Sao Paulo")
			body.Neighborhood = utils.StringPtr("Vila Mariana")

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.Equal(t, "Sao Paulo", resp["city"])
			assert.Equal(t, "Vila Mariana", resp["neighborhood"])

			stops, ok := resp["stops"].([]interface{})
			assert.True(t, ok)
			assert.Len(t, stops, 3)

			// stops come back ordered by stop_order and every position is used once
			ordered := make([]string, 0, len(stops))
			for index, item := range stops {
				stop := item.(map[string]interface{})
				assert.NotEmpty(t, stop["id_route_donation_step"])
				assert.Equal(t, resp["id_route"], stop["id_route"])

				stopOrder, ok := stop["stop_order"].(float64)
				assert.True(t, ok)
				assert.Equal(t, float64(index), stopOrder)

				ordered = append(ordered, stop["id_donation_step"].(string))
			}
			assert.ElementsMatch(t, body.Stops, ordered)
			// the mocked provider swaps the two first stops, so the order sent
			// on the request is not the order stored on the route
			assert.Equal(t, []string{idStepTwo, idStepOne, idStepThree}, ordered)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User does not have permission to create route", func(t *testing.T) {
			body := makeBody([]string{idStepOne}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create route", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Driver not found", func(t *testing.T) {
			body := makeBody([]string{idStepOne}, futureDate)
			body.IdDriver = "usr_2veL1FPpuXxUaZcFaEC57Bfpd53"

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Driver not found", resp["message"])
		})

		t.Run("User is not a driver", func(t *testing.T) {
			body := makeBody([]string{idStepOne}, futureDate)
			body.IdDriver = idCommonUser

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "User is not a driver", resp["message"])
			assert.Equal(t, "driver.invalid_type", resp["code"])
		})

		t.Run("Date set must be in the future", func(t *testing.T) {
			pastDate := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
			body := makeBody([]string{idStepOne}, pastDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Date set must be in the future", resp["message"])
			assert.Equal(t, "date_set.invalid", resp["code"])
		})

		t.Run("Stop is not a donation step id", func(t *testing.T) {
			body := makeBody([]string{idDriver}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Stop is not a donation step id", resp["message"])
			assert.Equal(t, "stops.invalid_id", resp["code"])
		})

		t.Run("Duplicated donation step on stops", func(t *testing.T) {
			body := makeBody([]string{idStepOne, idStepOne}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Duplicated donation step on stops", resp["message"])
			assert.Equal(t, "stops.duplicated", resp["code"])
		})

		t.Run("Donation step not found", func(t *testing.T) {
			body := makeBody([]string{idStepNotFound}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation step not found", resp["message"])
		})

		t.Run("Donation is not active", func(t *testing.T) {
			body := makeBody([]string{idStepInactive}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Donation is not active", resp["message"])
			assert.Equal(t, "donation.inactive", resp["code"])
		})

		t.Run("Donation step without coordinates", func(t *testing.T) {
			body := makeBody([]string{idStepNoCoords}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Donation step does not have an address with coordinates", resp["message"])
			assert.Equal(t, "donation_step.missing_coordinates", resp["code"])
		})

		t.Run("Empty stops", func(t *testing.T) {
			body := makeBody([]string{}, futureDate)

			status, _ := fluxgo.RunTestRequest(app, "POST", endpoint, body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("User not found", func(t *testing.T) {
			body := makeBody([]string{idStepOne}, futureDate)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpoint, body, &utils.InvalidTestHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "User not found", resp["message"])
		})
	})
}
