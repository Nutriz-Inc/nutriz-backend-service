package dto

import (
	"nutriz-backend-service/shared/entities"
)

type CreateRouteReq struct {
	ActionBy     string    `reqHeader:"action-by" validate:"required,id"`
	IdDriver     string    `json:"id_driver" validate:"required,id"`
	DateSet      string    `json:"date_set" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Stops        *[]string `json:"stops" validate:"omitempty,dive,required,id"`
	Name         string    `json:"name" validate:"required,max=150"`
	Description  string    `json:"description" validate:"required"`
	Neighborhood *string   `json:"neighborhood" validate:"omitempty,max=100"`
	City         *string   `json:"city" validate:"omitempty,max=100"`
}

type CreateRouteRes struct {
	entities.Route
	Stops []entities.RouteDonationStep `json:"stops"`
}

type CreateRouteOptionalFields struct {
	HasStops        bool
	HasCity         bool
	HasNeighborhood bool
}

func (c CreateRouteReq) ValidateCreateRouteOptionalFields() CreateRouteOptionalFields {
	return CreateRouteOptionalFields{
		HasStops:        c.Stops != nil && len(*c.Stops) > 0,
		HasCity:         c.City != nil && *c.City != "",
		HasNeighborhood: c.Neighborhood != nil && *c.Neighborhood != "",
	}
}
