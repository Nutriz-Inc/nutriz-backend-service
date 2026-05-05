package utils

import "github.com/segmentio/ksuid"

func IdValidation(id string) bool {
	size := len(id)

	if size != 31 {
		return false
	}

	start := id[0:4]

	switch start {
	case "adr_", "dpt_", "fil_", "dst_", "don_", "job_", "usr_":
		return true
	}

	return false
}

type EnumEntityType string

const (
	AddressEntity       EnumEntityType = "address"
	DonationPointEntity EnumEntityType = "donation_point"
	DonationStepEntity  EnumEntityType = "donation_step"
	DonationEntity      EnumEntityType = "donation"
	FileEntity          EnumEntityType = "file"
	JobEntity           EnumEntityType = "job"
	UserEntity          EnumEntityType = "user"
)

func GetIdEntity(id string) EnumEntityType {
	start := id[0:4]

	switch start {
	case "adr_":
		return AddressEntity
	case "dpt_":
		return DonationPointEntity
	case "fil_":
		return FileEntity
	case "dst_":
		return DonationStepEntity
	case "don_":
		return DonationEntity
	case "job_":
		return JobEntity
	case "usr_":
		return UserEntity
	}

	return ""
}

func IdGenerate(entity EnumEntityType) string {
	id := ksuid.New().String()

	switch entity {
	case AddressEntity:
		return "adr_" + id
	case DonationPointEntity:
		return "dpt_" + id
	case FileEntity:
		return "fil_" + id
	case DonationStepEntity:
		return "dst_" + id
	case DonationEntity:
		return "don_" + id
	case JobEntity:
		return "job_" + id
	case UserEntity:
		return "usr_" + id
	}

	return id
}
