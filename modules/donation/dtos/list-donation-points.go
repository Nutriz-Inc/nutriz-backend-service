package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationPointsReq struct {
	Name    *string `query:"name" validate:"omitempty,max=120"`
	Cnpj    *string `query:"cnpj" validate:"omitempty,max=14,document"`
	HasHome *bool   `query:"has_home" validate:"omitempty"`
	//to do: add location filter
	utils.PaginationReq
}

type ListDonationPointsRes struct {
	Data []entities.DonationPoint `json:"data"`
	utils.PaginationRes
}
