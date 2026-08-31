package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type GetDonationReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type GetDonationRes struct {
	entities.Donation
	Steps   *[]entities.DonationStep `json:"steps"`
	Bottles *[]entities.Bottle       `json:"bottles"`
}
