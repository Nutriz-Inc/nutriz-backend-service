package route

import (
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/modules/route/handlers"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("route")

	mod.AddHandler(handlers.HandlerListRoutesStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerListRoutes) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/route",
			fluxgo.RouteIncome{
				Entity:     dto.ListRoutesReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerCreateRouteStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerCreateRoute) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/route",
			fluxgo.RouteIncome{
				Entity:     dto.CreateRouteReq{},
				FromBody:   true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateRouteStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerUpdateRoute) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"PUT",
			"/route/:id_route",
			fluxgo.RouteIncome{
				Entity:     dto.UpdateRouteReq{},
				FromBody:   true,
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerRemoveRouteStopStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerRemoveRouteStop) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"DELETE",
			"/route/stop/:id_stop",
			fluxgo.RouteIncome{
				Entity:     dto.RemoveRouteStopReq{},
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateRouteStopStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerUpdateRouteStop) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"PUT",
			"/route/stop/:id_stop",
			fluxgo.RouteIncome{
				Entity:     dto.UpdateRouteStopReq{},
				FromBody:   true,
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerCreateRouteStopStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerCreateRouteStop) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/route/:id_route/stop",
			fluxgo.RouteIncome{
				Entity:     dto.CreateRouteStopReq{},
				FromBody:   true,
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerGetRouteStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerGetRoute) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/route/:id",
			fluxgo.RouteIncome{
				Entity:     dto.GetRouteReq{},
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
