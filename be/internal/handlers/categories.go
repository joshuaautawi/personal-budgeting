package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/httpjson"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/services"
)

type Categories struct {
	Svc       *services.CategoryService
	BudgetSvc *services.BudgetService
	TxnSvc    *services.TxnService
}

func (h Categories) List(c *fiber.Ctx) error {
	out, err := h.Svc.List(c.Context())
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h Categories) Create(c *fiber.Ctx) error {
	var in services.CreateCategoryInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.Name == "" {
		return httpjson.WriteError(c, errs.ErrValidation)
	}
	if in.Type != models.CategoryIncome && in.Type != models.CategoryExpense {
		return httpjson.WriteError(c, errs.ErrValidation)
	}
	out, err := h.Svc.Create(c.Context(), in)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(out)
}

func (h Categories) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var in services.UpdateCategoryInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}
	if in.Name != nil {
		trimmed := strings.TrimSpace(*in.Name)
		if trimmed == "" {
			return httpjson.WriteError(c, errs.ErrValidation)
		}
		in.Name = &trimmed
	}
	if in.Description != nil {
		trimmed := strings.TrimSpace(*in.Description)
		in.Description = &trimmed
	}
	out, err := h.Svc.Update(c.Context(), id, in)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h Categories) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validation/business rule: disallow delete if referenced.
	budgets, err := h.BudgetSvc.List(c.Context())
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	for _, b := range budgets {
		if b.CategoryID == id {
			return httpjson.WriteError(c, errs.ErrConflict)
		}
	}
	n, err := h.TxnSvc.CountByCategory(c.Context(), id)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	if n > 0 {
		return httpjson.WriteError(c, errs.ErrConflict)
	}

	if err := h.Svc.Delete(c.Context(), id); err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
