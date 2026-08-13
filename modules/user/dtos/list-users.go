package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListUsersReq struct {
	Name               *string                `query:"name" validate:"omitempty,max=120"`
	Type               *entities.EnumUserType `query:"type" validate:"omitempty,oneof=common adm nurse"`
	InternalIdentifier *string                `query:"internal_identifier" validate:"omitempty,max=36"`
	Cpf                *string                `query:"cpf" validate:"omitempty,document"`
	ActionBy           string                 `reqHeader:"action-by" validate:"required,id"`
	utils.PaginationReq
}

type ListUsersRes struct {
	Data []sharedDto.UserOut `json:"data"`
	utils.PaginationRes
}
