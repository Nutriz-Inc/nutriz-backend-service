package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

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
	if user.Type != entities.EnumUserTypeDonor {
		return nil, utils.ErrorForbidden("User does not have permission to create baby", "user.forbidden")
	}

	layout := "2006-01-02"

	birthDateParsed, err := time.Parse(layout, data.BirthDate)
	if err != nil {
		return nil, fluxgo.ErrorBadRequest("Invalid birth date format. Use YYYY-MM-DD HH:MM:SS", "user_baby.invalid_format")
	}

	if birthDateParsed.After(time.Now()) {
		return nil, fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user_baby.invalid_birth_date")
	}

	userBabyId := utils.IdGenerate(utils.UserBabyEntity)

	err = h.userBabyRepo.CreateUserBaby(ctx, &repositories.CreateUserBabyRepositoryReq{
		IdUserBaby: userBabyId,
		IdUser:     data.ActionBy,
		Name:       data.Name,
		BirthDate:  birthDateParsed,
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
