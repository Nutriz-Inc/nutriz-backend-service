package entities

import "time"

type UserBaby struct {
	IdUserBaby string     `json:"id_user_baby" db:"id_user_baby"`
	IdUser     string     `json:"id_user" db:"id_user"`
	Name       *string    `json:"name" db:"name"`
	BirthDate  time.Time  `json:"birth_date" db:"birth_date"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" db:"updated_at"`
	RemovedAt  *time.Time `json:"removed_at" db:"removed_at"`
}

func (u UserBaby) TableName() string {
	return "user_baby"
}

func (u UserBaby) PrimaryKey() string {
	return "id_user_baby"
}

const MAX_BABY_QUANTITY_PER_USER = 20

type UserBabyOut struct {
	IdUserBaby string     `json:"id_user_baby"`
	IdUser     string     `json:"id_user"`
	Name       *string    `json:"name"`
	BirthDate  time.Time  `json:"birth_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	RemovedAt  *time.Time `json:"removed_at"`
}

func NewUserBabyOut(b UserBaby) UserBabyOut {
	return UserBabyOut{
		IdUserBaby: b.IdUserBaby,
		IdUser:     b.IdUser,
		Name:       b.Name,
		BirthDate:  b.BirthDate,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
		RemovedAt:  b.RemovedAt,
	}
}

func NewUserBabiesOut(babies []UserBaby) []UserBabyOut {
	out := make([]UserBabyOut, 0, len(babies))
	for _, b := range babies {
		out = append(out, NewUserBabyOut(b))
	}
	return out
}
