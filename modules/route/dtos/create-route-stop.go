package dto

import (
	"nutriz-backend-service/shared/entities"
)

type CreateRouteStopReq struct {
	IdRoute        string `params:"id_route" validate:"required,id"`
	ActionBy       string `reqHeader:"action-by" validate:"required,id"`
	IdDonationStep string `json:"id_donation_step" validate:"required,id"`
}

type CreateRouteStopRes struct {
	entities.Route
	Stops []entities.RouteDonationStep `json:"stops"`
}
