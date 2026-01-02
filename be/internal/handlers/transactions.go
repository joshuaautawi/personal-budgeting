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

type Transactions struct {
	Svc    *services.TxnService
	CatSvc *services.CategoryService
}

func (h Transactions) List(c *fiber.Ctx) error {
	out, err := h.Svc.List(c.Context())
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h Transactions) Create(c *fiber.Ctx) error {
	var in services.CreateTxnInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}

	// Validation belongs in handlers.
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	if in.Kind != models.KindIncome && in.Kind != models.KindExpense {
		return httpjson.WriteError(c, errs.ErrValidation)
	}
	if !validate.DateKey(in.Date) || in.CategoryID == "" || in.AmountCents <= 0 {
		return httpjson.WriteError(c, errs.ErrValidation)
	}
	cat, err := h.CatSvc.Get(c.Context(), in.CategoryID)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	if (in.Kind == models.KindIncome && cat.Type != models.CategoryIncome) || (in.Kind == models.KindExpense && cat.Type != models.CategoryExpense) {
		return httpjson.WriteError(c, errs.ErrValidation)
	}

	out, err := h.Svc.Create(c.Context(), in)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(out)
}

func (h Transactions) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var in services.UpdateTxnInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httpjson.ErrorResponse{Error: "bad_json"})
	}

	// Validate patch + enforce kind/category compatibility.
	existing, err := h.Svc.Get(c.Context(), id)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	nextKind := existing.Kind
	nextCatID := existing.CategoryID

	if in.Kind != nil {
		if *in.Kind != models.KindIncome && *in.Kind != models.KindExpense {
			return httpjson.WriteError(c, errs.ErrValidation)
		}
		nextKind = *in.Kind
	}
	if in.Date != nil {
		if !validate.DateKey(*in.Date) {
			return httpjson.WriteError(c, errs.ErrValidation)
		}
	}
	if in.CategoryID != nil {
		trimmed := strings.TrimSpace(*in.CategoryID)
		if trimmed == "" {
			return httpjson.WriteError(c, errs.ErrValidation)
		}
		nextCatID = trimmed
		in.CategoryID = &trimmed
	}
	if in.AmountCents != nil {
		if *in.AmountCents <= 0 {
			return httpjson.WriteError(c, errs.ErrValidation)
		}
	}
	if in.Note != nil {
		trimmed := strings.TrimSpace(*in.Note)
		in.Note = &trimmed
	}

	cat, err := h.CatSvc.Get(c.Context(), nextCatID)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	if (nextKind == models.KindIncome && cat.Type != models.CategoryIncome) || (nextKind == models.KindExpense && cat.Type != models.CategoryExpense) {
		return httpjson.WriteError(c, errs.ErrValidation)
	}

	out, err := h.Svc.Update(c.Context(), id, in)
	if err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.JSON(out)
}

func (h Transactions) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.Svc.Delete(c.Context(), id); err != nil {
		return httpjson.WriteError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}


