package dtos

import (
	"nutriz-backend-service/shared/entities"
)

type CreateUserReq struct {
	Type        entities.EnumUserType `json:"type" validate:"required,oneof=common adm nurse"`
	Name        string                `json:"name" validate:"required,max=120"`
	Cpf         string                `json:"cpf" validate:"required,document"`
	Email       string                `json:"email" validate:"required,email"`
	Password    string                `json:"password" validate:"required,min=8"`
	PhoneNumber string                `json:"phone_number" validate:"required,e164"`

	//common
	BirthDate  *string               `json:"birth_date" validate:"omitempty,datetime=2006-01-02"`
	Address    *AddressCreateBase    `json:"address" validate:"omitempty"`
	UserBaby   *[]UserBabyCreateBase `json:"user_baby" validate:"omitempty,max=20,dive"`
	ConsentLog *CreateConsentBase    `json:"consent_log" validate:"omitempty"`

	//nurse/admn
	ActionBy   *string `reqHeader:"action-by" validate:"omitempty,id"`
	Identifier *string `json:"identifier" validate:"omitempty,max=50"`
}

type CreateUserOptionalFields struct {
	CanCreateCommon bool
	CanCreateWorker bool
}

func (c CreateUserReq) ValidateCreateUserOptionalFields() CreateUserOptionalFields {
	return CreateUserOptionalFields{
		CanCreateCommon: c.Type == entities.EnumUserTypeCommon && c.BirthDate != nil && c.Address != nil && c.ActionBy == nil && c.ConsentLog != nil,
		CanCreateWorker: (c.Type == entities.EnumUserTypeAdmin || c.Type == entities.EnumUserTypeNurse) && c.Identifier != nil && c.ActionBy != nil,
	}
}

type CreateUserRes struct {
	entities.User
}
