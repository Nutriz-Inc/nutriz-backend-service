package dto

import (
	"nutriz-backend-service/shared/entities"
)

type UpdateRouteStopReq struct {
	IdStop   string `params:"id_stop" validate:"required,id"`
	ActionBy string `reqHeader:"action-by" validate:"required,id"`

	DateStart *bool `json:"date_start" validate:"omitempty"`
}

type UpdateRouteStopRes struct {
	Stop entities.RouteDonationStep `json:"stop"`
}
