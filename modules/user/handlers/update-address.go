package handlers

import (
	"context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerUpdateAddress struct {
	config      *config.Env
	addressRepo *repositories.AddressRepository
	userRepo    *repositories.UserRepository
}

func HandlerUpdateAddressStart(
	config *config.Env,
	addressRepo *repositories.AddressRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateAddress {
	return &HandlerUpdateAddress{
		config,
		addressRepo,
		userRepo,
	}
}

func (h *HandlerUpdateAddress) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateAddressReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateAddress) Execute(ctx context.Context, data *dto.UpdateAddressReq) (*dto.UpdateAddressRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to update address", "user.forbidden")
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

	fieldsToUpdate := 0
	req := repositories.UpdateAddressRepositoryReq{
		IdAddress: data.Id,
		IdUser:    data.ActionBy,
	}

	validator := data.ValidateUpdateAddressOptionalFields()

	canUpdateZipCode := validator.HasZipCode && (*data.ZipCode != address.Zipcode)
	canUpdateNumber := validator.HasNumber && (address.Number == nil || *data.Number != *address.Number)
	canUpdateComplement := validator.HasComplement && (address.Complement == nil || *data.Complement != *address.Complement)

	if canUpdateZipCode {
		addressWithSameZipcode, err := h.addressRepo.GetAddressByZipcodeAndIdUser(ctx, *data.ZipCode, user.IdUser)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get address by zipcode")
		}
		if addressWithSameZipcode != nil {
			return nil, fluxgo.ErrorBadRequest("Address with same zipcode already exists", "address.already_exists")
		}

		addressData, err := utils.GetAddressByZipCode(ctx, *data.ZipCode, h.config)
		if err != nil {
			return nil, fluxgo.ErrorInternalError(err.Error())
		}

		req.Longitude = addressData.Longitude
		req.Latitude = addressData.Latitude
		req.Zipcode = data.ZipCode
		req.City = &addressData.City
		req.State = &addressData.State
		req.Neighborhood = &addressData.Neighborhood
		req.Street = &addressData.Street

		fieldsToUpdate++
	}

	if canUpdateNumber {
		req.Number = data.Number
		fieldsToUpdate++
	}

	if canUpdateComplement {
		req.Complement = data.Complement
		fieldsToUpdate++
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be different to update", "no_fields_to_update")
	}

	err = h.addressRepo.UpdateAddress(ctx, req)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update address")
	}

	address, err = h.addressRepo.GetAddressById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address")
	}
	if address == nil {
		return nil, fluxgo.ErrorNotFound("Address not found")
	}

	return &dto.UpdateAddressRes{
		Address: *address,
	}, nil
}
