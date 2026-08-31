package dto

import (
	"nutriz-backend-service/shared/utils"
)

type RemoveRouteStopReq struct {
	IdStop   string `params:"id_stop" validate:"required,id"`
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
}

type RemoveRouteStopRes = utils.DeleteRes
