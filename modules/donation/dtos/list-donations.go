package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationReq struct {
	IsActive *bool  `query:"is_active" validate:"omitempty"`
	ActionBy string `reqHeader:"action-by" validate:"required"`
	utils.PaginationReq
}

type ListDonationRes struct {
	Data []entities.Donation `json:"data"`
	utils.PaginationRes
}
