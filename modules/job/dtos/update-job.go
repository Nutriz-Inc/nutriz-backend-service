package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateJobReq struct {
	ActionBy     string                  `reqHeader:"action-by" validate:"required,id"`
	IdUser       *string                 `json:"id_user" validate:"omitempty,id"`
	Status       *entities.EnumJobStatus `json:"status" validate:"omitempty,oneof=done failed"`
	Description  *string                 `json:"description" validate:"omitempty,max=100"`
	DateSet      *string                 `json:"date_set" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	UserFeedback *string                 `json:"user_feedback" validate:"omitempty,max=500"`
	utils.GetReq
}

type UpdateJobRes struct {
	sharedDto.JobOut
}

type UpdateJobOptionalFields struct {
	HasIdUser       bool
	HasStatus       bool
	HasDescription  bool
	HasDateSet      bool
	HasUserFeedback bool
}

func (c UpdateJobReq) ValidateUpdateJobOptionalFields() UpdateJobOptionalFields {
	return UpdateJobOptionalFields{
		HasIdUser:       c.IdUser != nil && *c.IdUser != "",
		HasStatus:       c.Status != nil && *c.Status != "",
		HasDescription:  c.Description != nil && *c.Description != "",
		HasDateSet:      c.DateSet != nil && *c.DateSet != "",
		HasUserFeedback: c.UserFeedback != nil && *c.UserFeedback != "",
	}
}
