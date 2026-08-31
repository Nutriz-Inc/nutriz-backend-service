package handlers

import (
	c "context"
	"errors"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/provider/location"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerRemoveRouteStop struct {
	db                    *fluxgo.Database
	config                *config.Env
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerRemoveRouteStopStart(
	db *fluxgo.Database,
	config *config.Env,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerRemoveRouteStop {
	return &HandlerRemoveRouteStop{
		db,
		config,
		routeRepo,
		routeDonationStepRepo,
		userRepo,
	}
}

func (h *HandlerRemoveRouteStop) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveRouteStopReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveRouteStop) Execute(ctx c.Context, data *dto.RemoveRouteStopReq) (*dto.RemoveRouteStopRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanRemoveRouteStop {
		return nil, utils.ErrorForbidden("User does not have permission to remove route stop", "user.forbidden")
	}

	stop, err := h.routeDonationStepRepo.GetRouteDonationStepById(ctx, data.IdStop)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stop")
	}
	if stop == nil {
		return nil, fluxgo.ErrorNotFound("Route stop not found")
	}

	remainingStops, globalErr := h.getRemainingStops(ctx, stop.IdRoute, data.IdStop)
	if globalErr != nil {
		return nil, globalErr
	}

	stopOrders, globalErr := h.getStopOrders(ctx, remainingStops)
	if globalErr != nil {
		return nil, globalErr
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.routeDonationStepRepo.RemoveRouteDonationStepTx(ctx, tx, data.IdStop, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to remove route stop: %w", err)
		}

		for index, remaining := range remainingStops {
			err = h.routeDonationStepRepo.UpdateStopOrderTx(
				ctx,
				tx,
				remaining.IdRouteDonationStep,
				stopOrders[index],
				data.ActionBy,
			)
			if err != nil {
				return fmt.Errorf("error to reorder route stops: %w", err)
			}
		}

		err = h.routeRepo.TouchRouteTx(ctx, tx, stop.IdRoute, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to update route: %w", err)
		}

		return nil
	})
	if err != nil {
		var txErr *utils.TxError
		if errors.As(err, &txErr) {
			return nil, txErr.Err
		}
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	removedStop, err := h.routeDonationStepRepo.GetRouteDonationStepById(ctx, data.IdStop)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stop")
	}

	return &dto.RemoveRouteStopRes{
		Success: removedStop == nil,
	}, nil
}

func (h *HandlerRemoveRouteStop) getRemainingStops(
	ctx c.Context,
	idRoute string,
	idRemovedStop string,
) ([]repositories.RouteDonationStepWithLocation, *fluxgo.GlobalError) {
	stops, err := h.routeDonationStepRepo.GetRouteDonationStepsWithLocationByIdRoute(ctx, idRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stops")
	}

	remaining := make([]repositories.RouteDonationStepWithLocation, 0, len(*stops))
	for _, stop := range *stops {
		if stop.IdRouteDonationStep == idRemovedStop {
			continue
		}

		remaining = append(remaining, stop)
	}

	return remaining, nil
}

func (h *HandlerRemoveRouteStop) getStopOrders(
	ctx c.Context,
	stops []repositories.RouteDonationStepWithLocation,
) ([]int16, *fluxgo.GlobalError) {
	if len(stops) == 0 {
		return nil, nil
	}

	coordinates := make([]location.Coordinate, 0, len(stops))
	for _, stop := range stops {
		latitude, longitude := utils.FillMissingCoordinates(stop.Latitude, stop.Longitude)

		coordinates = append(coordinates, location.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		})
	}

	optimizedRoute, err := utils.GetOptimizedRoute(ctx, coordinates, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to reorder the route: " + err.Error())
	}

	stopOrders := make([]int16, 0, len(optimizedRoute.StopOrders))
	for _, position := range optimizedRoute.StopOrders {
		stopOrders = append(stopOrders, int16(position))
	}

	return stopOrders, nil
}
