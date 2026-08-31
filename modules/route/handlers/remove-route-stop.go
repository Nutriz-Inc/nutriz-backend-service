package handlers

import (
	c "context"
	"errors"
	"fmt"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerRemoveRouteStop struct {
	db                    *fluxgo.Database
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerRemoveRouteStopStart(
	db *fluxgo.Database,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerRemoveRouteStop {
	return &HandlerRemoveRouteStop{
		db,
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

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.routeDonationStepRepo.RemoveRouteDonationStepTx(ctx, tx, data.IdStop, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to remove route stop: %w", err)
		}

		if stop.StopOrder != nil {
			err = h.routeDonationStepRepo.ShiftStopOrdersAfterTx(ctx, tx, stop.IdRoute, *stop.StopOrder, data.ActionBy)
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
