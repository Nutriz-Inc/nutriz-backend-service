package user

import (
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/modules/user/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("user")

	//user
	mod.AddHandler(handlers.HandlerGetUserStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerGetUser) error {
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
				Cache:      redis,
				CacheTTL:   time.Hour,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerListUsersStart)
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
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerCreateUserStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateUser) error {
		return mod.HttpRoute(
			f,
			"/public",
			"POST",
			"/user",
			fluxgo.RouteIncome{
				Entity:          dto.CreateUserReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateUser) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/user",
			fluxgo.RouteIncome{
				Entity:          dto.CreateUserReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerRemoveUserStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerRemoveUser) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"DELETE",
			"/user/:id",
			fluxgo.RouteIncome{
				Entity:          dto.RemoveUserReq{},
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateUserStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerUpdateUser) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"PATCH",
			"/user/:id",
			fluxgo.RouteIncome{
				Entity:				dto.UpdateUserReq{},
				FromBody:			true,
				FromParam:			true,
				FromHeader:			true,
				Validate:			true,
				Cache:				redis,
				CacheInvalidate:	[]string{"/internal/user"},
			},
			handler.HandleHttp,
		)	
	})

	// baby
	mod.AddHandler(handlers.HandlerCreateUserBabyStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateUserBaby) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/user/baby",
			fluxgo.RouteIncome{
				Entity:          dto.CreateUserBabyReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerRemoveUserBabyStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerRemoveUserBaby) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"DELETE",
			"/user/baby/:id",
			fluxgo.RouteIncome{
				Entity:          dto.RemoveUserBabyReq{},
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateUserBabyStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerUpdateUserBaby) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"PUT",
			"/user/baby/:id",
			fluxgo.RouteIncome{
				Entity:          dto.UpdateUserBabyReq{},
				FromBody:        true,
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	// address
	mod.AddHandler(handlers.HandlerCreateAddressStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateAddress) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/user/address",
			fluxgo.RouteIncome{
				Entity:          dto.CreateAddressReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateAddressStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerUpdateAddress) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"PUT",
			"/user/address/:id",
			fluxgo.RouteIncome{
				Entity:          dto.UpdateAddressReq{},
				FromBody:        true,
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerRemoveAddressStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerRemoveAddress) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"DELETE",
			"/user/address/:id",
			fluxgo.RouteIncome{
				Entity:          dto.RemoveAddressReq{},
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	// consent log
	mod.AddHandler(handlers.HandlerCreateConsentLogStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateConsentLog) error {
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
				Cache:      redis,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
