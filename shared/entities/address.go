package entities

import (
	"time"
)

type Address struct {
	IdAddress       string     `json:"id_address" db:"id_address"`
	IdUser          *string    `json:"id_user" db:"id_user"`
	IdDonationPoint *string    `json:"id_donation_point" db:"id_donation_point"`
	Zipcode         string     `json:"zipcode" db:"zipcode"`
	Street          string     `json:"street" db:"street"`
	Number          *string    `json:"number" db:"number"`
	City            string     `json:"city" db:"city"`
	State           string     `json:"state" db:"state"`
	Neighborhood    string     `json:"neighborhood" db:"neighborhood"`
	Complement      *string    `json:"complement" db:"complement"`
	Latitude        *float64   `json:"latitude" db:"latitude"`
	Longitude       *float64   `json:"longitude" db:"longitude"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	RemovedAt       *time.Time `json:"removed_at" db:"removed_at"`
}

func (a Address) TableName() string {
	return "address"
}

func (a Address) PrimaryKey() string {
	return "id_address"
}

const MAX_ADDRESS_QUANTITY_PER_USER = 1
