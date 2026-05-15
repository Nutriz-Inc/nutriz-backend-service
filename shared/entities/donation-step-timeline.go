package entities

import "time"

type DonationStepTimeline struct {
	IdDonationStepTimeline string                 `json:"id_donation_step_timeline" db:"id_donation_step_timeline"`
	IdDonationStep         string                 `json:"id_donation_step" db:"id_donation_step"`
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
