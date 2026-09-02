package tests

import (
	"fmt"
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

func TestCreateRouteStop(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	driverHeaders := &utils.TestHeadersDriver
	commonHeaders := &utils.TestHeaders

	const (
		idDriver       = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		idStepOne      = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idStepTwo      = "dst_2veL1FPpuXxUaZcFaEC57BfpcR2"
		idStepThree    = "dst_2veL1FPpuXxUaZcFaEC57BfpcR3"
		idStepInactive = "dst_2veL1FPpuXxUaZcFaEC57BfpcR4"
		idStepFarAway  = "dst_2veL1FPpuXxUaZcFaEC57BfpcR6"
		idStepNotFound = "dst_2veL1FPpuXxUaZcFaEC57Bfpd53"
		idRouteFake    = "rot_2veL1FPpuXxUaZcFaEC57Bfpd53"
	)

	futureDate := time.Now().UTC().AddDate(0, 0, 3).Format(time.RFC3339)

	createRoute := func(t *testing.T, name string, stops []string, city, neighborhood *string) string {
		status, resp := fluxgo.RunTestRequest(
			app,
			"POST",
			"/internal/route",
			dto.CreateRouteReq{
				IdDriver:     idDriver,
				DateSet:      futureDate,
				Stops:        &stops,
				Name:         name,
				Description:  "Rota criada para adicionar parada",
				City:         city,
				Neighborhood: neighborhood,
			},
			adminHeaders,
		)

		assert.Equal(t, http.StatusCreated, status)

		return resp["id_route"].(string)
	}

	endpointOf := func(idRoute string) string {
		return fmt.Sprintf("/internal/route/%s/stop", idRoute)
	}

	body := func(idDonationStep string) dto.CreateRouteStopReq {
		return dto.CreateRouteStopReq{IdDonationStep: idDonationStep}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Adds a stop and reorders the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota para adicionar parada", []string{idStepOne, idStepTwo}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepThree), adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.Equal(t, idRoute, resp["id_route"])
			assert.NotNil(t, resp["updated_at"])

			estimatedTime, ok := resp["estimated_time"].(float64)
			assert.True(t, ok)
			assert.Greater(t, estimatedTime, float64(0))

			stops := fluxgo.ConvertToList(resp["stops"])
			assert.Len(t, stops, 3)

			// the stops come back ordered and keep a sequence without holes
			donationSteps := make([]string, 0, len(stops))
			for index, item := range stops {
				stop := fluxgo.ConvertToMap(item)

				assert.Equal(t, idRoute, stop["id_route"])
				assert.Equal(t, float64(index), stop["stop_order"])
				assert.Equal(t, string(entities.EnumRouteDonationStepStatusPending), stop["status"])

				donationSteps = append(donationSteps, stop["id_donation_step"].(string))
			}

			assert.ElementsMatch(t, []string{idStepOne, idStepTwo, idStepThree}, donationSteps)
		})

		t.Run("Adds a stop inside the city of the route", func(t *testing.T) {
			idRoute := createRoute(
				t,
				"Rota com cidade",
				[]string{idStepOne},
				utils.StringPtr("Sao Paulo"),
				nil,
			)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepThree), adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.Equal(t, "Sao Paulo", resp["city"])
			assert.Len(t, fluxgo.ConvertToList(resp["stops"]), 2)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada sem permissao comum", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepTwo), commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create route stop", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Driver does not have permission", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada sem permissao driver", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepTwo), driverHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Route not found", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRouteFake), body(idStepTwo), adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route not found", resp["message"])
		})

		t.Run("Canceled route cannot receive stops", func(t *testing.T) {
			idRoute := createRoute(t, "Rota cancelada com parada", []string{idStepOne}, nil, nil)

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				fmt.Sprintf("/internal/route/%s", idRoute),
				dto.UpdateRouteReq{
					Status:      utils.RouteStatusPtr(entities.EnumRouteStatusCanceled),
					Description: utils.StringPtr("Cancelada"),
				},
				adminHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepTwo), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.canceled", resp["code"])
		})

		t.Run("Donation step not found", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada inexistente", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepNotFound), adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Donation step not found", resp["message"])
		})

		t.Run("Stop is not a donation step id", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada com id invalido", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idDriver), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "stops.invalid_id", resp["code"])
		})

		t.Run("Donation step is already a stop of the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada duplicada", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepOne), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "stops.duplicated", resp["code"])
		})

		t.Run("Donation is not active", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada inativa", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepInactive), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "donation.inactive", resp["code"])
		})

		t.Run("Stop is not in the city of the route", func(t *testing.T) {
			idRoute := createRoute(
				t,
				"Rota parada fora da cidade",
				[]string{idStepOne},
				utils.StringPtr("Sao Paulo"),
				nil,
			)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepFarAway), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "stops.invalid_city", resp["code"])
		})

		t.Run("Stop is not in the neighborhood of the route", func(t *testing.T) {
			idRoute := createRoute(
				t,
				"Rota parada fora do bairro",
				[]string{idStepOne},
				nil,
				utils.StringPtr("Bela Vista"),
			)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepTwo), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "stops.invalid_neighborhood", resp["code"])
		})

		t.Run("Route takes longer than the maximum duration", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada longa demais", []string{idStepOne}, nil, nil)

			status, resp := fluxgo.RunTestRequest(app, "POST", endpointOf(idRoute), body(idStepFarAway), adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.max_duration_exceeded", resp["code"])
		})

		t.Run("Invalid donation step id", func(t *testing.T) {
			idRoute := createRoute(t, "Rota parada id invalido", []string{idStepOne}, nil, nil)

			status, _ := fluxgo.RunTestRequest(
				app,
				"POST",
				endpointOf(idRoute),
				map[string]any{"id_donation_step": "invalid-id"},
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
		})
	})
}
