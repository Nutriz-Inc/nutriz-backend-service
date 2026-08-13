package entities

import "time"

type DonationStepTimeline struct {
	IdDonationStepTimeline string                 `json:"id_donation_step_timeline" db:"id_donation_step_timeline"`
	IdDonationStep         string                 `json:"id_donation_step" db:"id_donation_step"`
	IdAddress              *string                `json:"id_address" db:"id_address"`
	Description            string                 `json:"description" db:"description"`
	Status                 EnumDonationStepStatus `json:"status" db:"status"`
	SetDate                *time.Time             `json:"set_date" db:"set_date"`
	CreatedAt              time.Time              `json:"created_at" db:"created_at"`
	CreatedBy              string                 `json:"created_by" db:"created_by"`
}

func (d DonationStepTimeline) TableName() string {
	return "donation_step_timeline"
}

func (d DonationStepTimeline) PrimaryKey() string {
	return "id_donation_step_timeline"
}

type DonationStepTimelineOut struct {
	IdDonationStepTimeline string                 `json:"id_donation_step_timeline"`
	IdDonationStep         string                 `json:"id_donation_step"`
	IdAddress              *string                `json:"id_address"`
	Description            string                 `json:"description"`
	Status                 EnumDonationStepStatus `json:"status"`
	SetDate                *time.Time             `json:"set_date"`
	CreatedAt              time.Time              `json:"created_at"`
	CreatedBy              string                 `json:"created_by"`
}

func NewDonationStepTimelineOut(t DonationStepTimeline) DonationStepTimelineOut {
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

func NewDonationStepTimelinesOut(timeline []DonationStepTimeline) []DonationStepTimelineOut {
	out := make([]DonationStepTimelineOut, 0, len(timeline))
	for _, t := range timeline {
		out = append(out, NewDonationStepTimelineOut(t))
	}
	return out
}
