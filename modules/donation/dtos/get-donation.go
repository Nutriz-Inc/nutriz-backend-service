package dto

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/utils"
)

type GetDonationReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type GetDonationRes struct {
	sharedDto.DonationOut
	Steps *[]sharedDto.DonationStepOut `json:"steps"`
}
