package dtos

import (
	"nutriz-backend-service/shared/utils"
)

type RemoveAddressReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type RemoveAddressRes = utils.DeleteRes
