package handlers

import (
	c "context"
	"errors"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerCreateRouteStop struct {
	db                    *fluxgo.Database
	config                *config.Env
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	donationStepRepo      *repositories.DonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerCreateRouteStopStart(
	db *fluxgo.Database,
	config *config.Env,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	donationStepRepo *repositories.DonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateRouteStop {
	return &HandlerCreateRouteStop{
		db,
		config,
		routeRepo,
		routeDonationStepRepo,
		donationStepRepo,
		userRepo,
	}
}

func (h *HandlerCreateRouteStop) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateRouteStopReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateRouteStop) Execute(ctx c.Context, data *dto.CreateRouteStopReq) (*dto.CreateRouteStopRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanCreateRouteStop {
		return nil, utils.ErrorForbidden("User does not have permission to create route stop", "user.forbidden")
	}

	route, err := h.routeRepo.GetRouteById(ctx, data.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if route == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}
	if route.Status == entities.EnumRouteStatusCanceled {
		return nil, fluxgo.ErrorBadRequest("Canceled routes cannot be updated", "route.canceled")
	}

	currentStops, err := h.routeDonationStepRepo.GetRouteDonationStepsWithLocationByIdRoute(ctx, data.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stops")
	}

	for _, stop := range *currentStops {
		if stop.IdDonationStep == data.IdDonationStep {
			return nil, fluxgo.ErrorBadRequest("Donation step is already a stop of the route", "stops.duplicated")
		}
	}

	donationStep, globalErr := h.getDonationStep(ctx, data, route)
	if globalErr != nil {
		return nil, globalErr
	}

	stops := append(routeStopCoordinates(*currentStops), utils.StopCoordinates{
		Latitude:  donationStep.Latitude,
		Longitude: donationStep.Longitude,
	})

	stopOrders, estimatedTime, globalErr := utils.OptimizeStops(ctx, h.config, stops)
	if globalErr != nil {
		return nil, globalErr
	}

	idRouteDonationStep := utils.IdGenerate(utils.RouteDonationStepEntity)

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		for index, stop := range *currentStops {
			err := h.routeDonationStepRepo.UpdateStopOrderTx(
				ctx,
				tx,
				stop.IdRouteDonationStep,
				stopOrders[index],
				data.ActionBy,
			)
			if err != nil {
				return fmt.Errorf("error to reorder route stops: %w", err)
			}
		}

		err := h.routeDonationStepRepo.CreateRouteDonationStepTx(ctx, tx, &repositories.CreateRouteDonationStepRepositoryReq{
			IdRouteDonationStep: idRouteDonationStep,
			IdRoute:             data.IdRoute,
			IdDonationStep:      data.IdDonationStep,
			IdUser:              data.ActionBy,
			StopOrder:           stopOrders[len(stopOrders)-1],
		})
		if err != nil {
			return fmt.Errorf("error to create route stop: %w", err)
		}

		err = h.routeRepo.UpdateEstimatedTimeTx(ctx, tx, data.IdRoute, &estimatedTime, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to update route estimated time: %w", err)
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
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if updatedRoute == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}

	updatedStops, err := h.routeDonationStepRepo.GetRouteDonationStepsByIdRoute(ctx, data.IdRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stops")
	}

	return &dto.CreateRouteStopRes{
		Route: *updatedRoute,
		Stops: *updatedStops,
	}, nil
}

func (h *HandlerCreateRouteStop) getDonationStep(
	ctx c.Context,
	data *dto.CreateRouteStopReq,
	route *entities.Route,
) (*repositories.DonationStepWithLocation, *fluxgo.GlobalError) {
	if utils.GetIdEntity(data.IdDonationStep) != utils.DonationStepEntity {
		return nil, fluxgo.ErrorBadRequest("Stop is not a donation step id", "stops.invalid_id")
	}

	donationSteps, err := h.donationStepRepo.GetDonationStepsWithLocationByIds(ctx, []string{data.IdDonationStep})
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if len(*donationSteps) == 0 {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	donationStep := (*donationSteps)[0]

	if !donationStep.IsDonationActive {
		return nil, fluxgo.ErrorBadRequest(
			fmt.Sprintf("Donation of %s is not active", donationStep.IdDonationStep),
			"donation.inactive",
		)
	}
	if !donationStep.HasAddress {
		return nil, fluxgo.ErrorBadRequest(
			fmt.Sprintf("Donation step %s has no address", donationStep.IdDonationStep),
			"stops.no_address",
		)
	}
	if donationStep.InActiveRoute {
		return nil, fluxgo.ErrorBadRequest(
			fmt.Sprintf("Donation step %s is already in another active route", donationStep.IdDonationStep),
			"stops.already_in_route",
		)
	}

	if route.City != nil && !utils.MatchesAddressField(donationStep.City, *route.City) {
		return nil, fluxgo.ErrorBadRequest(
			fmt.Sprintf("Donation step %s is not in the city %s", donationStep.IdDonationStep, *route.City),
			"stops.invalid_city",
		)
	}
	if route.Neighborhood != nil && !utils.MatchesAddressField(donationStep.Neighborhood, *route.Neighborhood) {
		return nil, fluxgo.ErrorBadRequest(
			fmt.Sprintf("Donation step %s is not in the neighborhood %s", donationStep.IdDonationStep, *route.Neighborhood),
			"stops.invalid_neighborhood",
		)
	}

	return &donationStep, nil
}
