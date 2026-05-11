package handlers

import (
	"nutriz-backend-service/modules/consent/dtos"
	"nutriz-backend-service/modules/consent/usecases"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateConsent struct {
	useCase *usecases.CreateConsentUseCase
}

func NewHandlerCreateConsent(
	useCase *usecases.CreateConsentUseCase,
) *HandlerCreateConsent {
	return &HandlerCreateConsent{
		useCase: useCase,
	}
}

func (h *HandlerCreateConsent) HandleHttp(
	c *fiber.Ctx,
	income interface{},
) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {

	idUser, ok := c.Locals("id_user").(string)
	if !ok || idUser == "" {
		return nil, fluxgo.ErrorBadRequest("missing or invalid user identity", idUser)
	}

	resp, err := h.useCase.Execute(
		c.UserContext(),
		idUser,
		income.(*dtos.CreateConsentReq),
	)
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{
		Content: resp,
		Status:  201,
	}, nil
}
