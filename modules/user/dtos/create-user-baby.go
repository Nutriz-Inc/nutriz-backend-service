package dtos

import (
	"nutriz-backend-service/shared/entities"
	"time"
)

type CreateUserBabyReq struct {
	ActionBy  string    `reqHeader:"action-by" validate:"required,id"`
	Name      *string   `json:"name" validate:"omitempty"`
	BirthDate time.Time `json:"birth_date" validate:"required,datetime=2006-01-02"`
}

type CreateUserBabyRes struct {
	entities.UserBaby
}
