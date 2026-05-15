package entities

import "time"

type Job struct {
	IdJob        string        `json:"id_job" db:"id_job"`
	IdUser       string        `json:"id_user" db:"id_user"`
	IdStep       string        `json:"id_step" db:"id_step"`
	Status       EnumJobStatus `json:"status" db:"status"`
	Name         string        `json:"name" db:"name"`
	Description  string        `json:"description" db:"description"`
	DateSet      *time.Time    `json:"date_set" db:"date_set"`
	UserFeedback *string       `json:"user_feedback" db:"user_feedback"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	CreatedBy    string        `json:"created_by" db:"created_by"`
	UpdatedAt    *time.Time    `json:"updated_at" db:"updated_at"`
	UpdatedBy    *string       `json:"updated_by" db:"updated_by"`
	RemovedAt    *time.Time    `json:"removed_at" db:"removed_at"`
	RemovedBy    *string       `json:"removed_by" db:"removed_by"`
}

type EnumJobStatus string

const (
	EnumJobStatusPending EnumJobStatus = "pending"
	EnumJobStatusDone    EnumJobStatus = "done"
	EnumJobStatusFailed  EnumJobStatus = "failed"
)

func (j Job) TableName() string {
	return "job"
}

func (j Job) PrimaryKey() string {
	return "id_job"
}
