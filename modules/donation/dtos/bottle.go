package dto

type BottleUpdateBase struct {
	IdDonation        string   `json:"id_donation" validate:"required,id"`
	QuantityDonatedMl *float64 `json:"quantity_donated_ml" validate:"required,gte=0"`
	Discarded         *bool    `json:"discarded" validate:"omitempty"`
	Description       *string  `json:"description" validate:"omitempty,max=255"`
}
