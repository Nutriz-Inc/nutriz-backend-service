package entities

import "time"

type Route struct {
	IdRoute   string          `json:"id_route" db:"id_route"`
	IdDriver  string          `json:"id_driver" db:"id_driver"`
	Status    EnumRouteStatus `json:"status" db:"status"`
	DateStart *time.Time      `json:"date_start" db:"date_start"`
	DateEnd   *time.Time      `json:"date_end" db:"date_end"`
	Mileage   *float64        `json:"mileage" db:"mileage"`
	DateSet   time.Time       `json:"date_set" db:"date_set"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	CreatedBy string          `json:"created_by" db:"created_by"`
	UpdatedAt *time.Time      `json:"updated_at" db:"updated_at"`
	UpdatedBy *string         `json:"updated_by" db:"updated_by"`
	RemovedAt *time.Time      `json:"removed_at" db:"removed_at"`
	RemovedBy *string         `json:"removed_by" db:"removed_by"`
}

type EnumRouteStatus string

const (
	EnumRouteStatusPending    EnumRouteStatus = "pending"
	EnumRouteStatusInProgress EnumRouteStatus = "in_progress"
	EnumRouteStatusDone       EnumRouteStatus = "done"
	EnumRouteStatusCanceled   EnumRouteStatus = "canceled"
)

func (r Route) TableName() string {
	return "route"
}

func (r Route) PrimaryKey() string {
	return "id_route"
}
