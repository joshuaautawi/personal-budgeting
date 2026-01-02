package handlers

import (
	"github.com/gofiber/fiber/v2"

	"personal-budgeting/be/internal/httpjson"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/services"
)

type State struct {
	Svc *services.StateService
}

func (h State) Get(c *fiber.Ctx) error {
	out, err := h.Svc.Get(c.Context())
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h State) Replace(c *fiber.Ctx) error {
	var st models.AppStateV1
	if err := c.BodyParser(&st); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}
	if err := h.Svc.Replace(c.Context(), st); err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}


