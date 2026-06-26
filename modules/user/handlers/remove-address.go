package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerRemoveAddress struct {
	addressRepo *repositories.AddressRepository
	userRepo    *repositories.UserRepository
}

func HandlerRemoveAddressStart(
	addressRepo *repositories.AddressRepository,
	userRepo *repositories.UserRepository,
) *HandlerRemoveAddress {
	return &HandlerRemoveAddress{
		addressRepo,
		userRepo,
	}
}

func (h *HandlerRemoveAddress) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveAddressReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveAddress) Execute(ctx c.Context, data *dto.RemoveAddressReq) (*dto.RemoveAddressRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to delete address", "user.forbidden")
	}

	address, err := h.addressRepo.GetAddressById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address")
	}
	if address == nil {
		return nil, fluxgo.ErrorNotFound("Address not found")
	}

	if *address.IdUser != user.IdUser {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "address.forbidden")
	}

	err = h.addressRepo.RemoveAddress(ctx, data.Id, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to remove address")
	}

	address, err = h.addressRepo.GetAddressById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address")
	}

	return &dto.RemoveAddressRes{
		Success: address == nil,
	}, nil
}
