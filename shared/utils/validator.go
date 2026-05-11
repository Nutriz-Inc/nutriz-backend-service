package utils

import (
	"net/mail"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/paemuri/brdoc/v2"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

func GetValidate() *validator.Validate {
	once.Do(func() {
		v := validator.New()

		_ = v.RegisterValidation("document", ValidDocument)
		_ = v.RegisterValidation("id", ValidId)
		_ = v.RegisterValidation("datetime", ValidDateTime)
		_ = v.RegisterValidation("email", ValidEmail)

		validateInstance = v
	})

	return validateInstance
}

// validators

func ValidDocument(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return true
	}

	if value == "00000000000" {
		return false
	}

	return brdoc.IsCPF(value) || brdoc.IsCNPJ(value)
}

func ValidId(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return true
	}

	return IdValidation(value)
}

func ValidDateTime(fl validator.FieldLevel) bool {
	field := fl.Field()

	if t, ok := field.Interface().(time.Time); ok {
		return !t.IsZero()
	}

	value := field.String()
	if value == "" {
		return true
	}
	layout := fl.Param()
	_, err := time.Parse(layout, value)
	return err == nil
}

func ValidEmail(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return true
	}

	_, err := mail.ParseAddress(value)
	return err == nil
}
