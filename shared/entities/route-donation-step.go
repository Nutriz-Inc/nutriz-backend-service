package entities

import "time"

type RouteDonationStep struct {
	IdRouteDonationStep string     `json:"id_route_donation_step" db:"id_route_donation_step"`
	IdRoute             string     `json:"id_route" db:"id_route"`
	IdDonationStep      string     `json:"id_donation_step" db:"id_donation_step"`
	StopOrder           *int16     `json:"stop_order" db:"stop_order"`
	DateStart           *time.Time `json:"date_start" db:"date_start"`
	DateEnd             *time.Time `json:"date_end" db:"date_end"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	CreatedBy           string     `json:"created_by" db:"created_by"`
	RemovedAt           *time.Time `json:"removed_at" db:"removed_at"`
	RemovedBy           *string    `json:"removed_by" db:"removed_by"`
}

func (r RouteDonationStep) TableName() string {
	return "route_donation_step"
}

func (r RouteDonationStep) PrimaryKey() string {
	return "id_route_donation_step"
}
