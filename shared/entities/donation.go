package entities

import "time"

type Donation struct {
	IdDonation      string     `json:"id_donation" db:"id_donation"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	QuantityDonated *float64   `json:"quantity_donated" db:"quantity_donated"`
	UserFeedback    *string    `json:"user_feedback" db:"user_feedback"`
	ScoreFeedback   *int16     `json:"score_feedback" db:"score_feedback"`
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

const MAX_DONATION_ACTIVE_QUANTITY_PER_USER = 1

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

func NewDonationOut(d Donation) DonationOut {
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
