package dtos

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type ListJobsReq struct {
	ActionBy       *string                 `reqHeader:"action-by" validate:"required,id"`
	DateSet        *string                 `query:"date_set" validate:"omitempty,datetime=2006-01-02"`
	IdStep         *string                 `query:"id_step" validate:"omitempty,id"`
	IdUserCommon   *string                 `query:"id_user_common" validate:"omitempty,id"`
	IdUserNurse    *string                 `query:"id_user_nurse" validate:"omitempty,id"`
	UserCommonName *string                 `query:"user_common_name" validate:"omitempty"`
	UserNurseName  *string                 `query:"user_nurse_name" validate:"omitempty"`
	Status         *entities.EnumJobStatus `query:"status" validate:"omitempty,oneof=pending done failed"`
	utils.PaginationReq
}

type JobRes struct {
	entities.Job
	UserCommonName *string           `json:"user_common_name,omitempty" db:"user_common_name"`
	UserNurseName  *string           `json:"user_nurse_name,omitempty" db:"user_nurse_name"`
	Address        *entities.Address `json:"address,omitempty" db:"-"`
	IdDonation     *string           `json:"id_donation,omitempty" db:"id_donation"`
}

type ListJobsRes struct {
	Data []JobRes `json:"data"`
	utils.PaginationRes
}
