package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateDonationStepReq struct {
	ActionBy    string                           `reqHeader:"action-by" validate:"required,id"`
	Description string                           `json:"description" validate:"required"`                                   //adm
	SetDate     *string                          `json:"set_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`  //adm
	Status      *entities.EnumDonationStepStatus `json:"status" validate:"omitempty,oneof=pending review done warn failed"` //adm
	utils.GetReq
}

type UpdateDonationStepRes struct {
	entities.DonationStep
}

type UpdateDonationStepOptionalFields struct {
	HasSetDate bool
	HasStatus  bool
}

func (c UpdateDonationStepReq) ValidateUpdateDonationStepOptionalFields() UpdateDonationStepOptionalFields {
	return UpdateDonationStepOptionalFields{
		HasSetDate: c.SetDate != nil,
		HasStatus:  c.Status != nil,
	}
}
