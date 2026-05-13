package dtos

import (
	"nutriz-backend-service/shared/entities"
)

type CreateUserBabyReq struct {
	ActionBy  string    `json:"-" reqHeader:"action-by" validate:"required,id"`
	Name      *string   `json:"name" validate:"omitempty"`
	BirthDate string	`json:"birth_date" validate:"required,datetime=2006-01-02"`
}

type CreateUserBabyRes struct {
	entities.UserBaby
}