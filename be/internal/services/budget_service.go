package services

import (
	"context"
	"time"

	"personal-budgeting/be/internal/clock"
	"personal-budgeting/be/internal/id"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/repositories"
)

type BudgetService struct {
	clk clock.Clock
	ids id.Generator

	budgets repositories.BudgetRepository
}

func NewBudgetService(clk clock.Clock, ids id.Generator, budgets repositories.BudgetRepository) *BudgetService {
	return &BudgetService{clk: clk, ids: ids, budgets: budgets}
}

func (s *BudgetService) List(ctx context.Context) ([]models.Budget, error) {
	return s.budgets.List(ctx)
}

type UpsertBudgetInput struct {
	Month       string `json:"month"`
	CategoryID  string `json:"categoryId"`
	AmountCents int64  `json:"amountCents"`
}

func (s *BudgetService) Upsert(ctx context.Context, in UpsertBudgetInput) (models.Budget, error) {
	existing, ok, err := s.budgets.FindByMonthCategory(ctx, in.Month, in.CategoryID)
	if err != nil {
		return models.Budget{}, err
	}

	now := s.clk.Now().Format(time.RFC3339)
	if ok {
		existing.AmountCents = in.AmountCents
		existing.UpdatedAt = now
		return s.budgets.Upsert(ctx, existing)
	}

	b := models.Budget{
		ID:          s.ids.NewID(),
		Month:       in.Month,
		CategoryID:  in.CategoryID,
		AmountCents: in.AmountCents,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.budgets.Upsert(ctx, b)
}

func (s *BudgetService) Delete(ctx context.Context, id string) error {
	return s.budgets.Delete(ctx, id)
}
