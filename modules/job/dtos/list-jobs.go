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
	entities.JobOut
	UserCommonName *string              `json:"user_common_name,omitempty"`
	UserNurseName  *string              `json:"user_nurse_name,omitempty"`
	Address        *entities.AddressOut `json:"address,omitempty"`
	IdDonation     *string              `json:"id_donation,omitempty"`
}

type ListJobsRes struct {
	Data []JobRes `json:"data"`
	utils.PaginationRes
}
