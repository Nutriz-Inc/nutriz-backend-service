package utils

import (
	"net/http"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

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
