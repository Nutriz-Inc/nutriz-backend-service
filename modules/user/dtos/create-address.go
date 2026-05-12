package dtos

import (
	"nutriz-backend-service/shared/entities"

	"github.com/paemuri/brdoc/v2"
)

type CreateAddressReq struct {
	ActionBy     string               `reqHeader:"action-by" validate:"required,id"`
	ZipCode      string               `json:"zip_code" validate:"required,cep"`
	Street       string               `json:"street" validate:"required,max=150"`
	Number       string               `json:"number" validate:"required,max=10"`
	City         string               `json:"city" validate:"required,max=100"`
	State        brdoc.FederativeUnit `json:"state" validate:"required,max=27"`
	Neighborhood string               `json:"neighborhood" validate:"required,max=100"`
	Complement   *string              `json:"complement" validate:"omitempty,max=150"`
}

type CreateAddressRes struct {
	entities.Address
}
