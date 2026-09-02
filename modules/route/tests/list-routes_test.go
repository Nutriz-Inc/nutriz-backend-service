package tests

import (
	"fmt"
	"net/http"
	"net/url"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"testing"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestListRoutes(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	endpoint := "/internal/route?page=1&page_size=25"
	adminHeaders := &utils.TestHeadersAdmin
	nurseHeaders := &utils.TestHeadersNurse
	driverHeaders := &utils.TestHeadersDriver
	commonHeaders := &utils.TestHeaders

	const (
		idDriver     = "usr_2veL1FPpuXxUaZcFaEC57BfpcDR"
		driverName   = "Carlos Motorista"
		idStepOne    = "dst_2veL1FPpuXxUaZcFaEC57BfpcR1"
		idOtherUser  = "usr_2veL1FPpuXxUaZcFaEC57BfpcKL"
		routeDateSet = 4
		routeName    = "Rota da listagem"
		routeCity    = "Sao Paulo"
		routeHood    = "Bela Vista"
	)

	dateSet := time.Now().UTC().AddDate(0, 0, routeDateSet)
	stops := []string{idStepOne}

	_, listSetupResp := fluxgo.RunTestRequest(
		app,
		"POST",
		"/internal/route",
		dto.CreateRouteReq{
			IdDriver:     idDriver,
			DateSet:      dateSet.Format(time.RFC3339),
			Stops:        &stops,
			Name:         routeName,
			Description:  "Rota criada para a listagem",
			City:         utils.StringPtr(routeCity),
			Neighborhood: utils.StringPtr(routeHood),
		},
		adminHeaders,
	)
	// deferred (not t.Cleanup) so it runs while the app is still up, before
	// the fx.RequireStop() defer tears it down
	if id, ok := listSetupResp["id_route"].(string); ok {
		defer cancelRouteNow(app, id)
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("No filters", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			first := fluxgo.ConvertToMap(data[0])
			assert.NotEmpty(t, first["id_route"])
			assert.NotEmpty(t, first["name"])
			assert.Equal(t, idDriver, first["id_driver"])
			assert.Equal(t, driverName, first["driver_name"])

			assert.Equal(t, float64(25), body["page_size"])
			assert.Equal(t, float64(1), body["page"])
			assert.GreaterOrEqual(t, body["total"], float64(1))
		})

		t.Run("Nurse can list routes", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, nurseHeaders)

			assert.Equal(t, int(http.StatusOK), status)
			assert.GreaterOrEqual(t, len(fluxgo.ConvertToList(body["data"])), 1)
		})

		t.Run("Driver can list routes", func(t *testing.T) {
			status, body := fluxgo.RunTestRequest(app, "GET", endpoint, nil, driverHeaders)

			assert.Equal(t, int(http.StatusOK), status)
			assert.GreaterOrEqual(t, len(fluxgo.ConvertToList(body["data"])), 1)
		})

		t.Run("id_driver filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&id_driver=%s", endpoint, idDriver)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, idDriver, row["id_driver"])
			}
		})

		t.Run("id_driver filter without routes", func(t *testing.T) {
			route := fmt.Sprintf("%s&id_driver=%s", endpoint, idOtherUser)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)
			assert.Len(t, fluxgo.ConvertToList(body["data"]), 0)
			assert.Equal(t, float64(0), body["total"])
		})

		t.Run("driver_name filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&driver_name=Carlos", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Contains(t, row["driver_name"], "Carlos")
			}
		})

		t.Run("status filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&status=%s", endpoint, entities.EnumRouteStatusPending)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, string(entities.EnumRouteStatusPending), row["status"])
			}
		})

		t.Run("date_set filter", func(t *testing.T) {
			date := dateSet.Format("2006-01-02")
			route := fmt.Sprintf("%s&date_set=%s", endpoint, date)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Contains(t, row["date_set"], date)
			}
		})

		t.Run("name filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&name=listagem", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Contains(t, row["name"], "listagem")
			}
		})

		t.Run("city filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&city=%s", endpoint, url.QueryEscape(routeCity))

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, routeCity, row["city"])
			}
		})

		t.Run("neighborhood filter", func(t *testing.T) {
			route := fmt.Sprintf("%s&neighborhood=%s", endpoint, url.QueryEscape(routeHood))

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)

			data := fluxgo.ConvertToList(body["data"])
			assert.GreaterOrEqual(t, len(data), 1, "unexpected data length")

			for _, item := range data {
				row := fluxgo.ConvertToMap(item)
				assert.Equal(t, routeHood, row["neighborhood"])
			}
		})

		t.Run("name filter without routes", func(t *testing.T) {
			route := fmt.Sprintf("%s&name=rota-inexistente", endpoint)

			status, body := fluxgo.RunTestRequest(app, "GET", route, nil, adminHeaders)

			assert.Equal(t, int(http.StatusOK), status)
			assert.Len(t, fluxgo.ConvertToList(body["data"]), 0)
			assert.Equal(t, float64(0), body["total"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User does not have permission to list routes", func(t *testing.T) {
			status, resp := fluxgo.RunTestRequest(app, "GET", endpoint, nil, commonHeaders)

			assert.Equal(t, int(http.StatusForbidden), status)
			assert.Equal(t, "User does not have permission to list routes", resp["message"])
			assert.Equal(t, "user.forbidden", resp["code"])
		})
	})
}
