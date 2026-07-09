package dto

import (
	"nutriz-backend-service/shared/entities"
)

type CreateDonationStepReq struct {
	ActionBy    string                     `reqHeader:"action-by" validate:"required,id"`
	IdDonation  string                     `json:"id_donation" validate:"required,id"`
	IdAddress   *string                    `json:"id_address" validate:"omitempty,id"`
	Address     *AddressCreateBase         `json:"address" validate:"omitempty"`
	Name        entities.EnumDonationSteps `json:"name" validate:"required,oneof='Exame de sangue' 'Entregar kit de ordenha' 'Coletar leite' 'Análise de leite'"`
	Description string                     `json:"description" validate:"required"`
	SetDate     *string                    `json:"set_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type AddressCreateBase struct {
	ZipCode    string  `json:"zip_code" validate:"required,cep"`
	Number     *string `json:"number" validate:"omitempty,max=10"`
	Complement *string `json:"complement" validate:"omitempty,max=150"`
}

func (c CreateDonationStepReq) HasAddress() bool {
	return c.IdAddress != nil || c.Address != nil
}

type CreateDonationStepRes struct {
	entities.DonationStep
}
