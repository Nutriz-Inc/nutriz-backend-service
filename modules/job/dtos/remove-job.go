package dtos

import (
	"nutriz-backend-service/shared/utils"
)

type RemoveJobReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type RemoveJobRes = utils.DeleteRes
