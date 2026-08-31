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

type HandlerUpdateRoute struct {
	db                    *fluxgo.Database
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerUpdateRouteStart(
	db *fluxgo.Database,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateRoute {
	return &HandlerUpdateRoute{
		db,
		routeRepo,
		routeDonationStepRepo,
		userRepo,
	}
}

func (h *HandlerUpdateRoute) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateRouteReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateRoute) Execute(ctx c.Context, data *dto.UpdateRouteReq) (*dto.UpdateRouteRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanUpdateRoute {
		return nil, utils.ErrorForbidden("User does not have permission to update route", "user.forbidden")
	}

	route, err := h.routeRepo.GetRouteById(ctx, data.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if route == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}
	if route.Status == entities.EnumRouteStatusCanceled || route.Status == entities.EnumRouteStatusDone {
		return nil, fluxgo.ErrorBadRequest("Canceled or done routes cannot be updated", "route.canceled_or_done")
	}

	repoData := &repositories.UpdateRouteRepositoryReq{
		IdRoute:   data.IdRoute,
		UpdatedBy: data.ActionBy,
	}

	var globalErr *fluxgo.GlobalError

	if user.Type == entities.EnumUserTypeAdmin {
		globalErr = h.handleAdmUpdate(data, repoData)
	} else {
		globalErr = h.handleDriverUpdate(data, route, user, repoData)
	}
	if globalErr != nil {
		return nil, globalErr
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.routeRepo.UpdateRouteTx(ctx, tx, repoData)
		if err != nil {
			return fmt.Errorf("error to update route: %w", err)
		}

		if repoData.Status != nil && *repoData.Status == entities.EnumRouteStatusCanceled {
			err = h.routeDonationStepRepo.RemoveRouteDonationStepsByIdRouteTx(ctx, tx, data.IdRoute, data.ActionBy)
			if err != nil {
				return fmt.Errorf("error to remove route donation steps: %w", err)
			}
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

	updatedRoute, err := h.routeRepo.GetRouteById(ctx, data.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get updated route")
	}
	if updatedRoute == nil {
		return nil, fluxgo.ErrorNotFound("Updated route not found")
	}

	return &dto.UpdateRouteRes{
		Route: *updatedRoute,
	}, nil
}

func (h *HandlerUpdateRoute) handleAdmUpdate(
	data *dto.UpdateRouteReq,
	repoData *repositories.UpdateRouteRepositoryReq,
) *fluxgo.GlobalError {
	validator := data.ValidateUpdateRouteOptionalFields()

	if validator.HasDriverFields() {
		return utils.ErrorForbidden(
			"Adms can only update name, city, neighborhood, status, description and date_set",
			"route.invalid_fields_for_adm",
		)
	}
	if !validator.HasAdmFields() {
		return fluxgo.ErrorBadRequest("At least one field must be sent to update", "route.no_fields_to_update")
	}

	if validator.HasStatus {
		if *data.Status != entities.EnumRouteStatusCanceled {
			return fluxgo.ErrorBadRequest(
				"Adms can only set status to canceled",
				"route.invalid_status_for_adm",
			)
		}
		if !validator.HasDescription {
			return fluxgo.ErrorBadRequest(
				"Description is required when canceling a route",
				"route.description_required",
			)
		}
	}

	if validator.HasDateSet {
		if !utils.IsFutureDate(*data.DateSet) {
			return fluxgo.ErrorBadRequest("Date set must be in the future", "date_set.invalid")
		}

		dateSet, err := utils.StringToTime(*data.DateSet)
		if err != nil {
			return fluxgo.ErrorBadRequest("Invalid date set format", "route.invalid_date_set_format")
		}

		repoData.DateSet = dateSet
	}

	repoData.Name = data.Name
	repoData.City = data.City
	repoData.Neighborhood = data.Neighborhood
	repoData.Status = data.Status
	repoData.Description = data.Description

	return nil
}

func (h *HandlerUpdateRoute) handleDriverUpdate(
	data *dto.UpdateRouteReq,
	route *entities.Route,
	user *entities.User,
	repoData *repositories.UpdateRouteRepositoryReq,
) *fluxgo.GlobalError {
	validator := data.ValidateUpdateRouteOptionalFields()

	if validator.HasAdmFields() {
		return utils.ErrorForbidden(
			"Drivers can only update date_start, date_end, mileage and user_feedback",
			"route.invalid_fields_for_driver",
		)
	}
	if !validator.HasDriverFields() {
		return fluxgo.ErrorBadRequest("At least one field must be sent to update", "route.no_fields_to_update")
	}

	if route.IdDriver != user.IdUser {
		return utils.ErrorForbidden("Route belongs to another driver", "route.forbidden")
	}

	if validator.HasDateStart && validator.HasDateEnd {
		return fluxgo.ErrorBadRequest(
			"Date start and date end cannot be sent together",
			"route.date_start_and_date_end",
		)
	}

	if validator.HasDateStart {
		if route.DateStart != nil {
			return fluxgo.ErrorBadRequest("Date start is already set", "route.date_start_already_set")
		}

		repoData.SetDateStart = true
	}

	if validator.HasDateEnd {
		if route.DateEnd != nil {
			return fluxgo.ErrorBadRequest("Date end is already set", "route.date_end_already_set")
		}
		if !validator.HasMileage || !validator.HasUserFeedback {
			return fluxgo.ErrorBadRequest(
				"Mileage and user feedback are required when setting date end",
				"route.mileage_and_user_feedback_required",
			)
		}

		repoData.SetDateEnd = true
	}

	repoData.Mileage = data.Mileage
	repoData.UserFeedback = data.UserFeedback

	return nil
}
