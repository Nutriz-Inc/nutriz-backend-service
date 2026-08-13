package dto

import (
	sharedDto "nutriz-backend-service/shared/dtos"
)

type ListDonationStepTimelineReq struct {
	ActionBy       string `reqHeader:"action-by" validate:"required,id"`
	IdDonationStep string `params:"id" validate:"required,id"`
}

type ListDonationStepTimelineRes struct {
	Data []sharedDto.DonationStepTimelineOut `json:"data"`
}
