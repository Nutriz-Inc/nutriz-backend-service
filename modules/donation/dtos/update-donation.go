package dto

import (
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
)

type UpdateDonationReq struct {
	ActionBy      string              `reqHeader:"action-by" validate:"required,id"`
	IsActive      *bool               `json:"is_active" validate:"omitempty"`                  //adm
	Bottles       *[]BottleUpdateBase `json:"bottles" validate:"omitempty,max=50,dive"`        //adm
	UserFeedback  *string             `json:"user_feedback" validate:"omitempty"`              //common
	ScoreFeedback *int16              `json:"score_feedback" validate:"omitempty,gte=0,lte=5"` //common
	utils.GetReq
}

type UpdateDonationRes struct {
	entities.Donation
}

type UpdateDonationOptionalFields struct {
	HasIsActive bool
	HasBottles  bool
	HasFeedback bool
}

func (c UpdateDonationReq) ValidateUpdateDonationOptionalFields() UpdateDonationOptionalFields {
	return UpdateDonationOptionalFields{
		HasIsActive: c.IsActive != nil && !*c.IsActive,
		HasBottles:  c.Bottles != nil && len(*c.Bottles) > 0,
		HasFeedback: c.UserFeedback != nil && *c.UserFeedback != "" && c.ScoreFeedback != nil,
	}
}
