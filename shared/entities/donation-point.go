package entities

import "time"

type DonationPoint struct {
	IdDonationPoint string     `json:"id_donation_point" db:"id_donation_point"`
	Name            string     `json:"name" db:"name"`
	Description     *string    `json:"description" db:"description"`
	Cnpj            string     `json:"cnpj" db:"cnpj"`
	HasHome         bool       `json:"has_home" db:"has_home"`
	PhoneNumber     *string    `json:"phone_number" db:"phone_number"`
	Email           *string    `json:"email" db:"email"`
	OpeningHours    *string    `json:"opening_hours" db:"opening_hours"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CreatedBy       string     `json:"created_by" db:"created_by"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy       *string    `json:"updated_by" db:"updated_by"`
	RemovedAt       *time.Time `json:"removed_at" db:"removed_at"`
	RemovedBy       *string    `json:"removed_by" db:"removed_by"`
}

func (d DonationPoint) TableName() string {
	return "donation_point"
}

func (d DonationPoint) PrimaryKey() string {
	return "id_donation_point"
}
