package entities

import "time"

type File struct {
	IdFile    string     `json:"id_file" db:"id_file"`
	IdJob     *string    `json:"id_job" db:"id_job"`
	FilePath  string     `json:"file_path" db:"file_path"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy *string    `json:"updated_by" db:"updated_by"`
	RemovedAt *time.Time `json:"removed_at" db:"removed_at"`
	RemovedBy *string    `json:"removed_by" db:"removed_by"`
}

func (f File) TableName() string {
	return "file"
}

func (f File) PrimaryKey() string {
	return "id_file"
}
