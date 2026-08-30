package entities

import "time"

type User struct {
	IdUser             string       `json:"id_user" db:"id_user"`
	InternalIdentifier *string      `json:"internal_identifier" db:"internal_identifier"`
	Type               EnumUserType `json:"type" db:"type"`
	Name               string       `json:"name" db:"name"`
	Cpf                string       `json:"cpf" db:"cpf"`
	BirthDate          time.Time    `json:"birth_date" db:"birth_date"`
	PhoneNumber        string       `json:"phone_number" db:"phone_number"`
	Email              string       `json:"email" db:"email"`
	Password           string       `json:"-" db:"password"`
	MilkDonated        *float64     `json:"milk_donated" db:"milk_donated"`
	CreatedAt          time.Time    `json:"created_at" db:"created_at"`
	CreatedBy          string       `json:"created_by" db:"created_by"`
	UpdatedAt          *time.Time   `json:"updated_at" db:"updated_at"`
	UpdatedBy          *string      `json:"updated_by" db:"updated_by"`
	RemovedAt          *time.Time   `json:"removed_at" db:"removed_at"`
	RemovedBy          *string      `json:"removed_by" db:"removed_by"`
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

type UserAction struct {
	CanCreateRoute              bool
	CanViewDashboard            bool
	CanCreateDonationStep       bool
	CanCreateDonation           bool
	CanViewDonationStepTimeline bool
	CanListDonations            bool
	CanUpdateDonationStep       bool
	CanUpdateDonation           bool
	CanCreateJob                bool
	CanViewJob                  bool
	CanListJobs                 bool
	CanRemoveJob                bool
	CanUpdateJob                bool
	CanCreateUser               bool
	CanListUsers                bool
	CanViewAnyUser              bool
	CanCreateAddress            bool
	CanUpdateAddress            bool
	CanRemoveAddress            bool
	CanCreateBaby               bool
	CanUpdateBaby               bool
	CanRemoveBaby               bool
	CanCreateConsentLog         bool
}

func (u User) Action() UserAction {
	return UserAction{
		CanCreateRoute:              u.Type == EnumUserTypeAdmin,
		CanViewDashboard:            u.Type == EnumUserTypeAdmin,
		CanCreateDonationStep:       u.Type == EnumUserTypeAdmin,
		CanCreateDonation:           u.Type == EnumUserTypeCommon,
		CanViewDonationStepTimeline: u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeNurse,
		CanListDonations:            u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeNurse || u.Type == EnumUserTypeCommon,
		CanUpdateDonationStep:       u.Type == EnumUserTypeAdmin,
		CanUpdateDonation:           u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeCommon,
		CanCreateJob:                u.Type == EnumUserTypeAdmin,
		CanViewJob:                  u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeNurse,
		CanListJobs:                 u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeNurse || u.Type == EnumUserTypeCommon,
		CanRemoveJob:                u.Type == EnumUserTypeAdmin,
		CanUpdateJob:                u.Type == EnumUserTypeAdmin || u.Type == EnumUserTypeNurse,
		CanCreateUser:               u.Type == EnumUserTypeAdmin,
		CanListUsers:                u.Type != EnumUserTypeCommon,
		CanViewAnyUser:              u.Type != EnumUserTypeCommon,
		CanCreateAddress:            u.Type == EnumUserTypeCommon || u.Type == EnumUserTypeAdmin,
		CanUpdateAddress:            u.Type == EnumUserTypeCommon,
		CanRemoveAddress:            u.Type == EnumUserTypeCommon,
		CanCreateBaby:               u.Type == EnumUserTypeCommon,
		CanUpdateBaby:               u.Type == EnumUserTypeCommon,
		CanRemoveBaby:               u.Type == EnumUserTypeCommon,
		CanCreateConsentLog:         u.Type == EnumUserTypeCommon,
	}
}

func (u User) CanRemoveUser(target User) bool {
	if u.Type == EnumUserTypeAdmin {
		return target.Type != EnumUserTypeCommon
	}

	if u.Type == EnumUserTypeCommon {
		return target.Type == EnumUserTypeCommon && target.IdUser == u.IdUser
	}

	return true
}
