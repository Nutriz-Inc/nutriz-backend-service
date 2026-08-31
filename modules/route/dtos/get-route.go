package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type GetRouteReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type RouteStop struct {
	entities.RouteDonationStep
	Address *entities.Address `json:"address"`
}

type GetRouteRes struct {
	entities.Route
	DriverName *string     `json:"driver_name"`
	Stops      []RouteStop `json:"stops"`
}
