package dto

import (
	"nutriz-backend-service/shared/entities"
)

type UpdateRouteReq struct {
	IdRoute  string `params:"id_route" validate:"required,id"`
	ActionBy string `reqHeader:"action-by" validate:"required,id"`

	// adm fields
	Name         *string                   `json:"name" validate:"omitempty,max=150"`
	City         *string                   `json:"city" validate:"omitempty,max=100"`
	Neighborhood *string                   `json:"neighborhood" validate:"omitempty,max=100"`
	Status       *entities.EnumRouteStatus `json:"status" validate:"omitempty,oneof=canceled"`
	Description  *string                   `json:"description" validate:"omitempty,max=500"`
	DateSet      *string                   `json:"date_set" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`

	// driver fields
	DateStart    *bool    `json:"date_start" validate:"omitempty"`
	DateEnd      *bool    `json:"date_end" validate:"omitempty"`
	Mileage      *float64 `json:"mileage" validate:"omitempty,gt=0"`
	UserFeedback *string  `json:"user_feedback" validate:"omitempty,max=500"`
}

type UpdateRouteRes struct {
	entities.Route
}

type UpdateRouteOptionalFields struct {
	HasName         bool
	HasCity         bool
	HasNeighborhood bool
	HasStatus       bool
	HasDescription  bool
	HasDateSet      bool
	HasDateStart    bool
	HasDateEnd      bool
	HasMileage      bool
	HasUserFeedback bool
}

func (c UpdateRouteReq) ValidateUpdateRouteOptionalFields() UpdateRouteOptionalFields {
	return UpdateRouteOptionalFields{
		HasName:         c.Name != nil && *c.Name != "",
		HasCity:         c.City != nil && *c.City != "",
		HasNeighborhood: c.Neighborhood != nil && *c.Neighborhood != "",
		HasStatus:       c.Status != nil && *c.Status != "",
		HasDescription:  c.Description != nil && *c.Description != "",
		HasDateSet:      c.DateSet != nil && *c.DateSet != "",
		HasDateStart:    c.DateStart != nil && *c.DateStart,
		HasDateEnd:      c.DateEnd != nil && *c.DateEnd,
		HasMileage:      c.Mileage != nil,
		HasUserFeedback: c.UserFeedback != nil && *c.UserFeedback != "",
	}
}

func (f UpdateRouteOptionalFields) HasAdmFields() bool {
	return f.HasName || f.HasCity || f.HasNeighborhood || f.HasStatus || f.HasDescription || f.HasDateSet
}

func (f UpdateRouteOptionalFields) HasDriverFields() bool {
	return f.HasDateStart || f.HasDateEnd || f.HasMileage || f.HasUserFeedback
}
