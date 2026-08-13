package dto

import (
	sharedDto "nutriz-backend-service/shared/dtos"
)

type CreateDonationReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
}

type CreateDonationRes struct {
	sharedDto.DonationOut
}
