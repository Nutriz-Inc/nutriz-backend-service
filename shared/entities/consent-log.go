package entities

import "time"

type ConsentLog struct {
	IdConsentLog string    `json:"id_consent_log" db:"id_consent_log"`
	IdUser       string    `json:"id_user" db:"id_user"`
	TermsVersion string    `json:"terms_version" db:"terms_version"`
	AcceptedAt   time.Time `json:"accepted_at" db:"accepted_at"`
	IpAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
}

func (c ConsentLog) TableName() string {
	return "consent_log"
}

func (c ConsentLog) PrimaryKey() string {
	return "id_consent_log"
}

type ConsentLogOut struct {
	IdConsentLog string    `json:"id_consent_log"`
	IdUser       string    `json:"id_user"`
	TermsVersion string    `json:"terms_version"`
	AcceptedAt   time.Time `json:"accepted_at"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

func NewConsentLogOut(c ConsentLog) ConsentLogOut {
	return ConsentLogOut{
		IdConsentLog: c.IdConsentLog,
		IdUser:       c.IdUser,
		TermsVersion: c.TermsVersion,
		AcceptedAt:   c.AcceptedAt,
		IpAddress:    c.IpAddress,
		UserAgent:    c.UserAgent,
	}
}
