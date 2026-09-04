package entities

import "time"

type User struct {
	IdUser              string       `json:"id_user" db:"id_user"`
	InternalIdentifier  *string      `json:"internal_identifier" db:"internal_identifier"`
	Type                EnumUserType `json:"type" db:"type"`
	Name                string       `json:"name" db:"name"`
	Cpf                 string       `json:"cpf" db:"cpf"`
	BirthDate           time.Time    `json:"birth_date" db:"birth_date"`
	PhoneNumber         string       `json:"phone_number" db:"phone_number"`
	Email               string       `json:"email" db:"email"`
	Password            string       `json:"-" db:"password"`
	MilkDonated         *float64     `json:"milk_donated" db:"milk_donated"`
	BloodExamValidUntil *time.Time   `json:"blood_exam_valid_until" db:"blood_exam_valid_until"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	CreatedBy           string       `json:"created_by" db:"created_by"`
	UpdatedAt           *time.Time   `json:"updated_at" db:"updated_at"`
	UpdatedBy           *string      `json:"updated_by" db:"updated_by"`
	RemovedAt           *time.Time   `json:"removed_at" db:"removed_at"`
	RemovedBy           *string      `json:"removed_by" db:"removed_by"`
}

type EnumUserType string

const (
	EnumUserTypeCommon EnumUserType = "common"
	EnumUserTypeAdmin  EnumUserType = "adm"
	EnumUserTypeNurse  EnumUserType = "nurse"
	EnumUserTypeDriver EnumUserType = "driver"
)

func (u User) TableName() string {
	return "user"
}

func (u User) PrimaryKey() string {
	return "id_user"
}

const BLOOD_EXAM_VALIDITY_MONTHS = 6

func BloodExamValidUntilFrom(from time.Time) time.Time {
	return from.AddDate(0, BLOOD_EXAM_VALIDITY_MONTHS, 0)
}

func (u User) IsBloodExamValid() bool {
	if u.BloodExamValidUntil == nil {
		return false
	}

	return u.BloodExamValidUntil.After(time.Now())
}

type UserAction struct {
	CanCreateRoute      bool
	CanListRoute        bool
	CanUpdateRoute      bool
	CanCreateRouteStop  bool
	CanRemoveRouteStop  bool
	CanUpdateRouteStop  bool
	CanListDonationStep bool
}

func (u User) Action() UserAction {
	return UserAction{
		CanCreateRoute:      u.Type == EnumUserTypeAdmin,
		CanListRoute:        u.Type != EnumUserTypeCommon,
		CanUpdateRoute:      u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeDriver,
		CanCreateRouteStop:  u.Type == EnumUserTypeAdmin,
		CanRemoveRouteStop:  u.Type == EnumUserTypeAdmin,
		CanUpdateRouteStop:  u.Type == EnumUserTypeDriver,
		CanListDonationStep: u.Type == EnumUserTypeAdmin,
	}
}
