package entities

import "time"

type Bottle struct {
	IdBottle          string     `json:"id_bottle" db:"id_bottle"`
	IdDonation        string     `json:"id_donation" db:"id_donation"`
	QuantityDonatedMl *float64   `json:"quantity_donated_ml" db:"quantity_donated_ml"`
	Discarded         *bool      `json:"discarded" db:"discarded"`
	Description       *string    `json:"description" db:"description"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	CreatedBy         string     `json:"created_by" db:"created_by"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy         *string    `json:"updated_by" db:"updated_by"`
	RemovedAt         *time.Time `json:"removed_at" db:"removed_at"`
	RemovedBy         *string    `json:"removed_by" db:"removed_by"`
}

func (b Bottle) TableName() string {
	return "bottle"
}

func (b Bottle) PrimaryKey() string {
	return "id_bottle"
}
