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
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (a Address) TableName() string {
	return "address"
}

func (a Address) PrimaryKey() string {
	return "id_address"
}

const MAX_ADDRESS_QUANTITY_PER_USER = 1

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

func NewAddressOut(a Address) AddressOut {
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

func NewAddressesOut(addresses []Address) []AddressOut {
	out := make([]AddressOut, 0, len(addresses))
	for _, a := range addresses {
		out = append(out, NewAddressOut(a))
	}
	return out
}
