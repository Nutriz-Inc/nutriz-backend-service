package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationStepTimelineReq struct {
	ActionBy       string `reqHeader:"action-by" validate:"required,id"`
	IdDonationStep string `params:"id" validate:"required,id"`
	utils.PaginationReq
}

type ListDonationStepTimelineRes struct {
	Data []entities.DonationStepTimeline `json:"data"`
	utils.PaginationRes
}
