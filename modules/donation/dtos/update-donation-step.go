package dto

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateDonationStepReq struct {
	ActionBy    string                           `reqHeader:"action-by" validate:"required,id"`
	Description string                           `json:"description" validate:"required"`                                   //adm
	SetDate     *string                          `json:"set_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`  //adm
	Status      *entities.EnumDonationStepStatus `json:"status" validate:"omitempty,oneof=pending review done warn failed"` //adm
	IdAddress   *string                          `json:"id_address" validate:"omitempty,id"`                                //adm
	Address     *AddressCreateBase               `json:"address" validate:"omitempty"`                                      //adm
	utils.GetReq
}

type UpdateDonationStepRes struct {
	sharedDto.DonationStepOut
}

type UpdateDonationStepOptionalFields struct {
	HasSetDate bool
	HasStatus  bool
	HasAddress bool
}

func (c UpdateDonationStepReq) ValidateUpdateDonationStepOptionalFields() UpdateDonationStepOptionalFields {
	return UpdateDonationStepOptionalFields{
		HasSetDate: c.SetDate != nil,
		HasStatus:  c.Status != nil,
		HasAddress: c.IdAddress != nil || c.Address != nil,
	}
}
