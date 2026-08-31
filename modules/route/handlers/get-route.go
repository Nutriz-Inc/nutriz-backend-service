package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetRoute struct {
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerGetRouteStart(
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerGetRoute {
	return &HandlerGetRoute{
		routeRepo,
		routeDonationStepRepo,
		userRepo,
	}
}

func (h *HandlerGetRoute) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetRouteReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetRoute) Execute(ctx c.Context, data *dto.GetRouteReq) (*dto.GetRouteRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanListRoute {
		return nil, utils.ErrorForbidden("User does not have permission to get route", "user.forbidden")
	}

	route, err := h.routeRepo.GetRouteWithDriverById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if route == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}

	stopRows, err := h.routeDonationStepRepo.GetRouteDonationStepsWithAddressByIdRoute(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stops")
	}

	stops := make([]dto.RouteStop, 0, len(*stopRows))
	for _, row := range *stopRows {
		stops = append(stops, dto.RouteStop{
			RouteDonationStep: row.RouteDonationStep,
			Address:           row.Address(),
		})
	}

	return &dto.GetRouteRes{
		Route:      route.Route,
		DriverName: route.DriverName,
		Stops:      stops,
	}, nil
}
