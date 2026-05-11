package handlers

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	dto "nutriz-backend-service/modules/user/dtos"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetUser struct {
	userRepo    *repositories.UserRepository
	addressRepo *repositories.AddressRepository
	babyRepo    *repositories.UserBabyRepository
}

func HandlerGetUserStart(userRepo *repositories.UserRepository, addressRepo *repositories.AddressRepository, babyRepo *repositories.UserBabyRepository) *HandlerGetUser {
	return &HandlerGetUser{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		babyRepo:    babyRepo,
	}
}

func (h *HandlerGetUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetUser) Execute(ctx c.Context, filters *dto.GetUserReq) (*dto.GetUserRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.IdUser != filters.ActionBy {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "user.forbidden")
	}

	var addresses *[]entities.Address
	var babies *[]entities.UserBaby

	if filters.ShowAddress {
		addresses, _, err = h.addressRepo.GetAddressesByUserId(ctx, filters.Id)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get addresses")
		}
	}

	if filters.ShowBaby {
		babies, _, err = h.babyRepo.GetUserBabyesByUserId(ctx, filters.Id)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get babies")
		}
	}

	return &dto.GetUserRes{
		User:      *user,
		Addresses: addresses,
		Babies:    babies,
	}, nil
}
