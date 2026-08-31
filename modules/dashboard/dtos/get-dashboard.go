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

type GetDashboardRes struct {
	TotalMilkCollected      float64                 `json:"total_milk_collected"`
	MilkCollectedByMonth    []MilkCollectedByMonth  `json:"milk_collected_by_month"`
	FeedbackByScore         []FeedbackScoreCount    `json:"feedback_by_score"`
	AverageServiceTimeHours *float64                `json:"average_service_time_hours"`
	DonationsWithError      int64                   `json:"donations_with_error"`
	DonorRecurrenceRate     float64                 `json:"donor_recurrence_rate"`
	ActiveDonationsByStep   []ActiveDonationsByStep `json:"active_donations_by_step"`
}
