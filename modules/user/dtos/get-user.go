package dtos

import (
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type GetUserReq struct {
	ShowAddress            bool   `query:"show_address"`
	ShowBaby               bool   `query:"show_baby"`
	ShowDonationsCompleted bool   `query:"show_donations_completed"`
	ShowCurrentDonation    bool   `query:"show_current_donation"`
	ActionBy               string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type GetUserRes struct {
	entities.User
	DonationsCompleted *int64               `json:"donations_completed"`
	CurrentDonation    *dto.GetDonationRes  `json:"current_donation"`
	Addresses          *[]entities.Address  `json:"addresses"`
	Babies             *[]entities.UserBaby `json:"babies"`
}
