package dtos

import (
	"nutriz-backend-service/shared/utils"
)

type RemoveUserReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type RemoveUserRes = utils.DeleteRes
