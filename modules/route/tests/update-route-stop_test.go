package tests

import (
	"fmt"
	"net/http"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRouteStop(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	driverHeaders := &utils.TestHeadersDriver
	commonHeaders := &utils.TestHeaders

	const (
		idDriver   = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		idStepOne  = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idStepTwo  = "dst_2veL1FPpuXxUaZcFaEC57BfpcR2"
		idStopFake = "rds_2veL1FPpuXxUaZcFaEC57Bfpd53"
	)

	futureDate := time.Now().UTC().AddDate(0, 0, 3).Format(time.RFC3339)

	createRoute := func(t *testing.T, name string, stops []string) map[string]map[string]interface{} {
		status, resp := fluxgo.RunTestRequest(
			app,
			"POST",
			"/internal/route",
			dto.CreateRouteReq{
				IdDriver:    idDriver,
				DateSet:     futureDate,
				Stops:       &stops,
				Name:        name,
				Description: "Rota criada para atualizacao de parada",
			},
			adminHeaders,
		)

		assert.Equal(t, http.StatusCreated, status)

		byDonationStep := map[string]map[string]interface{}{}
		for _, item := range resp["stops"].([]interface{}) {
			stop := item.(map[string]interface{})
			byDonationStep[stop["id_donation_step"].(string)] = stop
		}

		return byDonationStep
	}

	endpointOf := func(idStop string) string {
		return fmt.Sprintf("/internal/route/stop/%s", idStop)
	}

	routeEndpointOf := func(idRoute string) string {
		return fmt.Sprintf("/internal/route/%s", idRoute)
	}

	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})
	db, err := sqlx.Connect("postgres", env.Database.Dsn)
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	stopDates := func(t *testing.T, idStop string) (start *time.Time, end *time.Time) {
		row := struct {
			DateStart *time.Time `db:"date_start"`
			DateEnd   *time.Time `db:"date_end"`
		}{}

		err := db.Get(
			&row,
			`SELECT date_start, date_end
			 FROM route_donation_step
			 WHERE id_route_donation_step = $1 AND removed_at IS NULL`,
			idStop,
		)
		assert.NoError(t, err)

		return row.DateStart, row.DateEnd
	}

	stopStatus := func(t *testing.T, idStop string) string {
		var status string
		err := db.Get(
			&status,
			`SELECT status FROM route_donation_step WHERE id_route_donation_step = $1`,
			idStop,
		)
		assert.NoError(t, err)
		return status
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Driver sets date_start of a stop", func(t *testing.T) {
			created := createRoute(t, "Rota para iniciar parada", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			stop := resp["stop"].(map[string]interface{})
			assert.NotNil(t, stop["date_start"])
			assert.Nil(t, stop["date_end"])
			assert.Equal(t, "in_progress", stop["status"])

			start, end := stopDates(t, idStop)
			assert.NotNil(t, start)
			assert.Nil(t, end)
			assert.Equal(t, "in_progress", stopStatus(t, idStop))
		})

		t.Run("Driver reports an error on a stop", func(t *testing.T) {
			created := createRoute(t, "Rota parada com erro", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{HasError: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			stop := resp["stop"].(map[string]interface{})
			assert.Equal(t, "error", stop["status"])
			assert.Nil(t, stop["date_start"])
			assert.Equal(t, "error", stopStatus(t, idStop))
		})

		t.Run("Route date_end propagates to started stops only", func(t *testing.T) {
			created := createRoute(t, "Rota para finalizar com paradas", []string{idStepOne, idStepTwo})
			idRoute := created[idStepOne]["id_route"].(string)
			idStartedStop := created[idStepOne]["id_route_donation_step"].(string)
			idPendingStop := created[idStepTwo]["id_route_donation_step"].(string)

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStartedStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, _ = fluxgo.RunTestRequest(
				app,
				"PUT",
				routeEndpointOf(idRoute),
				dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, _ = fluxgo.RunTestRequest(
				app,
				"PUT",
				routeEndpointOf(idRoute),
				dto.UpdateRouteReq{
					DateEnd:      utils.BoolPtr(true),
					Mileage:      utils.Float64Ptr(12.3),
					UserFeedback: utils.StringPtr("Rota concluida"),
				},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			startedStart, startedEnd := stopDates(t, idStartedStop)
			assert.NotNil(t, startedStart)
			assert.NotNil(t, startedEnd)
			assert.Equal(t, "done", stopStatus(t, idStartedStop))

			_, pendingEnd := stopDates(t, idPendingStop)
			assert.Nil(t, pendingEnd)
			assert.Equal(t, "pending", stopStatus(t, idPendingStop))
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Admin does not have permission", func(t *testing.T) {
			created := createRoute(t, "Rota parada sem permissao adm", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				adminHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update route stop", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Common user does not have permission", func(t *testing.T) {
			created := createRoute(t, "Rota parada sem permissao comum", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("date_start must be sent", func(t *testing.T) {
			created := createRoute(t, "Rota parada sem campos", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{},
				driverHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route_stop.no_fields_to_update", resp["code"])
		})

		t.Run("Route stop not found", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStopFake),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route stop not found", resp["message"])
		})

		t.Run("Invalid stop id", func(t *testing.T) {
			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf("invalid-id"),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
		})

		t.Run("Date start can only be set once", func(t *testing.T) {
			created := createRoute(t, "Rota parada iniciada duas vezes", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idStop),
				dto.UpdateRouteStopReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route_stop.date_start_already_set", resp["code"])
		})
	})
}
