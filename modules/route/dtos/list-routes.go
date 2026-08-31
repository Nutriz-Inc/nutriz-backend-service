package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListRoutesReq struct {
	ActionBy     string                    `reqHeader:"action-by" validate:"required,id"`
	IdDriver     *string                   `query:"id_driver" validate:"omitempty,id"`
	DriverName   *string                   `query:"driver_name" validate:"omitempty,max=120"`
	Status       *entities.EnumRouteStatus `query:"status" validate:"omitempty,oneof=pending in_progress done canceled"`
	DateSet      *string                   `query:"date_set" validate:"omitempty,datetime=2006-01-02"`
	Name         *string                   `query:"name" validate:"omitempty,max=120"`
	City         *string                   `query:"city" validate:"omitempty,max=120"`
	Neighborhood *string                   `query:"neighborhood" validate:"omitempty,max=120"`
	utils.PaginationReq
}

type RouteRes struct {
	entities.Route
	DriverName *string `json:"driver_name" db:"driver_name"`
}

type ListRoutesRes struct {
	Data []RouteRes `json:"data"`
	utils.PaginationRes
}
