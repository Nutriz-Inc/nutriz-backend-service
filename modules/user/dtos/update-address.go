package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateAddressReq struct {
	ActionBy   string  `reqHeader:"action-by" validate:"required,id"`
	ZipCode    *string `json:"zip_code" validate:"omitempty,cep"`
	Number     *string `json:"number" validate:"omitempty,max=10"`
	Complement *string `json:"complement" validate:"omitempty,max=150"`
	utils.GetReq
}

type UpdateAddressRes struct {
	entities.Address
}

type UpdateAddressOptionalFields struct {
	HasZipCode    bool
	HasNumber     bool
	HasComplement bool
}

func (c UpdateAddressReq) ValidateUpdateAddressOptionalFields() UpdateAddressOptionalFields {
	return UpdateAddressOptionalFields{
		HasZipCode:    c.ZipCode != nil && *c.ZipCode != "",
		HasNumber:     c.Number != nil && *c.Number != "",
		HasComplement: c.Complement != nil && *c.Complement != "",
	}
}
