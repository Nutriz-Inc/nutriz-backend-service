package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationPointsReq struct {
	ShowAddress bool     `query:"show_address"`
	Name        *string  `query:"name" validate:"omitempty,max=120"`
	HasHome     *bool    `query:"has_home" validate:"omitempty"`
	Longitude   *float64 `query:"longitude" validate:"omitempty,required_with=Latitude"`
	Latitude    *float64 `query:"latitude" validate:"omitempty,required_with=Longitude"`
	ZipCode     *string  `query:"zipcode" validate:"omitempty,cep"`
	utils.PaginationReq
}

type DonationPointsRes struct {
	entities.DonationPointOut
	Address         *entities.AddressOut `json:"address,omitempty"`
	DistanceFromYou *float64             `json:"distance_from_you,omitempty"`
}

type ListDonationPointsRes struct {
	Data []DonationPointsRes `json:"data"`
	utils.PaginationRes
}
