package dashboard

import (
	dto "nutriz-backend-service/modules/dashboard/dtos"
	"nutriz-backend-service/modules/dashboard/handlers"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("dashboard")

	mod.AddHandler(handlers.HandlerGetDashboardStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerGetDashboard) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"GET",
			"/dashboard",
			fluxgo.RouteIncome{
				Entity:     dto.GetDashboardReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
