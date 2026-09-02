package tests

import (
	"fmt"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
)

func cancelRouteNow(app *fluxgo.Http, idRoute string) {
	if idRoute == "" {
		return
	}

	_, _ = fluxgo.RunTestRequestRaw(
		app,
		"PUT",
		fmt.Sprintf("/internal/route/%s", idRoute),
		dto.UpdateRouteReq{
			Status:      utils.RouteStatusPtr(entities.EnumRouteStatusCanceled),
			Description: utils.StringPtr("cleanup: liberar donation steps"),
		},
		&utils.TestHeadersAdmin,
	)
}

func cancelRouteOnCleanup(t *testing.T, app *fluxgo.Http, idRoute interface{}) {
	t.Helper()

	id, _ := idRoute.(string)
	if id == "" {
		return
	}

	t.Cleanup(func() { cancelRouteNow(app, id) })
}
