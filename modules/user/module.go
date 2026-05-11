package user

import (
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/modules/user/handlers"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("user")

	mod.AddHandler(handlers.HandlerCreateUserBabyStart)
	mod.AddHandler(handlers.HandlerGetUserStart)

	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateUserBaby) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/userbaby",
			fluxgo.RouteIncome{
				Entity: dto.CreateUserBabyReq{},
				FromBody: true,
				FromHeader: true,
				Validate: true,
				Cache: nil,
				CacheTTL: 0,
			},
			handler.HandleHttp,
		)
	})

	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerGetUser) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/user/:id",
			fluxgo.RouteIncome{
				Entity:     dto.GetUserReq{},
				FromParam:  true,
				FromQuery:  true,
				FromHeader: true,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
