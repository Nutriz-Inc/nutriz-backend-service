package dto

import "nutriz-backend-service/shared/entities"

type GetDashboardReq struct {
	StartDate *string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   *string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	ActionBy  *string `reqHeader:"action-by" validate:"required,id"`
}

type MilkCollectedByMonth struct {
	Month string  `json:"month" db:"month"`
	Total float64 `json:"total" db:"total"`
}

type FeedbackScoreCount struct {
	Score int16 `json:"score" db:"score"`
	Count int64 `json:"count" db:"count"`
}

type ActiveDonationsByStep struct {
	Step       entities.EnumDonationSteps `json:"step"`
	Count      int64                      `json:"count"`
	Percentage float64                    `json:"percentage"`
}

type BottleStats struct {
	BottlesCount           int64   `json:"bottles_count" db:"bottles_count"`
	DiscardedBottlesCount  int64   `json:"discarded_bottles_count" db:"discarded_bottles_count"`
	AverageBottlesPerDonor float64 `json:"average_bottles_per_donor" db:"average_bottles_per_donor"`
	BottlesUtilizationRate float64 `json:"bottles_utilization_rate" db:"bottles_utilization_rate"`
}

type RouteStats struct {
	AverageMileagePerRoute    *float64 `json:"average_mileage_per_route" db:"average_mileage_per_route"`
	AverageStopsPerRoute      *float64 `json:"average_stops_per_route" db:"average_stops_per_route"`
	AverageRouteDurationHours *float64 `json:"average_route_duration_hours" db:"average_route_duration_hours"`
}

type GetDashboardRes struct {
	TotalMilkCollected      float64                 `json:"total_milk_collected"`
	MilkCollectedByMonth    []MilkCollectedByMonth  `json:"milk_collected_by_month"`
	FeedbackByScore         []FeedbackScoreCount    `json:"feedback_by_score"`
	AverageServiceTimeHours *float64                `json:"average_service_time_hours"`
	DonationsWithError      int64                   `json:"donations_with_error"`
	DonorRecurrenceRate     float64                 `json:"donor_recurrence_rate"`
	ActiveDonationsByStep   []ActiveDonationsByStep `json:"active_donations_by_step"`
	BottleStats
	RouteStats
}
