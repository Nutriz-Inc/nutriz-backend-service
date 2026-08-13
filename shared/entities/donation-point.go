package entities

import "time"

type DonationPoint struct {
	IdDonationPoint string     `json:"id_donation_point" db:"id_donation_point"`
	Name            string     `json:"name" db:"name"`
	Description     *string    `json:"description" db:"description"`
	HasHome         bool       `json:"has_home" db:"has_home"`
	PhoneNumber     *string    `json:"phone_number" db:"phone_number"`
	Email           *string    `json:"email" db:"email"`
	OpeningHours    *string    `json:"opening_hours" db:"opening_hours"`
	RemovedAt       *time.Time `json:"removed_at" db:"removed_at"`
}

func (d DonationPoint) TableName() string {
	return "donation_point"
}

func (d DonationPoint) PrimaryKey() string {
	return "id_donation_point"
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

func NewDonationPointOut(dp DonationPoint) DonationPointOut {
	return DonationPointOut(dp)
}
