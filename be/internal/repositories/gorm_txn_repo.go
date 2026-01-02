package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"personal-budgeting/be/internal/dbmodel"
	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/models"
)

type GormTxnRepo struct {
	db *gorm.DB
}

func NewGormTxnRepo(db *gorm.DB) *GormTxnRepo {
	return &GormTxnRepo{db: db}
}

var _ TxnRepository = (*GormTxnRepo)(nil)

func (r *GormTxnRepo) List(ctx context.Context) ([]models.Txn, error) {
	var rows []dbmodel.Transaction
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.Txn, 0, len(rows))
	for _, t := range rows {
		out = append(out, toAPITxn(t))
	}
	return out, nil
}

func (r *GormTxnRepo) Get(ctx context.Context, id string) (models.Txn, error) {
	var row dbmodel.Transaction
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if isNotFound(err) {
			return models.Txn{}, errs.ErrNotFound
		}
		return models.Txn{}, err
	}
	return toAPITxn(row), nil
}

func (r *GormTxnRepo) Create(ctx context.Context, t models.Txn) (models.Txn, error) {
	row, err := toDBTxn(t)
	if err != nil {
		return models.Txn{}, err
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		if isUniqueViolation(err) {
			return models.Txn{}, errs.ErrConflict
		}
		if isForeignKeyViolation(err) {
			return models.Txn{}, errs.ErrValidation
		}
		return models.Txn{}, err
	}
	return r.Get(ctx, t.ID)
}

func (r *GormTxnRepo) Update(ctx context.Context, id string, patch TxnPatch) (models.Txn, error) {
	updates := map[string]any{}
	if patch.Kind != nil {
		updates["kind"] = string(*patch.Kind)
	}
	if patch.Date != nil {
		updates["date"] = *patch.Date
	}
	if patch.CategoryID != nil {
		updates["category_id"] = *patch.CategoryID
	}
	if patch.AmountCents != nil {
		updates["amount_cents"] = *patch.AmountCents
	}
	if patch.Note != nil {
		updates["note"] = *patch.Note
	}
	if patch.UpdatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *patch.UpdatedAt); err == nil {
			updates["updated_at"] = t
		}
	}
	if len(updates) == 0 {
		return r.Get(ctx, id)
	}
	tx := r.db.WithContext(ctx).Model(&dbmodel.Transaction{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		if isForeignKeyViolation(tx.Error) {
			return models.Txn{}, errs.ErrValidation
		}
		return models.Txn{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return models.Txn{}, errs.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *GormTxnRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&dbmodel.Transaction{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *GormTxnRepo) CountByCategory(ctx context.Context, categoryID string) (int, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&dbmodel.Transaction{}).Where("category_id = ?", categoryID).Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
}

func (r *GormTxnRepo) Reset() {
	_ = r.db.WithContext(context.Background()).Exec("DELETE FROM transactions").Error
}

func (r *GormTxnRepo) BulkUpsert(items []models.Txn) {
	ctx := context.Background()
	for _, t := range items {
		row, err := toDBTxn(t)
		if err != nil {
			continue
		}
		_ = r.db.WithContext(ctx).Save(&row).Error
	}
}

func toAPITxn(t dbmodel.Transaction) models.Txn {
	return models.Txn{
		ID:          t.ID,
		Kind:        models.TransactionKind(t.Kind),
		Date:        t.Date,
		CategoryID:  t.CategoryID,
		AmountCents: t.AmountCents,
		Note:        t.Note,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toDBTxn(t models.Txn) (dbmodel.Transaction, error) {
	createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}
	updatedAt, err := time.Parse(time.RFC3339, t.UpdatedAt)
	if err != nil {
		updatedAt = createdAt
	}
	return dbmodel.Transaction{
		ID:          t.ID,
		Kind:        string(t.Kind),
		Date:        t.Date,
		CategoryID:  t.CategoryID,
		AmountCents: t.AmountCents,
		Note:        t.Note,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}


