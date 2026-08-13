package dtos

import (
	"nutriz-backend-service/shared/entities"
)

type CreateAddressReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	AddressCreateBase
}

type AddressCreateBase struct {
	ZipCode    string  `json:"zip_code" validate:"required,cep"`
	Number     *string `json:"number" validate:"omitempty,max=10"`
	Complement *string `json:"complement" validate:"omitempty,max=150"`
}

type CreateAddressRes struct {
	entities.AddressOut
}
