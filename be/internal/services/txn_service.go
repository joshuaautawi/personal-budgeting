package services

import (
	"context"
	"strings"
	"time"

	"personal-budgeting/be/internal/clock"
	"personal-budgeting/be/internal/id"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/repositories"
)

type TxnService struct {
	clk clock.Clock
	ids id.Generator

	txns repositories.TxnRepository
}

func NewTxnService(clk clock.Clock, ids id.Generator, txns repositories.TxnRepository) *TxnService {
	return &TxnService{clk: clk, ids: ids, txns: txns}
}

func (s *TxnService) List(ctx context.Context) ([]models.Txn, error) {
	return s.txns.List(ctx)
}

func (s *TxnService) Get(ctx context.Context, id string) (models.Txn, error) {
	return s.txns.Get(ctx, id)
}

func (s *TxnService) CountByCategory(ctx context.Context, categoryID string) (int, error) {
	return s.txns.CountByCategory(ctx, categoryID)
}

type CreateTxnInput struct {
	Kind        models.TransactionKind `json:"kind"`
	Date        string                 `json:"date"`
	CategoryID  string                 `json:"categoryId"`
	AmountCents int64                  `json:"amountCents"`
	Note        string                 `json:"note,omitempty"`
}

func (s *TxnService) Create(ctx context.Context, in CreateTxnInput) (models.Txn, error) {
	now := s.clk.Now().Format(time.RFC3339)
	t := models.Txn{
		ID:          s.ids.NewID(),
		Kind:        in.Kind,
		Date:        in.Date,
		CategoryID:  strings.TrimSpace(in.CategoryID),
		AmountCents: in.AmountCents,
		Note:        strings.TrimSpace(in.Note),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.txns.Create(ctx, t)
}

type UpdateTxnInput struct {
	Kind        *models.TransactionKind `json:"kind"`
	Date        *string                 `json:"date"`
	CategoryID  *string                 `json:"categoryId"`
	AmountCents *int64                  `json:"amountCents"`
	Note        *string                 `json:"note"`
}

func (s *TxnService) Update(ctx context.Context, id string, in UpdateTxnInput) (models.Txn, error) {
	now := s.clk.Now().Format(time.RFC3339)
	patch := repositories.TxnPatch{
		UpdatedAt: &now,
	}
	if in.Kind != nil {
		patch.Kind = in.Kind
	}
	if in.Date != nil {
		patch.Date = in.Date
	}
	if in.CategoryID != nil {
		trimmed := strings.TrimSpace(*in.CategoryID)
		patch.CategoryID = &trimmed
	}
	if in.AmountCents != nil {
		patch.AmountCents = in.AmountCents
	}
	if in.Note != nil {
		trimmed := strings.TrimSpace(*in.Note)
		patch.Note = &trimmed
	}
	return s.txns.Update(ctx, id, patch)
}

func (s *TxnService) Delete(ctx context.Context, id string) error {
	return s.txns.Delete(ctx, id)
}
