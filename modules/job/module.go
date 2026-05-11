package job

import (
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/modules/job/handlers"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("job")

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

	return mod
}
