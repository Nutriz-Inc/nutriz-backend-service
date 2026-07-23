package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type GetJobReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type JobInfoRes struct {
	entities.Job
	IdUserCommon *string `json:"id_user_common,omitempty" db:"id_user_common"`
	IdAddress    *string `json:"id_address,omitempty" db:"id_address"`
}

type GetJobRes struct {
	JobInfoRes
}
