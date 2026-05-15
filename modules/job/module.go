package job

import (
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/modules/job/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("job")

	mod.AddHandler(handlers.HandlerListJobsStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerListJobs) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/job",
			fluxgo.RouteIncome{
				Entity:     dto.ListJobsReq{},
				FromQuery:  true,
				FromHeader: true,
				Validate:   true,
				Cache:      redis,
			},
			handler.HandleHttp,
		)
	})

	mod.AddHandler(handlers.HandlerGetJobStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerGetJob) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"GET",
			"/job/:id",
			fluxgo.RouteIncome{
				Entity:     dto.GetJobReq{},
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

	mod.AddHandler(handlers.HandlerCreateJobStart)
	mod.AddRoute(func(f *fluxgo.FluxGo, redis *fluxgo.Redis, handler *handlers.HandlerCreateJob) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/job",
			fluxgo.RouteIncome{
				Entity:          dto.CreateJobReq{},
				FromBody:        true,
				FromHeader:      true,
				Validate:        true,
				Cache:           redis,
				CacheInvalidate: []string{"/internal/job"},
			},
			handler.HandleHttp,
		)
	})

	return mod
}
