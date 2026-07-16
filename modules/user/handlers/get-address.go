package handlers

import (
	c "context"
	"nutriz-backend-service/shared/repositories"

	dto "nutriz-backend-service/modules/user/dtos"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetAddress struct {
	addressRepo *repositories.AddressRepository
}

func HandlerGetAddressStart(addressRepo *repositories.AddressRepository) *HandlerGetAddress {
	return &HandlerGetAddress{
		addressRepo,
	}
}

func (h *HandlerGetAddress) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetAddressReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetAddress) Execute(ctx c.Context, data *dto.GetAddressReq) (*dto.GetAddressRes, *fluxgo.GlobalError) {
	address, err := h.addressRepo.GetAddressById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address")
	}
	if address == nil {
		return nil, fluxgo.ErrorNotFound("Address not found")
	}

	return &dto.GetAddressRes{
		Address: *address,
	}, nil
}
