package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationStepsReq struct {
	ActionBy     string                           `reqHeader:"action-by" validate:"required,id"`
	Status       *entities.EnumDonationStepStatus `query:"status" validate:"omitempty,oneof=pending review done warn failed"`
	IdDonation   *string                          `query:"id_donation" validate:"omitempty,id"`
	Name         *entities.EnumDonationSteps      `query:"name" validate:"omitempty,oneof='Exame de sangue' 'Entregar kit de ordenha' 'Coletar leite' 'Análise de leite'"`
	SetDate      *string                          `query:"set_date" validate:"omitempty,datetime=2006-01-02"`
	Neighborhood *string                          `query:"neighborhood" validate:"omitempty,max=120"`
	City         *string                          `query:"city" validate:"omitempty,max=120"`
	HasAddress   *bool                            `query:"has_address" validate:"omitempty"`
	utils.PaginationReq
}

type DonationStepRes struct {
	entities.DonationStep
	Address *entities.Address `json:"address"`
}

type ListDonationStepsRes struct {
	Data []DonationStepRes `json:"data"`
	utils.PaginationRes
}
