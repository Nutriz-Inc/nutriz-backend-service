package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type GetUserReq struct {
	ShowAddress bool   `query:"show_address"`
	ShowBaby    bool   `query:"show_baby"`
	ActionBy    string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type GetUserRes struct {
	entities.User
	Addresses *[]entities.Address  `json:"addresses"`
	Babies    *[]entities.UserBaby `json:"babies"`
}
