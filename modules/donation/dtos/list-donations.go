package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationReq struct {
	IsActive     *bool   `query:"is_active" validate:"omitempty"`
	UserDocument *string `query:"user_document" validate:"omitempty,document"`
	ActionBy     *string `reqHeader:"action-by" validate:"required,id"`
	utils.PaginationReq
}

type ListDonationRes struct {
	Data []entities.Donation `json:"data"`
	utils.PaginationRes
}
