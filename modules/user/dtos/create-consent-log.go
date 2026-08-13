package dtos

import sharedDto "nutriz-backend-service/shared/dtos"

type CreateConsentReq struct {
	ActionBy string `reqHeader:"action-by" validate:"required,id"`
	CreateConsentBase
}

type CreateConsentBase struct {
	TermsVersion string `json:"terms_version" validate:"required"`
	IpAddress    string `json:"ip_address" validate:"required,ip"`
	UserAgent    string `json:"user_agent" validate:"required"`
}

type CreateConsentRes struct {
	sharedDto.ConsentLogOut
}
