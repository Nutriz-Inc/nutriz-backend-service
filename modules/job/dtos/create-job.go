package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
)

type CreateJobReq struct {
	ActionBy    string  `reqHeader:"action-by" validate:"required,id"`
	IdUser      string  `json:"id_user" validate:"required,id"`
	IdStep      string  `json:"id_step" validate:"required,id"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description" validate:"required,max=100"`
	DateSet     *string `json:"date_set" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type CreateJobRes struct {
	sharedDto.JobOut
}
