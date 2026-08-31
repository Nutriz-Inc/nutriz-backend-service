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

func values(orders map[string]int) []int {
	out := make([]int, 0, len(orders))
	for _, order := range orders {
		out = append(out, order)
	}

	return out
}

func TestRemoveRouteStop(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	driverHeaders := &utils.TestHeadersDriver
	commonHeaders := &utils.TestHeaders

	const (
		idDriver    = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		idStepOne   = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idStepTwo   = "dst_2veL1FPpuXxUaZcFaEC57BfpcR2"
		idStepThree = "dst_2veL1FPpuXxUaZcFaEC57BfpcR3"
		idStopFake  = "rds_2veL1FPpuXxUaZcFaEC57Bfpd53"
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
				Stops:       stops,
				Name:        name,
				Description: "Rota criada para remocao de parada",
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

	env := fluxgo.ParseEnv[config.Env](fluxgo.EnvOptions{LoadFromFile: fluxgo.Pointer(".env.development"), Validate: true})
	db, err := sqlx.Connect("postgres", env.Database.Dsn)
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	remainingStops := func(t *testing.T, idRoute string) map[string]int {
		rows := []struct {
			IdRouteDonationStep string `db:"id_route_donation_step"`
			StopOrder           int    `db:"stop_order"`
		}{}

		err := db.Select(
			&rows,
			`SELECT id_route_donation_step, stop_order
			 FROM route_donation_step
			 WHERE id_route = $1 AND removed_at IS NULL
			 ORDER BY stop_order ASC`,
			idRoute,
		)
		assert.NoError(t, err)

		orders := map[string]int{}
		for _, row := range rows {
			orders[row.IdRouteDonationStep] = row.StopOrder
		}

		return orders
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Removes a stop and reorders the remaining ones", func(t *testing.T) {
			created := createRoute(t, "Rota para remover parada", []string{idStepOne, idStepTwo, idStepThree})

			removed := created[idStepTwo]
			idRoute := removed["id_route"].(string)
			removedOrder := removed["stop_order"].(float64)

			status, resp := fluxgo.RunTestRequest(
				app,
				"DELETE",
				endpointOf(removed["id_route_donation_step"].(string)),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, resp["success"])

			orders := remainingStops(t, idRoute)

			assert.Len(t, orders, 2)
			assert.NotContains(t, orders, removed["id_route_donation_step"].(string))

			// every stop that came after the removed one moves one position up
			for idDonationStep, stop := range created {
				if idDonationStep == idStepTwo {
					continue
				}

				order := stop["stop_order"].(float64)
				expected := int(order)
				if order > removedOrder {
					expected = int(order) - 1
				}

				assert.Equal(t, expected, orders[stop["id_route_donation_step"].(string)])
			}

			assert.ElementsMatch(t, []int{0, 1}, values(orders))
		})

		t.Run("Removes the only stop of a route", func(t *testing.T) {
			created := createRoute(t, "Rota com uma parada", []string{idStepOne})

			stop := created[idStepOne]

			status, resp := fluxgo.RunTestRequest(
				app,
				"DELETE",
				endpointOf(stop["id_route_donation_step"].(string)),
				nil,
				adminHeaders,
			)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, true, resp["success"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission", func(t *testing.T) {
			created := createRoute(t, "Rota parada sem permissao comum", []string{idStepOne})
			stop := created[idStepOne]

			status, resp := fluxgo.RunTestRequest(
				app,
				"DELETE",
				endpointOf(stop["id_route_donation_step"].(string)),
				nil,
				commonHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to remove route stop", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Driver does not have permission", func(t *testing.T) {
			created := createRoute(t, "Rota parada sem permissao driver", []string{idStepOne})
			stop := created[idStepOne]

			status, resp := fluxgo.RunTestRequest(
				app,
				"DELETE",
				endpointOf(stop["id_route_donation_step"].(string)),
				nil,
				driverHeaders,
			)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Route stop not found", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "DELETE", endpointOf(idStopFake), nil, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route stop not found", resp["message"])
		})

		t.Run("Route stop already removed", func(t *testing.T) {
			created := createRoute(t, "Rota parada removida duas vezes", []string{idStepOne})
			idStop := created[idStepOne]["id_route_donation_step"].(string)

			status, _ := fluxgo.RunTestRequest(app, "DELETE", endpointOf(idStop), nil, adminHeaders)
			assert.Equal(t, http.StatusOK, status)

			status, resp := fluxgo.RunTestRequest(app, "DELETE", endpointOf(idStop), nil, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route stop not found", resp["message"])
		})

		t.Run("Invalid stop id", func(t *testing.T) {
			status, _ := fluxgo.RunTestRequest(app, "DELETE", endpointOf("invalid-id"), nil, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
		})
	})
}
