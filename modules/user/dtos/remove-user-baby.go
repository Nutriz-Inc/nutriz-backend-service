package dtos

import "nutriz-backend-service/shared/utils"

type RemoveUserBabyReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type RemoveUserBabyRes = utils.DeleteRes
