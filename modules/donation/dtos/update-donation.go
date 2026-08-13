package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateDonationReq struct {
	ActionBy        string   `reqHeader:"action-by" validate:"required,id"`
	IsActive        *bool    `json:"is_active" validate:"omitempty"`                  //adm
	QuantityDonated *float64 `json:"quantity_donated" validate:"omitempty,gt=0"`      //adm
	UserFeedback    *string  `json:"user_feedback" validate:"omitempty"`              //common
	ScoreFeedback   *int16   `json:"score_feedback" validate:"omitempty,gte=0,lte=5"` //common
	utils.GetReq
}

type UpdateDonationRes struct {
	entities.DonationOut
}

type UpdateDonationOptionalFields struct {
	HasIsActive        bool
	HasQuantityDonated bool
	HasFeedback        bool
	HasScoreFeedback   bool
}

func (c UpdateDonationReq) ValidateUpdateDonationOptionalFields() UpdateDonationOptionalFields {
	return UpdateDonationOptionalFields{
		HasIsActive:        c.IsActive != nil && !*c.IsActive,
		HasQuantityDonated: c.QuantityDonated != nil && *c.QuantityDonated > 0,
		HasFeedback:        c.UserFeedback != nil && *c.UserFeedback != "" && c.ScoreFeedback != nil,
	}
}
