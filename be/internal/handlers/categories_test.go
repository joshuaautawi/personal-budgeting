package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"personal-budgeting/be/internal/repositories"
	"personal-budgeting/be/internal/router"
	"personal-budgeting/be/internal/services"
	"personal-budgeting/be/internal/testutil"
)

func TestCategoriesHandler_CreateAndList(t *testing.T) {
	clk := testutil.FixedClock{T: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)}
	ids := &testutil.SeqID{}

	gdb := testutil.NewTestGormDB(t)
	catRepo := repositories.NewGormCategoryRepo(gdb)
	budgetRepo := repositories.NewGormBudgetRepo(gdb)
	txnRepo := repositories.NewGormTxnRepo(gdb)

	categorySvc := services.NewCategoryService(clk, ids, catRepo)
	budgetSvc := services.NewBudgetService(clk, ids, budgetRepo)
	txnSvc := services.NewTxnService(clk, ids, txnRepo)
	stateSvc := services.NewStateService(catRepo, budgetRepo, txnRepo)

	app := router.New(router.Deps{
		Category:    categorySvc,
		Budget:      budgetSvc,
		Transaction: txnSvc,
		State:       stateSvc,
	})

	body, _ := json.Marshal(map[string]any{
		"type": "expense",
		"name": "Groceries",
	})
	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /categories: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	req2 := httptest.NewRequest("GET", "/api/v1/categories", nil)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("GET /categories: %v", err)
	}
	if resp2.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
}
