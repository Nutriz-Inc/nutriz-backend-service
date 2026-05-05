package utils

import (
	"net/mail"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/paemuri/brdoc/v2"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

func getValidate() *validator.Validate {
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

func Validate(data interface{}) (bool, []ErrorResponse) {
	validate := getValidate()

	err := validate.Struct(data)
	if err == nil {
		return false, nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return true, []ErrorResponse{
			{
				Error: true,
			},
		}
	}

	errors := make([]ErrorResponse, 0, len(validationErrors))

	for _, e := range validationErrors {
		errors = append(errors, ErrorResponse{
			FailedField: e.Field(),
			Tag:         e.Tag(),
			Value:       e.Value(),
			Error:       true,
		})
	}

	return true, errors
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
	value := fl.Field().String()

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
