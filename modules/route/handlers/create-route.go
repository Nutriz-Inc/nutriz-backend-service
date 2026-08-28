package handlers

import (
	c "context"
	"errors"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/provider/location"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerCreateRoute struct {
	db                    *fluxgo.Database
	config                *config.Env
	routeRepo             *repositories.RouteRepository
	routeDonationStepRepo *repositories.RouteDonationStepRepository
	donationStepRepo      *repositories.DonationStepRepository
	userRepo              *repositories.UserRepository
}

func HandlerCreateRouteStart(
	db *fluxgo.Database,
	config *config.Env,
	routeRepo *repositories.RouteRepository,
	routeDonationStepRepo *repositories.RouteDonationStepRepository,
	donationStepRepo *repositories.DonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateRoute {
	return &HandlerCreateRoute{
		db,
		config,
		routeRepo,
		routeDonationStepRepo,
		donationStepRepo,
		userRepo,
	}
}

func (h *HandlerCreateRoute) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateRouteReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateRoute) Execute(ctx c.Context, data *dto.CreateRouteReq) (*dto.CreateRouteRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanCreateRoute {
		return nil, utils.ErrorForbidden("User does not have permission to create route", "user.forbidden")
	}

	driver, err := h.userRepo.GetUserById(ctx, data.IdDriver)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get driver")
	}
	if driver == nil {
		return nil, fluxgo.ErrorNotFound("Driver not found")
	}
	if driver.Type != entities.EnumUserTypeDriver {
		return nil, fluxgo.ErrorBadRequest("User is not a driver", "driver.invalid_type")
	}

	if !utils.IsFutureDate(data.DateSet) {
		return nil, fluxgo.ErrorBadRequest("Date set must be in the future", "date_set.invalid")
	}
	dateSet, err := utils.StringToTime(data.DateSet)
	if err != nil {
		return nil, fluxgo.ErrorBadRequest("Invalid date set format", "route.invalid_date_set_format")
	}

	donationSteps, globalErr := h.getDonationSteps(ctx, data.Stops)
	if globalErr != nil {
		return nil, globalErr
	}

	stopOrders, globalErr := h.getStopOrders(ctx, donationSteps)
	if globalErr != nil {
		return nil, globalErr
	}

	idRoute := utils.IdGenerate(utils.RouteEntity)

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.routeRepo.CreateRouteTx(ctx, tx, &repositories.CreateRouteRepositoryReq{
			IdRoute:      idRoute,
			IdDriver:     data.IdDriver,
			IdUser:       data.ActionBy,
			Name:         data.Name,
			Description:  data.Description,
			City:         data.City,
			Neighborhood: data.Neighborhood,
			Status:       entities.EnumRouteStatusPending, // Always starts as pending
			DateSet:      *dateSet,
		})
		if err != nil {
			return fmt.Errorf("error to create route: %w", err)
		}

		for index, donationStep := range *donationSteps {
			err = h.routeDonationStepRepo.CreateRouteDonationStepTx(ctx, tx, &repositories.CreateRouteDonationStepRepositoryReq{
				IdRouteDonationStep: utils.IdGenerate(utils.RouteDonationStepEntity),
				IdRoute:             idRoute,
				IdDonationStep:      donationStep.IdDonationStep,
				IdUser:              data.ActionBy,
				StopOrder:           stopOrders[index],
			})
			if err != nil {
				return fmt.Errorf("error to create route donation step: %w", err)
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

	route, err := h.routeRepo.GetRouteById(ctx, idRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route")
	}
	if route == nil {
		return nil, fluxgo.ErrorNotFound("Route not found")
	}

	stops, err := h.routeDonationStepRepo.GetRouteDonationStepsByIdRoute(ctx, idRoute)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route donation steps")
	}
	if stops == nil || len(*stops) == 0 {
		return nil, fluxgo.ErrorNotFound("No donation steps found for the route")
	}

	return &dto.CreateRouteRes{
		Route: *route,
		Stops: *stops,
	}, nil
}

func (h *HandlerCreateRoute) getDonationSteps(
	ctx c.Context,
	stops []string,
) (*[]repositories.DonationStepWithLocation, *fluxgo.GlobalError) {
	ids := make([]string, 0, len(stops))
	seen := make(map[string]bool, len(stops))

	for _, id := range stops {
		if utils.GetIdEntity(id) != utils.DonationStepEntity {
			return nil, fluxgo.ErrorBadRequest("Stop is not a donation step id", "stops.invalid_id")
		}
		if seen[id] {
			return nil, fluxgo.ErrorBadRequest("Duplicated donation step on stops", "stops.duplicated")
		}

		seen[id] = true
		ids = append(ids, id)
	}

	donationSteps, err := h.donationStepRepo.GetDonationStepsWithLocationByIds(ctx, ids)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation steps")
	}

	byId := make(map[string]repositories.DonationStepWithLocation, len(*donationSteps))
	for _, donationStep := range *donationSteps {
		byId[donationStep.IdDonationStep] = donationStep
	}

	ordered := make([]repositories.DonationStepWithLocation, 0, len(ids))

	for _, id := range ids {
		donationStep, ok := byId[id]
		if !ok {
			return nil, fluxgo.ErrorNotFound("Donation step not found")
		}
		if !donationStep.IsDonationActive {
			return nil, fluxgo.ErrorBadRequest(fmt.Sprintf("Donation of %s is not active", donationStep.IdDonationStep), "donation.inactive")
		}
		if donationStep.Latitude == nil || donationStep.Longitude == nil {
			return nil, fluxgo.ErrorBadRequest(
				"Donation step does not have an address with coordinates",
				"donation_step.missing_coordinates",
			)
		}

		ordered = append(ordered, donationStep)
	}

	return &ordered, nil
}

func (h *HandlerCreateRoute) getStopOrders(
	ctx c.Context,
	donationSteps *[]repositories.DonationStepWithLocation,
) ([]int16, *fluxgo.GlobalError) {
	coordinates := make([]location.Coordinate, 0, len(*donationSteps))

	for _, donationStep := range *donationSteps {
		coordinates = append(coordinates, location.Coordinate{
			Latitude:  *donationStep.Latitude,
			Longitude: *donationStep.Longitude,
		})
	}

	order, err := utils.GetOptimizedStopOrder(ctx, coordinates, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to build the route: " + err.Error())
	}

	stopOrders := make([]int16, 0, len(order))
	for _, position := range order {
		stopOrders = append(stopOrders, int16(position))
	}

	return stopOrders, nil
}
