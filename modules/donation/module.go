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
			"/v1/public",
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
	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerListDonations) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"GET",
			"/donation",
			fluxgo.RouteIncome{
				Entity:     dto.ListDonationReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerGetDonationStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerGetDonation) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
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
			"/v1/internal",
			"POST",
			"/donation",
			fluxgo.RouteIncome{
				Entity:          dto.CreateDonationReq{},
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/v1/internal/donation", "v1/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateDonationStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerUpdateDonation) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"PUT",
			"/donation/:id",
			fluxgo.RouteIncome{
				Entity:          dto.UpdateDonationReq{},
				FromBody:        true,
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/v1/internal/donation", "/v1/internal/user"},
			},
			handler.HandleHttp,
		)
	})

	// donation step
	mod.AddHandler(handlers.HandlerCreateDonationStepStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateDonationStep) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"POST",
			"/donation/step",
			fluxgo.RouteIncome{
				Entity:          dto.CreateDonationStepReq{},
				FromHeader:      true,
				FromBody:        true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/v1/internal/donation", "/v1/internal/user", "/v1/internal/donation/step/:id/timeline"},
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerUpdateDonationStepStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerUpdateDonationStep) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"PUT",
			"/donation/step/:id",
			fluxgo.RouteIncome{
				Entity:          dto.UpdateDonationStepReq{},
				FromBody:        true,
				FromParam:       true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/v1/internal/donation", "/v1/internal/user", "/v1/internal/donation/step/:id/timeline"},
			},
			handler.HandleHttp,
		)
	})

	// donation step timeline
	mod.AddHandler(handlers.HandlerListDonationStepTimelineStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListDonationStepTimeline) error {
		return mod.HttpRoute(
			f,
			"/v1/internal",
			"GET",
			"/donation/step/:id/timeline",
			fluxgo.RouteIncome{
				Entity:     dto.ListDonationStepTimelineReq{},
				FromParam:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
				CacheTTL:   time.Hour,
			},
			handler.HandleHttp,
		)
	})

	return mod
}
