package utils

import (
	"net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string      `json:"field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
	Error       bool        `json:"error"`
}

func ErrorValidation(err []ErrorResponse) *fluxgo.GlobalError {
	return &fluxgo.GlobalError{
		Message: "Error on validate data",
		Code:    "validation.err",
		Status:  http.StatusBadRequest,
		Success: false,
		Errors:  err,
	}
}

func ErrorParser(err error, message, code string) *fluxgo.GlobalError {
	return &fluxgo.GlobalError{
		Message: message,
		Code:    code,
		Status:  http.StatusBadRequest,
		Success: false,
		Errors:  err,
	}
}

func ErrorRes(err fluxgo.GlobalError, message, code string, statusCode int) *fluxgo.GlobalError {
	return &fluxgo.GlobalError{
		Message: message,
		Code:    code,
		Status:  statusCode,
		Success: false,
		Errors:  err,
	}
}

func ErrorUnauthorizedFiber(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"message": "Unauthorized",
		"code":    "unauthorized",
		"status":  http.StatusUnauthorized,
		"success": false,
	})
}

func ErrorUnauthorized(message, code string) *fluxgo.GlobalError {
	return &fluxgo.GlobalError{
		Message: message,
		Code:    code,
		Status:  http.StatusUnauthorized,
		Success: false,
	}
}

func ErrorForbidden(message, code string) *fluxgo.GlobalError {
	return &fluxgo.GlobalError{
		Message: message,
		Code:    code,
		Status:  http.StatusForbidden,
		Success: false,
	}
}
