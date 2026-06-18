package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateUserReq struct {
	ActionBy			string					`reqHeader:"action-by" validate:"omitempty,id"`
	InternalIdentifier	*string					`json:"internal_identifier" validate:"omitempty,max=50"`
	Type				*entities.EnumUserType	`json:"type" validate:"omitempty,oneof=adm nurse"`
	Name				*string					`json:"name" validate:"omitempty,max=120"`
	PhoneNumber			*string					`json:"phone_number" validate:"omitempty,e164"`
	Email				*string					`json:"email" validate:"omitempty,email"`
	Password			*string					`json:"password" validate:"omitempty,min=8"`
	utils.GetReq
}

type UpdateUserRes struct {
	entities.User
}