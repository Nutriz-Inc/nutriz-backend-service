package dto

import (
	"nutriz-backend-service/shared/entities"
)

type ListDonationStepTimelineReq struct {
	ActionBy       string `reqHeader:"action-by" validate:"required,id"`
	IdDonationStep string `params:"id" validate:"required,id"`
}

type ListDonationStepTimelineRes struct {
	Data []entities.DonationStepTimeline `json:"data"`
}
