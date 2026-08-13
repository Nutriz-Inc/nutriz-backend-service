package auth

import (
	dto "nutriz-backend-service/modules/auth/dtos"
	"nutriz-backend-service/modules/auth/handlers"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("auth")

	mod.AddHandler(handlers.HandlerLoginStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerLogin) error {
		return mod.HttpRoute(
			f,
			"/v1/public",
			"POST",
			"/auth/login",
			fluxgo.RouteIncome{
				Entity:   dto.LoginReq{},
				FromBody: true,
				Validate: true,
				Cache:    nil,
				CacheTTL: 0,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
