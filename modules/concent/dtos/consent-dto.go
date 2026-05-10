package dtos

import "time"

type CreateConsentReq struct {
	IdUser       string    `json:"id_user" validate:"required"`
	TermsVersion string    `json:"terms_version" validate:"required"`
	AcceptedAt   time.Time `json:"accepted_at" validate:"required"`
	IpAddress    string    `json:"ip_address" validate:"required"`
	UserAgent    string    `json:"user_agent" validate:"required"`
}

type CreateConsentRes struct {
	IdConsentLog string    `json:"id_consent_log"`
	IdUser       string    `json:"id_user"`
	TermsVersion string    `json:"terms_version"`
	AcceptedAt   time.Time `json:"accepted_at"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}