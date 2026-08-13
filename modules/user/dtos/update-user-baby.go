package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateUserBabyReq struct {
	ActionBy  string  `reqHeader:"action-by" validate:"required,id"`
	Name      *string `json:"name" validate:"omitempty"`
	BirthDate *string `json:"birth_date" validate:"omitempty,datetime=2006-01-02"`
	utils.GetReq
}

type UpdateUserBabyRes struct {
	entities.UserBabyOut
}

type UpdateUserBabyOptionalFields struct {
	HasName      bool
	HasBirthDate bool
}

func (c UpdateUserBabyReq) ValidateUpdateUserBabyOptionalFields() UpdateUserBabyOptionalFields {
	return UpdateUserBabyOptionalFields{
		HasName:      c.Name != nil && *c.Name != "",
		HasBirthDate: c.BirthDate != nil && *c.BirthDate != "",
	}
}
