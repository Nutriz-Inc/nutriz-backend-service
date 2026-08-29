package route

import (
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/modules/route/handlers"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("route")

	// route
	mod.AddHandler(handlers.HandlerCreateRouteStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateRoute) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/route",
			fluxgo.RouteIncome{
				Entity:          dto.CreateRouteReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/route"},
			},
			handler.HandleHttp,
		)
	})

	return mod
}
