package handlers

import (
	"context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/provider/location"
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
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerCreateAddress) Execute(ctx context.Context, data *dto.CreateAddressReq) (*dto.CreateAddressRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	provider, err := location.NewLocationProvider(h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error initializing location provider")
	}

	addressData, err := provider.GetAddressByZipCode(ctx, data.ZipCode)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error getting address by zipcode")
	}

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateAddressRepositoryReq{
		IdAddress:    idAddress,
		IdUser:       user.IdUser,
		Zipcode:      data.ZipCode,
		Street:       data.Street,
		Number:       data.Number,
		City:         data.City,
		State:        data.State,
		Neighborhood: data.Neighborhood,
		Complement:   data.Complement,
		Latitude:     addressData.Location.Coordinates.Latitude,
		Longitude:    addressData.Location.Coordinates.Longitude,
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
