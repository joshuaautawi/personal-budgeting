package services

import (
	"context"

	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/repositories"
)

type StateService struct {
	cats    repositories.CategoryRepository
	budgets repositories.BudgetRepository
	txns    repositories.TxnRepository
}

func NewStateService(cats repositories.CategoryRepository, budgets repositories.BudgetRepository, txns repositories.TxnRepository) *StateService {
	return &StateService{cats: cats, budgets: budgets, txns: txns}
}

func (s *StateService) Get(ctx context.Context) (models.AppStateV1, error) {
	cats, err := s.cats.List(ctx)
	if err != nil {
		return models.AppStateV1{}, err
	}
	budgets, err := s.budgets.List(ctx)
	if err != nil {
		return models.AppStateV1{}, err
	}
	txns, err := s.txns.List(ctx)
	if err != nil {
		return models.AppStateV1{}, err
	}
	return models.AppStateV1{
		Version:      1,
		Categories:   cats,
		Budgets:      budgets,
		Transactions: txns,
	}, nil
}

// Replace replaces all stored data with the provided state.
// This is intended for local/dev sync; production apps should use proper auth and per-user storage.
func (s *StateService) Replace(ctx context.Context, st models.AppStateV1) error {
	if st.Version != 1 {
		return errs.ErrValidation
	}

	// Only supported when repos expose reset/bulk operations (memory repos do).
	type catBulk interface {
		Reset()
		BulkUpsert([]models.Category)
	}
	type budgetBulk interface {
		Reset()
		BulkUpsert([]models.Budget)
	}
	type txnBulk interface {
		Reset()
		BulkUpsert([]models.Txn)
	}

	catR, ok1 := s.cats.(catBulk)
	budgetR, ok2 := s.budgets.(budgetBulk)
	txnR, ok3 := s.txns.(txnBulk)
	if !ok1 || !ok2 || !ok3 {
		return errs.ErrConflict
	}

	_ = ctx // future: validate referential integrity using ctx

	catR.Reset()
	budgetR.Reset()
	txnR.Reset()
	catR.BulkUpsert(st.Categories)
	budgetR.BulkUpsert(st.Budgets)
	txnR.BulkUpsert(st.Transactions)
	return nil
}
