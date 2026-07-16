package handlers

import (
	c "context"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateAddress struct {
	config      *config.Env
	addressRepo *repositories.AddressRepository
	userRepo    *repositories.UserRepository
}

func HandlerCreateAddressStart(
	config *config.Env,
	addressRepo *repositories.AddressRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateAddress {
	return &HandlerCreateAddress{
		config,
		addressRepo,
		userRepo,
	}
}

func (h *HandlerCreateAddress) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateAddressReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateAddress) Execute(ctx c.Context, data *dto.CreateAddressReq) (*dto.CreateAddressRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to create address", "user.forbidden")
	}

	_, totalAdresses, err := h.addressRepo.GetAddressesByUserId(ctx, user.IdUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get addresses by user id")
	}
	if totalAdresses >= entities.MAX_ADDRESS_QUANTITY_PER_USER {
		return nil, fluxgo.ErrorBadRequest(fmt.Sprintf("User can have up to %d addresses", entities.MAX_ADDRESS_QUANTITY_PER_USER), "address.max_quantity_reached")
	}

	// addressWithSameZipcode, err := h.addressRepo.GetAddressByZipcodeAndIdUser(ctx, data.ZipCode, user.IdUser)
	// if err != nil {
	// 	return nil, fluxgo.ErrorInternalError("Error to get address by zipcode")
	// }
	// if addressWithSameZipcode != nil {
	// 	return nil, fluxgo.ErrorBadRequest("Address with same zipcode already exists", "address.already_exists")
	// }

	addressData, err := utils.GetAddressByZipCodeOptionalCoordinates(ctx, data.ZipCode, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateAddressRepositoryReq{
		IdAddress:    idAddress,
		IdUser:       user.IdUser,
		Zipcode:      data.ZipCode,
		Street:       addressData.Street,
		Number:       data.Number,
		City:         addressData.City,
		State:        addressData.State,
		Neighborhood: addressData.Neighborhood,
		Complement:   data.Complement,
		Latitude:     addressData.Latitude,
		Longitude:    addressData.Longitude,
	}

	err = h.addressRepo.CreateAddress(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create address")
	}

	address, err := h.addressRepo.GetAddressById(ctx, idAddress)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address")
	}
	if address == nil {
		return nil, fluxgo.ErrorNotFound("Address not found")
	}

	return &dto.CreateAddressRes{
		Address: *address,
	}, nil
}
