package httpjson

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"personal-budgeting/be/internal/errs"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errs.ErrValidation):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "validation"})
	case errors.Is(err, errs.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "not_found"})
	case errors.Is(err, errs.ErrConflict):
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Error: "conflict"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "internal"})
	}
}
