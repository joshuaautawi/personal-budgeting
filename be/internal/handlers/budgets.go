package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/httpjson"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/services"
	"personal-budgeting/be/internal/validate"
)

type Budgets struct {
	Svc     *services.BudgetService
	CatSvc  *services.CategoryService
}

func (h Budgets) List(c *fiber.Ctx) error {
	out, err := h.Svc.List(c.Context())
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

// Upsert matches the frontend action `upsertBudget`.
func (h Budgets) Upsert(c *fiber.Ctx) error {
	var in services.UpsertBudgetInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}

	// Validation belongs in handlers.
	in.Month = strings.TrimSpace(in.Month)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	if !validate.MonthKey(in.Month) || in.CategoryID == "" || in.AmountCents < 0 {
		return httpjson.WriteError(c, errs.ErrValidation)
	}
	cat, err := h.CatSvc.Get(c.Context(), in.CategoryID)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	if cat.Type != models.CategoryExpense {
		return httpjson.WriteError(c, errs.ErrValidation)
	}

	out, err := h.Svc.Upsert(c.Context(), in)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h Budgets) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.Svc.Delete(c.Context(), id); err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
