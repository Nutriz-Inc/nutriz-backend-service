// Package dtos contains API output structures shared across modules.
//
// These types exist so that HTTP responses never serialize persisted
// entities (shared/entities) directly. Each Out type mirrors the public
// fields of its entity and is built via the New*Out constructor, which is
// the single place deciding what a persisted record exposes over the API.
package dtos

import (
	"nutriz-backend-service/shared/entities"
	"time"
)

type UserOut struct {
	IdUser             string                `json:"id_user"`
	InternalIdentifier *string               `json:"internal_identifier"`
	Type               entities.EnumUserType `json:"type"`
	Name               string                `json:"name"`
	Cpf                string                `json:"cpf"`
	BirthDate          time.Time             `json:"birth_date"`
	PhoneNumber        string                `json:"phone_number"`
	Email              string                `json:"email"`
	MilkDonated        *float64              `json:"milk_donated"`
	CreatedAt          time.Time             `json:"created_at"`
	CreatedBy          string                `json:"created_by"`
	UpdatedAt          *time.Time            `json:"updated_at"`
	UpdatedBy          *string               `json:"updated_by"`
	RemovedAt          *time.Time            `json:"removed_at"`
	RemovedBy          *string               `json:"removed_by"`
}

func NewUserOut(u entities.User) UserOut {
	return UserOut{
		IdUser:             u.IdUser,
		InternalIdentifier: u.InternalIdentifier,
		Type:               u.Type,
		Name:               u.Name,
		Cpf:                u.Cpf,
		BirthDate:          u.BirthDate,
		PhoneNumber:        u.PhoneNumber,
		Email:              u.Email,
		MilkDonated:        u.MilkDonated,
		CreatedAt:          u.CreatedAt,
		CreatedBy:          u.CreatedBy,
		UpdatedAt:          u.UpdatedAt,
		UpdatedBy:          u.UpdatedBy,
		RemovedAt:          u.RemovedAt,
		RemovedBy:          u.RemovedBy,
	}
}

func NewUsersOut(users []entities.User) []UserOut {
	out := make([]UserOut, 0, len(users))
	for _, u := range users {
		out = append(out, NewUserOut(u))
	}
	return out
}

type AddressOut struct {
	IdAddress       string     `json:"id_address"`
	IdUser          *string    `json:"id_user"`
	IdDonationPoint *string    `json:"id_donation_point"`
	Zipcode         string     `json:"zipcode"`
	Street          string     `json:"street"`
	Number          *string    `json:"number"`
	City            string     `json:"city"`
	State           string     `json:"state"`
	Neighborhood    string     `json:"neighborhood"`
	Complement      *string    `json:"complement"`
	Latitude        *float64   `json:"latitude"`
	Longitude       *float64   `json:"longitude"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	RemovedAt       *time.Time `json:"removed_at"`
}

func NewAddressOut(a entities.Address) AddressOut {
	return AddressOut{
		IdAddress:       a.IdAddress,
		IdUser:          a.IdUser,
		IdDonationPoint: a.IdDonationPoint,
		Zipcode:         a.Zipcode,
		Street:          a.Street,
		Number:          a.Number,
		City:            a.City,
		State:           a.State,
		Neighborhood:    a.Neighborhood,
		Complement:      a.Complement,
		Latitude:        a.Latitude,
		Longitude:       a.Longitude,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
		RemovedAt:       a.RemovedAt,
	}
}

func NewAddressesOut(addresses []entities.Address) []AddressOut {
	out := make([]AddressOut, 0, len(addresses))
	for _, a := range addresses {
		out = append(out, NewAddressOut(a))
	}
	return out
}

type UserBabyOut struct {
	IdUserBaby string     `json:"id_user_baby"`
	IdUser     string     `json:"id_user"`
	Name       *string    `json:"name"`
	BirthDate  time.Time  `json:"birth_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	RemovedAt  *time.Time `json:"removed_at"`
}

func NewUserBabyOut(b entities.UserBaby) UserBabyOut {
	return UserBabyOut{
		IdUserBaby: b.IdUserBaby,
		IdUser:     b.IdUser,
		Name:       b.Name,
		BirthDate:  b.BirthDate,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
		RemovedAt:  b.RemovedAt,
	}
}

func NewUserBabiesOut(babies []entities.UserBaby) []UserBabyOut {
	out := make([]UserBabyOut, 0, len(babies))
	for _, b := range babies {
		out = append(out, NewUserBabyOut(b))
	}
	return out
}

