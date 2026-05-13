package dto

import (
	"nutriz-backend-service/shared/entities"
)

type CreateDonationStepReq struct {
	ActionBy    string                     `reqHeader:"action-by" validate:"required,id"`
	IdDonation  string                     `json:"id_donation" validate:"required,id"`
	Name        entities.EnumDonationSteps `json:"name" validate:"required,oneof='Exame de sangue' 'Entregar kit de ordenha' 'Coletar leite' 'Análise de leite'"`
	Description string                     `json:"description" validate:"required"`
	SetDate     *string                    `json:"set_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type CreateDonationStepRes struct {
	entities.DonationStep
}
