package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/utils"
)

type GetAddressReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type GetAddressRes struct {
	sharedDto.AddressOut
}