type ConsentLogOut struct {
	IdConsentLog string    `json:"id_consent_log"`
	IdUser       string    `json:"id_user"`
	TermsVersion string    `json:"terms_version"`
	AcceptedAt   time.Time `json:"accepted_at"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

func NewConsentLogOut(c entities.ConsentLog) ConsentLogOut {
	return ConsentLogOut{
		IdConsentLog: c.IdConsentLog,
		IdUser:       c.IdUser,
		TermsVersion: c.TermsVersion,
		AcceptedAt:   c.AcceptedAt,
		IpAddress:    c.IpAddress,
		UserAgent:    c.UserAgent,
	}
}

type DonationOut struct {
	IdDonation      string     `json:"id_donation"`
	IsActive        bool       `json:"is_active"`
	QuantityDonated *float64   `json:"quantity_donated"`
	UserFeedback    *string    `json:"user_feedback"`
	ScoreFeedback   *int16     `json:"score_feedback"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       string     `json:"created_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
	UpdatedBy       *string    `json:"updated_by"`
	RemovedAt       *time.Time `json:"removed_at"`
	RemovedBy       *string    `json:"removed_by"`
}

func NewDonationOut(d entities.Donation) DonationOut {
	return DonationOut{
		IdDonation:      d.IdDonation,
		IsActive:        d.IsActive,
		QuantityDonated: d.QuantityDonated,
		UserFeedback:    d.UserFeedback,
		ScoreFeedback:   d.ScoreFeedback,
		CreatedAt:       d.CreatedAt,
		CreatedBy:       d.CreatedBy,
		UpdatedAt:       d.UpdatedAt,
		UpdatedBy:       d.UpdatedBy,
		RemovedAt:       d.RemovedAt,
		RemovedBy:       d.RemovedBy,
	}
}

type DonationStepOut struct {
	IdDonationStep string                          `json:"id_donation_step"`
	IdDonation     string                          `json:"id_donation"`
	IdAddress      *string                         `json:"id_address"`
	Name           entities.EnumDonationSteps      `json:"name"`
	Description    string                          `json:"description"`
	Status         entities.EnumDonationStepStatus `json:"status"`
	SetDate        *time.Time                      `json:"set_date"`
	CreatedAt      time.Time                       `json:"created_at"`
	CreatedBy      string                          `json:"created_by"`
	UpdatedAt      *time.Time                      `json:"updated_at"`
	UpdatedBy      *string                         `json:"updated_by"`
	CompletedAt    *time.Time                      `json:"completed_at"`
}

func NewDonationStepOut(s entities.DonationStep) DonationStepOut {
	return DonationStepOut{
		IdDonationStep: s.IdDonationStep,
		IdDonation:     s.IdDonation,
		IdAddress:      s.IdAddress,
		Name:           s.Name,
		Description:    s.Description,
		Status:         s.Status,
		SetDate:        s.SetDate,
		CreatedAt:      s.CreatedAt,
		CreatedBy:      s.CreatedBy,
		UpdatedAt:      s.UpdatedAt,
		UpdatedBy:      s.UpdatedBy,
		CompletedAt:    s.CompletedAt,
	}
}

func NewDonationStepsOut(steps []entities.DonationStep) []DonationStepOut {
	out := make([]DonationStepOut, 0, len(steps))
	for _, s := range steps {
		out = append(out, NewDonationStepOut(s))
	}
	return out
}

type DonationStepTimelineOut struct {
	IdDonationStepTimeline string                          `json:"id_donation_step_timeline"`
	IdDonationStep         string                          `json:"id_donation_step"`
	IdAddress              *string                         `json:"id_address"`
	Description            string                          `json:"description"`
	Status                 entities.EnumDonationStepStatus `json:"status"`
	SetDate                *time.Time                      `json:"set_date"`
	CreatedAt              time.Time                       `json:"created_at"`
	CreatedBy              string                          `json:"created_by"`
}

func NewDonationStepTimelineOut(t entities.DonationStepTimeline) DonationStepTimelineOut {
	return DonationStepTimelineOut{
		IdDonationStepTimeline: t.IdDonationStepTimeline,
		IdDonationStep:         t.IdDonationStep,
		IdAddress:              t.IdAddress,
		Description:            t.Description,
		Status:                 t.Status,
		SetDate:                t.SetDate,
		CreatedAt:              t.CreatedAt,
		CreatedBy:              t.CreatedBy,
	}
}

func NewDonationStepTimelinesOut(timeline []entities.DonationStepTimeline) []DonationStepTimelineOut {
	out := make([]DonationStepTimelineOut, 0, len(timeline))
	for _, t := range timeline {
		out = append(out, NewDonationStepTimelineOut(t))
	}
	return out
}

type DonationPointOut struct {
	IdDonationPoint string     `json:"id_donation_point"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	HasHome         bool       `json:"has_home"`
	PhoneNumber     *string    `json:"phone_number"`
	Email           *string    `json:"email"`
	OpeningHours    *string    `json:"opening_hours"`
	RemovedAt       *time.Time `json:"removed_at"`
}

func NewDonationPointOut(dp entities.DonationPoint) DonationPointOut {
	return DonationPointOut{
		IdDonationPoint: dp.IdDonationPoint,
		Name:            dp.Name,
		Description:     dp.Description,
		HasHome:         dp.HasHome,
		PhoneNumber:     dp.PhoneNumber,
		Email:           dp.Email,
		OpeningHours:    dp.OpeningHours,
		RemovedAt:       dp.RemovedAt,
	}
}

type JobOut struct {
	IdJob        string                 `json:"id_job"`
	IdUser       string                 `json:"id_user"`
	IdStep       string                 `json:"id_step"`
	Status       entities.EnumJobStatus `json:"status"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	DateSet      *time.Time             `json:"date_set"`
	UserFeedback *string                `json:"user_feedback"`
	CreatedAt    time.Time              `json:"created_at"`
	CreatedBy    string                 `json:"created_by"`
	UpdatedAt    *time.Time             `json:"updated_at"`
	UpdatedBy    *string                `json:"updated_by"`
	RemovedAt    *time.Time             `json:"removed_at"`
	RemovedBy    *string                `json:"removed_by"`
}

func NewJobOut(j entities.Job) JobOut {
	return JobOut{
		IdJob:        j.IdJob,
		IdUser:       j.IdUser,
		IdStep:       j.IdStep,
		Status:       j.Status,
		Name:         j.Name,
		Description:  j.Description,
		DateSet:      j.DateSet,
		UserFeedback: j.UserFeedback,
		CreatedAt:    j.CreatedAt,
		CreatedBy:    j.CreatedBy,
		UpdatedAt:    j.UpdatedAt,
		UpdatedBy:    j.UpdatedBy,
		RemovedAt:    j.RemovedAt,
		RemovedBy:    j.RemovedBy,
	}
}
