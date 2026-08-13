package dtos

import (
	dto "nutriz-backend-service/modules/donation/dtos"
	sharedDto "nutriz-backend-service/shared/dtos"
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
	sharedDto.UserOut
	DonationsCompleted *int64                   `json:"donations_completed"`
	CurrentDonation    *dto.GetDonationRes      `json:"current_donation"`
	Addresses          *[]sharedDto.AddressOut  `json:"addresses"`
	Babies             *[]sharedDto.UserBabyOut `json:"babies"`
}
