package handlers

import (
	c "context"
	"fmt"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateUserBaby struct {
	userBabyRepo *repositories.UserBabyRepository
	userRepo     *repositories.UserRepository
}

func HandlerCreateUserBabyStart(userBabyRepo *repositories.UserBabyRepository,
	userRepo *repositories.UserRepository) *HandlerCreateUserBaby {
	return &HandlerCreateUserBaby{userBabyRepo, userRepo}
}

func (h *HandlerCreateUserBaby) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateUserBabyReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateUserBaby) Execute(ctx c.Context, data *dto.CreateUserBabyReq) (*dto.CreateUserBabyRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to create baby", "user.forbidden")
	}

	_, babyquantity, err := h.userBabyRepo.GetUserBabiesByUserId(ctx, user.IdUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user babies by user id")
	}
	if babyquantity >= entities.MAX_BABY_QUANTITY_PER_USER {
		return nil, fluxgo.ErrorBadRequest(fmt.Sprintf("User can have up to %d babies", entities.MAX_BABY_QUANTITY_PER_USER), "user_baby.max_quantity_reached")
	}

	if utils.IsFutureDate(data.BirthDate) {
		return nil, fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user_baby.invalid_birth_date")
	}

	birthDateTime, err := utils.StringToDate(data.BirthDate)
	if err != nil {
		return nil, fluxgo.ErrorBadRequest("Invalid birth date format", "user_baby.invalid_birth_date_format")
	}

	userBabyId := utils.IdGenerate(utils.UserBabyEntity)

	err = h.userBabyRepo.CreateUserBaby(ctx, &repositories.CreateUserBabyRepositoryReq{
		IdUserBaby: userBabyId,
		IdUser:     data.ActionBy,
		Name:       data.Name,
		BirthDate:  *birthDateTime,
	})
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create baby")
	}

	userBaby, err := h.userBabyRepo.GetUserBabyById(ctx, userBabyId)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}
	if userBaby == nil {
		return nil, fluxgo.ErrorNotFound("User baby not found")
	}

	return &dto.CreateUserBabyRes{UserBaby: *userBaby}, nil
}
