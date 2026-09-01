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

	validator := data.ValidateCreateRouteOptionalFields()

	var donationSteps []repositories.DonationStepWithLocation
	var stopOrders []int16

	if validator.HasStops {
		steps, globalErr := h.getDonationSteps(ctx, data)
		if globalErr != nil {
			return nil, globalErr
		}
		donationSteps = *steps

		stopOrders, globalErr = utils.OptimizeStops(ctx, h.config, stopCoordinates(steps))
		if globalErr != nil {
			return nil, globalErr
		}
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

		if validator.HasStops {
			for index, donationStep := range donationSteps {
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

	stops := make([]entities.RouteDonationStep, 0)

	if validator.HasStops {
		routeStops, err := h.routeDonationStepRepo.GetRouteDonationStepsByIdRoute(ctx, idRoute)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get route donation steps")
		}
		if routeStops == nil || len(*routeStops) == 0 {
			return nil, fluxgo.ErrorNotFound("No donation steps found for the route")
		}
		stops = *routeStops
	}

	return &dto.CreateRouteRes{
		Route: *route,
		Stops: stops,
	}, nil
}

func (h *HandlerCreateRoute) getDonationSteps(
	ctx c.Context,
	data *dto.CreateRouteReq,
) (*[]repositories.DonationStepWithLocation, *fluxgo.GlobalError) {
	ids := make([]string, 0, len(*data.Stops))
	seen := make(map[string]bool, len(*data.Stops))

	for _, id := range *data.Stops {
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

		if data.City != nil && !utils.MatchesAddressField(donationStep.City, *data.City) {
			return nil, fluxgo.ErrorBadRequest(
				fmt.Sprintf("Donation step %s is not in the city %s", donationStep.IdDonationStep, *data.City),
				"stops.invalid_city",
			)
		}
		if data.Neighborhood != nil && !utils.MatchesAddressField(donationStep.Neighborhood, *data.Neighborhood) {
			return nil, fluxgo.ErrorBadRequest(
				fmt.Sprintf("Donation step %s is not in the neighborhood %s", donationStep.IdDonationStep, *data.Neighborhood),
				"stops.invalid_neighborhood",
			)
		}

		ordered = append(ordered, donationStep)
	}

	return &ordered, nil
}

func stopCoordinates(donationSteps *[]repositories.DonationStepWithLocation) []utils.StopCoordinates {
	stops := make([]utils.StopCoordinates, 0, len(*donationSteps))

	for _, donationStep := range *donationSteps {
		stops = append(stops, utils.StopCoordinates{
			Latitude:  donationStep.Latitude,
			Longitude: donationStep.Longitude,
		})
	}

	return stops
}
