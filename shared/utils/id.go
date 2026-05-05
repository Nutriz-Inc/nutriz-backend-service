package utils

import (
	"github.com/segmentio/ksuid"
)

const PREFIX_LENGTH = 4

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

var prefixToEntity = map[string]EnumEntityType{
	"adr_": AddressEntity,
	"dpt_": DonationPointEntity,
	"fil_": FileEntity,
	"dst_": DonationStepEntity,
	"don_": DonationEntity,
	"job_": JobEntity,
	"usr_": UserEntity,
}

var entityToPrefix = func() map[EnumEntityType]string {
	m := make(map[EnumEntityType]string)
	for prefix, entity := range prefixToEntity {
		m[entity] = prefix
	}
	return m
}()

func IdValidation(id string) bool {
	if len(id) < PREFIX_LENGTH {
		return false
	}

	prefix := id[:PREFIX_LENGTH]

	if _, ok := prefixToEntity[prefix]; !ok {
		return false
	}

	_, err := ksuid.Parse(id[PREFIX_LENGTH:])
	return err == nil
}

func GetIdEntity(id string) EnumEntityType {
	if len(id) < PREFIX_LENGTH {
		return ""
	}

	return prefixToEntity[id[:PREFIX_LENGTH]]
}

func IdGenerate(entity EnumEntityType) string {
	prefix, ok := entityToPrefix[entity]
	if !ok {
		panic("invalid entity type: " + string(entity))
	}

	return prefix + ksuid.New().String()
}
