package dtos

import (
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/utils"
)

type GetJobReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	utils.GetReq
}

type JobInfoRes struct {
	sharedDto.JobOut
	IdUserCommon *string `json:"id_user_common,omitempty"`
	IdAddress    *string `json:"id_address,omitempty"`
}

type GetJobRes struct {
	JobInfoRes
}
