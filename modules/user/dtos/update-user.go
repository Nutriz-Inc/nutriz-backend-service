package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateUserReq struct {
	ActionBy    string  `reqHeader:"action-by" validate:"omitempty,id"`
	Name        *string `json:"name" validate:"omitempty,max=120"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Password    *string `json:"password" validate:"omitempty,min=8"`
	utils.GetReq
}

type UpdateUserRes struct {
	entities.User
}

type UpdateUserOptionalFields struct {
	HasInternalIdentifier bool
	HasType               bool
	HasName               bool
	HasPhoneNumber        bool
	HasEmail              bool
	HasPassword           bool
}

func (c UpdateUserReq) ValidateUpdateUserOptionalFields() UpdateUserOptionalFields {
	return UpdateUserOptionalFields{
		HasName:        c.Name != nil && *c.Name != "",
		HasPhoneNumber: c.PhoneNumber != nil && *c.PhoneNumber != "",
		HasEmail:       c.Email != nil && *c.Email != "",
		HasPassword:    c.Password != nil && *c.Password != "",
	}
}
