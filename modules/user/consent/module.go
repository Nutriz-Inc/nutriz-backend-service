package consent

import (
	dto "nutriz-backend-service/modules/user/consent/dtos"
	"nutriz-backend-service/modules/user/consent/handlers"
	"nutriz-backend-service/modules/user/consent/usecases"
	"nutriz-backend-service/shared/repositories"

	fluxgo "github.com/MMortari/FluxGo"
)

func Module() *fluxgo.FluxModule {
	mod := fluxgo.Module("consent")

	mod.AddHandler(func(db *fluxgo.Database) *handlers.HandlerCreateConsent {
		consentRepo := repositories.ConsentRepositoryStart(db)
		createUseCase := usecases.NewCreateConsentUseCase(consentRepo)
		return handlers.NewHandlerCreateConsent(createUseCase)
	})

	mod.AddRoute(func(f *fluxgo.FluxGo, handler *handlers.HandlerCreateConsent) error {
		return mod.HttpRoute(
			f,
			"/internal",
			"POST",
			"/consent",
			fluxgo.RouteIncome{
				Entity:   dto.CreateConsentReq{},
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
