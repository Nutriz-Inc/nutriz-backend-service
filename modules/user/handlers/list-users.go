package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListUsers struct {
	userRepo *repositories.UserRepository
}

func HandlerListUsersStart(userRepo *repositories.UserRepository) *HandlerListUsers {
	return &HandlerListUsers{userRepo}
}

func (h *HandlerListUsers) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListUsersReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListUsers) Execute(ctx c.Context, filters *dto.ListUsersReq) (*dto.ListUsersRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanListUsers {
		return nil, utils.ErrorForbidden("User does not have permission to list users", "user.forbidden")
	}

	users, total, err := h.userRepo.ListUsersByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list users")
	}

	return &dto.ListUsersRes{
		Data: *users,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
