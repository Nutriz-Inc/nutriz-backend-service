package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListDonationReq struct {
	IsActive     *bool                       `query:"is_active" validate:"omitempty"`
	UserDocument *string                     `query:"user_document" validate:"omitempty,document"`
	UserName     *string                     `query:"user_name" validate:"omitempty"`
	CurrentStep  *entities.EnumDonationSteps `query:"current_step" validate:"omitempty,oneof='Exame de sangue' 'Entregar kit de ordenha' 'Coletar leite' 'Análise de leite'"`
	ActionBy     *string                     `reqHeader:"action-by" validate:"required,id"`
	utils.PaginationReq
}

type DonationRes struct {
	entities.Donation
	UserName     *string                    `json:"user_name,omitempty" db:"user_name"`
	UserDocument *string                    `json:"user_document,omitempty" db:"user_document"`
	CurrentStep  entities.EnumDonationSteps `json:"current_step" db:"current_step"`
}

type ListDonationRes struct {
	Data []DonationRes `json:"data"`
	utils.PaginationRes
}
