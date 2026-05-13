package entities

import "time"

type DonationStep struct {
	IdDonationStep string                 `json:"id_donation_step" db:"id_donation_step"`
	IdDonation     string                 `json:"id_donation" db:"id_donation"`
	Name           EnumDonationSteps      `json:"name" db:"name"`
	Description    string                 `json:"description" db:"description"`
	Status         EnumDonationStepStatus `json:"status" db:"status"`
	SetDate        *time.Time             `json:"set_date" db:"set_date"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time             `json:"updated_at" db:"updated_at"`
	CompletedAt    *time.Time             `json:"completed_at" db:"completed_at"`
}

type EnumDonationSteps string

const (
	EnumDonationStepBloodTest         EnumDonationSteps = "Exame de sangue"
	EnumDonationStepDeliverMilkingKit EnumDonationSteps = "Entregar kit de ordenha"
	EnumDonationStepCollectMilk       EnumDonationSteps = "Coletar leite"
	EnumDonationStepMilkAnalysis      EnumDonationSteps = "Análise de leite"
)

type EnumDonationStepStatus string

const (
	EnumDonationStepStatusPending EnumDonationStepStatus = "pending"
	EnumDonationStepStatusReview  EnumDonationStepStatus = "review"
	EnumDonationStepStatusDone    EnumDonationStepStatus = "done"
	EnumDonationStepStatusWarn    EnumDonationStepStatus = "warn"
	EnumDonationStepStatusFailed  EnumDonationStepStatus = "failed"
)

func (d DonationStep) TableName() string {
	return "donation_step"
}

func (d DonationStep) PrimaryKey() string {
	return "id_step_donation"
}

const NUMBER_OF_DONATION_STEPS = 4 //still not sure
