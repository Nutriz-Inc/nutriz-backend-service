package dto

import (
	"nutriz-backend-service/shared/entities"
)

type CreateDonationReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
}

type CreateDonationRes struct {
	entities.Donation
}
