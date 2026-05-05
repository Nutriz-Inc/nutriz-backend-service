package entities

import "time"

type Donation struct {
	IdDonation      string     `json:"id_donation" db:"id_donation"`
	IdDonationPoint *string    `json:"id_donation_point" db:"id_donation_point"`
	Quantity        *float64   `json:"quantity" db:"quantity"`
	UserFeedback    *string    `json:"user_feedback" db:"user_feedback"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CreatedBy       string     `json:"created_by" db:"created_by"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy       *string    `json:"updated_by" db:"updated_by"`
	RemovedAt       *time.Time `json:"removed_at" db:"removed_at"`
	RemovedBy       *string    `json:"removed_by" db:"removed_by"`
}

func (d Donation) TableName() string {
	return "donation"
}

func (d Donation) PrimaryKey() string {
	return "id_donation"
}
