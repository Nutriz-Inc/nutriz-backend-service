package handlers

import (
	"context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerRemoveUserBaby struct {
	userBabyRepo	*repositories.UserBabyRepository
}

func HandlerRemoveUserBabyStart(
	userBabyRepo *repositories.UserBabyRepository,
) *HandlerRemoveUserBaby {
	return &HandlerRemoveUserBaby{
		userBabyRepo,
	}
}
func (h *HandlerRemoveUserBaby) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveUserBabyReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveUserBaby) Execute(ctx context.Context, data *dto.RemoveUserBabyReq) (*dto.RemoveUserBabyRes, *fluxgo.GlobalError) {
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