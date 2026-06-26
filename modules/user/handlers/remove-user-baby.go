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

type HandlerRemoveUserBaby struct {
	userBabyRepo *repositories.UserBabyRepository
	userRepo     *repositories.UserRepository
}

func HandlerRemoveUserBabyStart(
	userBabyRepo *repositories.UserBabyRepository,
	userRepo *repositories.UserRepository,
) *HandlerRemoveUserBaby {
	return &HandlerRemoveUserBaby{
		userBabyRepo,
		userRepo,
	}
}

func (h *HandlerRemoveUserBaby) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveUserBabyReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveUserBaby) Execute(ctx c.Context, data *dto.RemoveUserBabyReq) (*dto.RemoveUserBabyRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to delete baby", "user.forbidden")
	}

	userBaby, err := h.userBabyRepo.GetUserBabyById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}
	if userBaby == nil {
		return nil, fluxgo.ErrorNotFound("User baby not found")
	}

	if userBaby.IdUser != data.ActionBy {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "user_baby.forbidden")
	}

	err = h.userBabyRepo.RemoveUserBaby(ctx, data.Id, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to remove user baby")
	}

	userBaby, err = h.userBabyRepo.GetUserBabyById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}

	return &dto.RemoveUserBabyRes{
		Success: userBaby == nil,
	}, nil
}
