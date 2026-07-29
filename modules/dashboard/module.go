package dashboard

import (
	dto "nutriz-backend-service/modules/dashboard/dtos"
	"nutriz-backend-service/modules/dashboard/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("dashboard")

	mod.AddHandler(handlers.HandlerGetDashboardStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerGetDashboard) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/dashboard",
			fluxgo.RouteIncome{
				Entity:     dto.GetDashboardReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
				CacheTTL:   5 * time.Minute,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
