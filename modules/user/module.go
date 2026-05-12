package user

import (
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/modules/user/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("user")

	mod.AddHandler(handlers.HandlerCreateUserBabyStart)
	mod.AddHandler(handlers.HandlerListUsersStart)
	mod.AddHandler(handlers.HandlerGetUserStart)
	mod.AddHandler(handlers.HandlerCreateConsentLogStart)

	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListUsers) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/user",
			fluxgo.RouteIncome{
				Entity:     dto.ListUsersReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
				CacheTTL:   time.Hour,
			},
			handler.HandleHttp,
		)
	})

	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerCreateUserBaby) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/user/baby",
			fluxgo.RouteIncome{
				Entity:     dto.CreateUserBabyReq{},
				FromBody:   true,
				FromHeader: true,
				Validate:   true,
				Cache:      nil,
				CacheTTL:   0,
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
				Validate:   true,
				Cache:      nil,
				CacheTTL:   0,
			},
			handler.HandleHttp,
		)
	})

	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerCreateConsentLog) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/user/consent",
			fluxgo.RouteIncome{
				Entity:     dto.CreateConsentReq{},
				FromBody:   true,
				FromHeader: true,
				Validate:   true,
				Cache:      nil,
				CacheTTL:   0,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
