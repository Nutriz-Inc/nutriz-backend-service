package dtos

import "nutriz-backend-service/shared/entities"

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
	entities.ConsentLogOut
}
