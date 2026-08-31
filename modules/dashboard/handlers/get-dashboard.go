package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/dashboard/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetDashboard struct {
	dashboardRepo *repositories.DashboardRepository
	userRepo      *repositories.UserRepository
}

func HandlerGetDashboardStart(dashboardRepo *repositories.DashboardRepository, userRepo *repositories.UserRepository) *HandlerGetDashboard {
	return &HandlerGetDashboard{dashboardRepo, userRepo}
}

func (h *HandlerGetDashboard) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetDashboardReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetDashboard) Execute(ctx c.Context, filters *dto.GetDashboardReq) (*dto.GetDashboardRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, *filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("Only adms can access the dashboard", "dashboard.forbidden")
	}

	var startDate, endDate *time.Time
	if filters.StartDate != nil {
		startDate, err = utils.StringToDate(*filters.StartDate)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest("Invalid start date", "dashboard.invalid_start_date")
		}
	}
	if filters.EndDate != nil {
		endDate, err = utils.StringToDate(*filters.EndDate)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest("Invalid end date", "dashboard.invalid_end_date")
		}
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, fluxgo.ErrorBadRequest("End date cannot be before start date", "dashboard.invalid_date_range")
	}

	var endDateExclusive *time.Time
	if endDate != nil {
		exclusive := endDate.AddDate(0, 0, 1)
		endDateExclusive = &exclusive
	}

	totalMilkCollected, err := h.dashboardRepo.GetTotalMilkCollected(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get total milk collected")
	}

	milkByMonth, err := h.dashboardRepo.GetMilkCollectedByMonth(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get milk collected by month")
	}

	feedbackByScore, err := h.dashboardRepo.GetFeedbackByScore(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get feedback metrics")
	}

	averageServiceTimeHours, err := h.dashboardRepo.GetAverageServiceTimeHours(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get average service time")
	}

	donationsWithError, err := h.dashboardRepo.GetDonationsWithErrorCount(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donations with error count")
	}

	donorRecurrenceRate, err := h.dashboardRepo.GetDonorRecurrenceRate(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donor recurrence rate")
	}

	activeDonationsByStep, err := h.dashboardRepo.GetActiveDonationsByCurrentStep(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get active donations by step")
	}

	bottleStats, err := h.dashboardRepo.GetBottleStats(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get bottle stats")
	}
	if bottleStats == nil {
		bottleStats = &dto.BottleStats{}
	}

	routeStats, err := h.dashboardRepo.GetRouteStats(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get route stats")
	}
	if routeStats == nil {
		routeStats = &dto.RouteStats{}
	}

	return &dto.GetDashboardRes{
		TotalMilkCollected:      totalMilkCollected,
		MilkCollectedByMonth:    milkByMonth,
		FeedbackByScore:         feedbackByScore,
		AverageServiceTimeHours: averageServiceTimeHours,
		DonationsWithError:      donationsWithError,
		DonorRecurrenceRate:     donorRecurrenceRate,
		ActiveDonationsByStep:   activeDonationsByStep,
		BottleStats:             *bottleStats,
		RouteStats:              *routeStats,
	}, nil
}
