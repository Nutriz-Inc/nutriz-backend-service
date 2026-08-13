package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
)

type CreateUserBabyReq struct {
	ActionBy string `json:"-" reqHeader:"action-by" validate:"required,id"`
	UserBabyCreateBase
}

type UserBabyCreateBase struct {
	Name      *string `json:"name" validate:"omitempty"`
	BirthDate string  `json:"birth_date" validate:"required,datetime=2006-01-02"`
}

type CreateUserBabyRes struct {
	sharedDto.UserBabyOut
}
