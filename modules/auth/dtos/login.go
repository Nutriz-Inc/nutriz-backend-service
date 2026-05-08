package dtos

import "nutriz-backend-service/shared/entities"

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRes struct {
	Token     string                `json:"token"`
	IdUser    string                `json:"id_user"`
	Name      string                `json:"name"`
	Type      entities.EnumUserType `json:"type"`
	Addresses *[]entities.Address   `json:"addresses,omitempty"`
}
