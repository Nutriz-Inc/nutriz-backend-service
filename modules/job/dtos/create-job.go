package dtos

import (
	"nutriz-backend-service/shared/entities"
	"time"
)

type CreateJobReq struct {
	ActionBy     string     `reqHeader:"action-by" validate:"required"`
	IdStep       string     `json:"id_step" validate:"required"`
	Name         string     `json:"name" validate:"required"`
	Description  string     `json:"description" validate:"required"`
	DateSet      *time.Time `json:"date_set" validate:"omitempty,future"`
	UserFeedback *string    `json:"user_feedback"`
}

type CreateJobRes struct {
	entities.Job
}
