package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerUpdateUserBaby struct {
	userBabyRepo *repositories.UserBabyRepository
	userRepo     *repositories.UserRepository
}

func HandlerUpdateUserBabyStart(
	userBabyRepo *repositories.UserBabyRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateUserBaby {
	return &HandlerUpdateUserBaby{userBabyRepo, userRepo}
}

func (h *HandlerUpdateUserBaby) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateUserBabyReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateUserBaby) Execute(ctx c.Context, data *dto.UpdateUserBabyReq) (*dto.UpdateUserBabyRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanUpdateBaby {
		return nil, utils.ErrorForbidden("User does not have permission to update baby", "user.forbidden")
	}

	userBaby, err := h.userBabyRepo.GetUserBabyById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}
	if userBaby == nil {
		return nil, fluxgo.ErrorNotFound("User baby not found")
	}
	if userBaby.IdUser != user.IdUser {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "user_baby.forbidden")
	}

	fieldsToUpdate := 0
	req := repositories.UpdateUserBabyRepositoryReq{
		IdUserBaby: data.Id,
		IdUser:     data.ActionBy,
	}

	validator := data.ValidateUpdateUserBabyOptionalFields()

	canUpdateName := validator.HasName && (userBaby.Name == nil || *data.Name != *userBaby.Name)
	canUpdateBirthDate := validator.HasBirthDate && (userBaby.BirthDate.Format("2006-01-02") != *data.BirthDate)

	if canUpdateName {
		req.Name = data.Name
		fieldsToUpdate++
	}

	if canUpdateBirthDate {
		if utils.IsFutureDate(*data.BirthDate) {
			return nil, fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user_baby.invalid_birth_date")
		}

		birthDateTime, err := utils.StringToDate(*data.BirthDate)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest("Invalid birth date format", "user_baby.invalid_birth_date_format")
		}

		req.BirthDate = birthDateTime
		fieldsToUpdate++
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be different to update", "no_fields_to_update")
	}

	err = h.userBabyRepo.UpdateUserBaby(ctx, req)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update user baby")
	}

	userBaby, err = h.userBabyRepo.GetUserBabyById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}
	if userBaby == nil {
		return nil, fluxgo.ErrorNotFound("User baby not found")
	}

	return &dto.UpdateUserBabyRes{UserBaby: *userBaby}, nil
}
