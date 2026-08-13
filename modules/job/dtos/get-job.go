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
	entities.JobOut
	IdUserCommon *string `json:"id_user_common,omitempty"`
	IdAddress    *string `json:"id_address,omitempty"`
}

type GetJobRes struct {
	JobInfoRes
}
