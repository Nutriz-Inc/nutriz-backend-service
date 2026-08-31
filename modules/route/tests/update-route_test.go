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

func TestUpdateRoute(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	adminHeaders := &utils.TestHeadersAdmin
	driverHeaders := &utils.TestHeadersDriver
	nurseHeaders := &utils.TestHeadersNurse
	commonHeaders := &utils.TestHeaders

	const (
		idDriver     = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		idStepOne    = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idRouteFake  = "rot_2veL1FPpuXxUaZcFaEC57Bfpd53"
		routeDateSet = 3
	)

	futureDate := time.Now().UTC().AddDate(0, 0, routeDateSet).Format(time.RFC3339)

	createRoute := func(t *testing.T, name string) string {
		status, resp := fluxgo.RunTestRequest(
			app,
			"POST",
			"/internal/route",
			dto.CreateRouteReq{
				IdDriver:    idDriver,
				DateSet:     futureDate,
				Stops:       []string{idStepOne},
				Name:        name,
				Description: "Rota criada para a atualizacao",
			},
			adminHeaders,
		)

		assert.Equal(t, http.StatusCreated, status)

		idRoute, ok := resp["id_route"].(string)
		assert.True(t, ok)

		return idRoute
	}

	endpointOf := func(idRoute string) string {
		return fmt.Sprintf("/internal/route/%s", idRoute)
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Adm updates the route data", func(t *testing.T) {
			idRoute := createRoute(t, "Rota para editar")
			newDateSet := time.Now().UTC().AddDate(0, 0, 10).Format(time.RFC3339)

			body := dto.UpdateRouteReq{
				Name:         utils.StringPtr("Rota editada"),
				City:         utils.StringPtr("Santo Andre"),
				Neighborhood: utils.StringPtr("Centro"),
				Description:  utils.StringPtr("Descricao editada"),
				DateSet:      &newDateSet,
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, idRoute, resp["id_route"])
			assert.Equal(t, "Rota editada", resp["name"])
			assert.Equal(t, "Santo Andre", resp["city"])
			assert.Equal(t, "Centro", resp["neighborhood"])
			assert.Equal(t, "Descricao editada", resp["description"])
			assert.Contains(t, resp["date_set"], newDateSet[:10])
			assert.Equal(t, string(entities.EnumRouteStatusPending), resp["status"])
		})

		t.Run("Adm cancels the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota para cancelar")

			body := dto.UpdateRouteReq{
				Status:      utils.RouteStatusPtr(entities.EnumRouteStatusCanceled),
				Description: utils.StringPtr("Cancelada por falta de motorista"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, string(entities.EnumRouteStatusCanceled), resp["status"])
			assert.Equal(t, "Cancelada por falta de motorista", resp["description"])
		})

		t.Run("Driver starts the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota para iniciar")

			body := dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.NotNil(t, resp["date_start"])
			assert.Nil(t, resp["date_end"])
			assert.Equal(t, string(entities.EnumRouteStatusInProgress), resp["status"])
		})

		t.Run("Driver finishes the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota para finalizar")

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			body := dto.UpdateRouteReq{
				DateEnd:      utils.BoolPtr(true),
				Mileage:      utils.Float64Ptr(42.5),
				UserFeedback: utils.StringPtr("Coletas realizadas sem problemas"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.NotNil(t, resp["date_end"])
			assert.Equal(t, float64(42.5), resp["mileage"])
			assert.Equal(t, "Coletas realizadas sem problemas", resp["user_feedback"])
			assert.Equal(t, string(entities.EnumRouteStatusDone), resp["status"])
		})

		t.Run("Driver reports an error on the route", func(t *testing.T) {
			idRoute := createRoute(t, "Rota com erro")

			body := dto.UpdateRouteReq{
				Status:       utils.RouteStatusPtr(entities.EnumRouteStatusError),
				UserFeedback: utils.StringPtr("Van quebrou no meio da rota"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, string(entities.EnumRouteStatusError), resp["status"])
			assert.Equal(t, "Van quebrou no meio da rota", resp["user_feedback"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Common user does not have permission", func(t *testing.T) {
			idRoute := createRoute(t, "Rota sem permissao comum")

			body := dto.UpdateRouteReq{Name: utils.StringPtr("Nova rota")}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to update route", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Nurse does not have permission", func(t *testing.T) {
			idRoute := createRoute(t, "Rota sem permissao enfermeira")

			body := dto.UpdateRouteReq{Name: utils.StringPtr("Nova rota")}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, nurseHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "user.forbidden", resp["code"])
		})

		t.Run("Route not found", func(t *testing.T) {
			body := dto.UpdateRouteReq{Name: utils.StringPtr("Nova rota")}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRouteFake), body, adminHeaders)

			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, "Route not found", resp["message"])
		})

		t.Run("Canceled route cannot be updated", func(t *testing.T) {
			idRoute := createRoute(t, "Rota cancelada duas vezes")

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{
					Status:      utils.RouteStatusPtr(entities.EnumRouteStatusCanceled),
					Description: utils.StringPtr("Cancelada"),
				},
				adminHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{Name: utils.StringPtr("Rota renomeada")},
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.canceled_or_done", resp["code"])
		})

		t.Run("Adm cannot send driver fields", func(t *testing.T) {
			idRoute := createRoute(t, "Rota adm com campos de driver")

			body := dto.UpdateRouteReq{Mileage: utils.Float64Ptr(10)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "route.invalid_fields_for_adm", resp["code"])
		})

		t.Run("Adm cannot set status error", func(t *testing.T) {
			idRoute := createRoute(t, "Rota adm com status error")

			body := dto.UpdateRouteReq{Status: utils.RouteStatusPtr(entities.EnumRouteStatusError)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.invalid_status_for_adm", resp["code"])
		})

		t.Run("Description is required when canceling", func(t *testing.T) {
			idRoute := createRoute(t, "Rota cancelada sem descricao")

			body := dto.UpdateRouteReq{Status: utils.RouteStatusPtr(entities.EnumRouteStatusCanceled)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.description_required", resp["code"])
		})

		t.Run("Date set must be in the future", func(t *testing.T) {
			idRoute := createRoute(t, "Rota com data passada")
			pastDate := time.Now().UTC().AddDate(0, 0, -1).Format(time.RFC3339)

			body := dto.UpdateRouteReq{DateSet: &pastDate}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "date_set.invalid", resp["code"])
		})

		t.Run("Adm must send at least one field", func(t *testing.T) {
			idRoute := createRoute(t, "Rota adm sem campos")

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), dto.UpdateRouteReq{}, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.no_fields_to_update", resp["code"])
		})

		t.Run("Driver cannot send adm fields", func(t *testing.T) {
			idRoute := createRoute(t, "Rota driver com campos de adm")

			body := dto.UpdateRouteReq{Name: utils.StringPtr("Rota renomeada")}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "route.invalid_fields_for_driver", resp["code"])
		})

		t.Run("Driver cannot set status canceled", func(t *testing.T) {
			idRoute := createRoute(t, "Rota driver cancelando")

			body := dto.UpdateRouteReq{Status: utils.RouteStatusPtr(entities.EnumRouteStatusCanceled)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.invalid_status_for_driver", resp["code"])
		})

		t.Run("Date start and date end cannot be sent together", func(t *testing.T) {
			idRoute := createRoute(t, "Rota com as duas datas")

			body := dto.UpdateRouteReq{
				DateStart:    utils.BoolPtr(true),
				DateEnd:      utils.BoolPtr(true),
				Mileage:      utils.Float64Ptr(10),
				UserFeedback: utils.StringPtr("Feedback"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.date_start_and_date_end", resp["code"])
		})

		t.Run("Route was not started yet", func(t *testing.T) {
			idRoute := createRoute(t, "Rota finalizada sem iniciar")

			body := dto.UpdateRouteReq{
				DateEnd:      utils.BoolPtr(true),
				Mileage:      utils.Float64Ptr(10),
				UserFeedback: utils.StringPtr("Feedback"),
			}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.not_started", resp["code"])
		})

		t.Run("Mileage and user feedback are required on date end", func(t *testing.T) {
			idRoute := createRoute(t, "Rota finalizada sem quilometragem")

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			body := dto.UpdateRouteReq{DateEnd: utils.BoolPtr(true)}

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), body, driverHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.mileage_and_user_feedback_required", resp["code"])
		})

		t.Run("Date start can only be set once", func(t *testing.T) {
			idRoute := createRoute(t, "Rota iniciada duas vezes")

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)
			assert.Equal(t, http.StatusOK, status)

			status, resp := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				dto.UpdateRouteReq{DateStart: utils.BoolPtr(true)},
				driverHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.date_start_already_set", resp["code"])
		})

		t.Run("Driver must send at least one field", func(t *testing.T) {
			idRoute := createRoute(t, "Rota driver sem campos")

			status, resp := fluxgo.RunTestRequest(app, "PUT", endpointOf(idRoute), dto.UpdateRouteReq{}, driverHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "route.no_fields_to_update", resp["code"])
		})

		t.Run("Invalid status", func(t *testing.T) {
			idRoute := createRoute(t, "Rota com status invalido")

			status, _ := fluxgo.RunTestRequest(
				app,
				"PUT",
				endpointOf(idRoute),
				map[string]any{"status": "in_progress"},
				adminHeaders,
			)

			assert.Equal(t, http.StatusBadRequest, status)
		})
	})
}
