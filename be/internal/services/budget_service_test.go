package services

import (
	"context"
	"testing"
	"time"

	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/repositories"
	"personal-budgeting/be/internal/testutil"
)

func TestBudgetService_Upsert_RejectsIncomeCategory(t *testing.T) {
	// Validation now lives in handlers, not services.
	t.Skip("validation is enforced in handlers")
}

func TestBudgetService_Upsert_CreatesAndUpdates(t *testing.T) {
	ctx := context.Background()
	clk := testutil.FixedClock{T: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)}
	ids := &testutil.SeqID{}

	gdb := testutil.NewTestGormDB(t)
	catRepo := repositories.NewGormCategoryRepo(gdb)
	budgetRepo := repositories.NewGormBudgetRepo(gdb)

	_, _ = catRepo.Create(ctx, models.Category{
		ID:        "cat-expense",
		Type:      models.CategoryExpense,
		Name:      "Groceries",
		CreatedAt: clk.Now().Format(time.RFC3339),
		UpdatedAt: clk.Now().Format(time.RFC3339),
	})

	svc := NewBudgetService(clk, ids, budgetRepo)

	b1, err := svc.Upsert(ctx, UpsertBudgetInput{
		Month:       "2026-01",
		CategoryID:  "cat-expense",
		AmountCents: 400_00,
	})
	if err != nil {
		t.Fatalf("upsert create: %v", err)
	}
	if b1.AmountCents != 400_00 {
		t.Fatalf("expected 40000, got %d", b1.AmountCents)
	}

	b2, err := svc.Upsert(ctx, UpsertBudgetInput{
		Month:       "2026-01",
		CategoryID:  "cat-expense",
		AmountCents: 500_00,
	})
	if err != nil {
		t.Fatalf("upsert update: %v", err)
	}
	if b2.ID != b1.ID {
		t.Fatalf("expected same ID on upsert update")
	}
	if b2.AmountCents != 500_00 {
		t.Fatalf("expected 50000, got %d", b2.AmountCents)
	}
}
