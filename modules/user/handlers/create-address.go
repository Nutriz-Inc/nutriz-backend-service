package handlers

import (
	"context"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
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

	_, totalAdresses, err := h.addressRepo.GetAddressesByUserId(ctx, user.IdUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get addresses by user id")
	}
	if totalAdresses >= entities.MAX_ADDRESS_QUANTITY_PER_USER {
		return nil, fluxgo.ErrorBadRequest(fmt.Sprintf("User can have up to %d addresses", entities.MAX_ADDRESS_QUANTITY_PER_USER), "address.max_quantity_reached")
	}

	addressWithSameZipcode, err := h.addressRepo.GetAddressByZipcode(ctx, data.ZipCode)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get address by zipcode")
	}
	if addressWithSameZipcode != nil {
		return nil, fluxgo.ErrorBadRequest("Address with same zipcode already exists", "address.already_exists")
	}

	coordinates, err := h.GetCoordinatesByZipCode(ctx, data.ZipCode)
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
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
		Latitude:     coordinates.Latitude,
		Longitude:    coordinates.Longitude,
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

func (h *HandlerCreateAddress) GetCoordinatesByZipCode(ctx context.Context, zipcode string) (*dto.Coordinates, error) {
	provider, err := location.NewLocationProvider(h.config)
	if err != nil {
		return nil, fmt.Errorf("error to initialize location provider: %v", err)
	}

	addressData, err := provider.GetAddressByZipCode(ctx, zipcode)
	if err != nil {
		return nil, fmt.Errorf("error getting address by zipcode: %v", err)
	}

	res := &dto.Coordinates{}

	if addressData.Location != nil && addressData.Location.Coordinates != nil && addressData.Location.Coordinates.Latitude != nil && addressData.Location.Coordinates.Longitude != nil {
		res.Latitude = utils.Float64Ptr(utils.StringToFloat64(*addressData.Location.Coordinates.Latitude))
		res.Longitude = utils.Float64Ptr(utils.StringToFloat64(*addressData.Location.Coordinates.Longitude))

		return res, nil
	}

	query := fmt.Sprintf(
		"%s %s %s Brazil",
		addressData.Street,
		addressData.City,
		addressData.State,
	)

	coordinates, err := provider.GetCoordinatesByAddress(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting coordinates by address: %v", err)
	}

	res.Latitude = utils.Float64Ptr(utils.StringToFloat64(coordinates.Lat))
	res.Longitude = utils.Float64Ptr(utils.StringToFloat64(coordinates.Lon))

	return res, nil
}
