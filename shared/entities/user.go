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
)

func (u User) TableName() string {
	return "user"
}

func (u User) PrimaryKey() string {
	return "id_user"
}

type UserOut struct {
	IdUser             string       `json:"id_user"`
	InternalIdentifier *string      `json:"internal_identifier"`
	Type               EnumUserType `json:"type"`
	Name               string       `json:"name"`
	Cpf                string       `json:"cpf"`
	BirthDate          time.Time    `json:"birth_date"`
	PhoneNumber        string       `json:"phone_number"`
	Email              string       `json:"email"`
	MilkDonated        *float64     `json:"milk_donated"`
	CreatedAt          time.Time    `json:"created_at"`
	CreatedBy          string       `json:"created_by"`
	UpdatedAt          *time.Time   `json:"updated_at"`
	UpdatedBy          *string      `json:"updated_by"`
	RemovedAt          *time.Time   `json:"removed_at"`
	RemovedBy          *string      `json:"removed_by"`
}

func NewUserOut(u User) UserOut {
	return UserOut{
		IdUser:             u.IdUser,
		InternalIdentifier: u.InternalIdentifier,
		Type:               u.Type,
		Name:               u.Name,
		Cpf:                u.Cpf,
		BirthDate:          u.BirthDate,
		PhoneNumber:        u.PhoneNumber,
		Email:              u.Email,
		MilkDonated:        u.MilkDonated,
		CreatedAt:          u.CreatedAt,
		CreatedBy:          u.CreatedBy,
		UpdatedAt:          u.UpdatedAt,
		UpdatedBy:          u.UpdatedBy,
		RemovedAt:          u.RemovedAt,
		RemovedBy:          u.RemovedBy,
	}
}

func NewUsersOut(users []User) []UserOut {
	out := make([]UserOut, 0, len(users))
	for _, u := range users {
		out = append(out, NewUserOut(u))
	}
	return out
}
