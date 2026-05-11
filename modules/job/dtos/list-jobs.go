package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListJobsReq struct {
	ActionBy string  `reqHeader:"action-by" validate:"required,id"`
	DateSet  *string `query:"date_set" validate:"omitempty,datetime=2006-01-02"`
	utils.PaginationReq
}

type ListJobsRes struct {
	Data []entities.Job `json:"data"`
	utils.PaginationRes
}
