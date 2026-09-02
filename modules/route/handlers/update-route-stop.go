package handlers

import (
	c "context"
	"errors"
	"fmt"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerUpdateRouteStop struct {
	db                    *fluxgo.Database
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerUpdateRouteStopStart(
	db *fluxgo.Database,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateRouteStop {
	return &HandlerUpdateRouteStop{
		db,
		routeRepo,
		routeDonationStepRepo,
		userRepo,
	}
}

func (h *HandlerUpdateRouteStop) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateRouteStopReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateRouteStop) Execute(ctx c.Context, data *dto.UpdateRouteStopReq) (*dto.UpdateRouteStopRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanUpdateRouteStop {
		return nil, utils.ErrorForbidden("User does not have permission to update route stop", "user.forbidden")
	}

	hasError := data.HasError != nil && *data.HasError
	setDateStart := data.DateStart != nil && *data.DateStart

	if !hasError && !setDateStart {
		return nil, fluxgo.ErrorBadRequest("date_start or has_error must be sent to update", "route_stop.no_fields_to_update")
	}
	if hasError && setDateStart {
		return nil, fluxgo.ErrorBadRequest("date_start and has_error cannot be sent together", "route_stop.date_start_and_has_error")
	}

	stop, err := h.routeDonationStepRepo.GetRouteDonationStepById(ctx, data.IdStop)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stop")
	}
	if stop == nil {
		return nil, fluxgo.ErrorNotFound("Route stop not found")
	}
	if setDateStart && stop.DateStart != nil {
		return nil, fluxgo.ErrorBadRequest("Date start is already set", "route_stop.date_start_already_set")
	}

	route, err := h.routeRepo.GetRouteById(ctx, stop.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if route == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}
	if route.Status == entities.EnumRouteStatusCanceled || route.Status == entities.EnumRouteStatusDone {
		return nil, fluxgo.ErrorBadRequest("Canceled or done routes cannot be updated", "route.canceled_or_done")
	}
	if route.IdDriver != user.IdUser {
		return nil, utils.ErrorForbidden("Route belongs to another driver", "route.forbidden")
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		if hasError {
			err := h.routeDonationStepRepo.SetStatusTx(ctx, tx, data.IdStop, entities.EnumRouteDonationStepStatusError, data.ActionBy)
			if err != nil {
				return fmt.Errorf("error to update route stop: %w", err)
			}
		} else {
			err := h.routeDonationStepRepo.UpdateDateStartTx(ctx, tx, data.IdStop, data.ActionBy)
			if err != nil {
				return fmt.Errorf("error to update route stop: %w", err)
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

	updatedStop, err := h.routeDonationStepRepo.GetRouteDonationStepById(ctx, data.IdStop)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stop")
	}
	if updatedStop == nil {
		return nil, fluxgo.ErrorNotFound("Route stop not found")
	}

	return &dto.UpdateRouteStopRes{
		Stop: *updatedStop,
	}, nil
}
