package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListRoutes struct {
	routeRepo *repositories.RouteRepository
	userRepo  *repositories.UserRepository
}

func HandlerListRoutesStart(
	routeRepo *repositories.RouteRepository,
	userRepo *repositories.UserRepository,
) *HandlerListRoutes {
	return &HandlerListRoutes{
		routeRepo,
		userRepo,
	}
}

func (h *HandlerListRoutes) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListRoutesReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListRoutes) Execute(ctx c.Context, filters *dto.ListRoutesReq) (*dto.ListRoutesRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanListRoute {
		return nil, utils.ErrorForbidden("User does not have permission to list routes", "user.forbidden")
	}

	var idNurse *string
	if user.Type == entities.EnumUserTypeNurse {
		idNurse = &user.IdUser
	}

	routes, total, err := h.routeRepo.ListRoutesByFilters(ctx, filters, idNurse)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list routes")
	}

	return &dto.ListRoutesRes{
		Data: *routes,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
