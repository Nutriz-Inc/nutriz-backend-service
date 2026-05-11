package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateUserBaby struct {
	userBabyRepo *repositories.UserBabyRepository
}

func HandlerCreateUserBabyStart(userBabyRepo *repositories.UserBabyRepository) *HandlerCreateUserBaby {
	return &HandlerCreateUserBaby{userBabyRepo}
}

func (h *HandlerCreateUserBaby) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateUserBabyReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateUserBaby) Execute(ctx c.Context, data *dto.CreateUserBabyReq) (*dto.CreateUserBabyRes, *fluxgo.GlobalError) {
	if data.BirthDate.After(time.Now()) {
        return nil, fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user_baby.invalid_birth_date")
    }

	userBabyId := utils.IdGenerate(utils.UserBabyEntity)

	err := h.userBabyRepo.CreateUserBaby(ctx, &repositories.CreateUserBabyRepositoryReq{
		IdUserBaby: userBabyId,
		IdUser:	data.ActionBy,
		Name:	data.Name,
		BirthDate: data.BirthDate,
	})
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create baby")
	}

	userBaby, err := h.userBabyRepo.GetUserBabybyUserBabyId(ctx, userBabyId)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user baby")
	}
	if userBaby == nil {
		return nil, fluxgo.ErrorNotFound("User baby not found")
	}

	return &dto.CreateUserBabyRes{UserBaby: *userBaby}, nil
}
