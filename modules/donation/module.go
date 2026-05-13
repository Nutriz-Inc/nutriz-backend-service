package donation

import (
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/modules/donation/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("donation")

	// donation point
	mod.AddHandler(handlers.HandlerListDonationPointsStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListDonationPoints) error {
		return mod.HttpRoute(
			f,
			"/public",
			"GET",
			"/donation/point",
			fluxgo.RouteIncome{
				Entity:    dto.ListDonationPointsReq{},
				FromQuery: true,
				Validate:  true,
				Cache:     redis,
				CacheTTL:  time.Hour,
			},
			handler.HandleHttp,
		)
	})

	// donation
	mod.AddHandler(handlers.HandlerListDonationsStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListDonations) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/donation",
			fluxgo.RouteIncome{
				Entity:     dto.ListDonationReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
				CacheTTL:   time.Hour,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerGetDonationStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerGetDonation) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/donation/:id",
			fluxgo.RouteIncome{
				Entity:     dto.GetDonationReq{},
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
				CacheTTL:   time.Hour,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerCreateDonationStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateDonation) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/donation",
			fluxgo.RouteIncome{
				Entity:          dto.CreateDonationReq{},
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/donation"},
			},
			handler.HandleHttp,
		)
	})

	return mod
}
