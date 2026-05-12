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
	users, total, err := h.userRepo.ListUsersByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list users")
	}
	if users == nil || len(*users) == 0 {
		return nil, fluxgo.ErrorNotFound("Users not found")
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
