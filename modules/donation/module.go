package donation

import (
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/modules/donation/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("donation")

	mod.AddHandler(handlers.HandlerListDonationPointsStart)

	mod.AddHandler(handlers.HandlerListDonationsStart)

	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListDonations) error {
		return mod.HttpRoute(
			f,
			"/public",
			"GET",
			"/donation",
			fluxgo.RouteIncome{
				Entity: dto.ListDonationReq{},
				FromQuery: true,
				Cache: redis,
				CacheTTL: time.Hour,
			},
			handler.HandleHttp,
		)
	})

	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListDonationPoints) error {
		return mod.HttpRoute(
			f,
			"/public",
			"GET",
			"/donation/point",
			fluxgo.RouteIncome{
				Entity:    dto.ListDonationPointsReq{},
				FromQuery: true,
				Cache:     redis,
				CacheTTL:  time.Hour,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
