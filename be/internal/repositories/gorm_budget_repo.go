package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"personal-budgeting/be/internal/dbmodel"
	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/models"
)

type GormBudgetRepo struct {
	db *gorm.DB
}

func NewGormBudgetRepo(db *gorm.DB) *GormBudgetRepo {
	return &GormBudgetRepo{db: db}
}

var _ BudgetRepository = (*GormBudgetRepo)(nil)

func (r *GormBudgetRepo) List(ctx context.Context) ([]models.Budget, error) {
	var rows []dbmodel.Budget
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.Budget, 0, len(rows))
	for _, b := range rows {
		out = append(out, toAPIBudget(b))
	}
	return out, nil
}

func (r *GormBudgetRepo) Get(ctx context.Context, id string) (models.Budget, error) {
	var row dbmodel.Budget
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if isNotFound(err) {
			return models.Budget{}, errs.ErrNotFound
		}
		return models.Budget{}, err
	}
	return toAPIBudget(row), nil
}

func (r *GormBudgetRepo) FindByMonthCategory(ctx context.Context, month string, categoryID string) (models.Budget, bool, error) {
	var row dbmodel.Budget
	err := r.db.WithContext(ctx).First(&row, "month = ? AND category_id = ?", month, categoryID).Error
	if err != nil {
		if isNotFound(err) {
			return models.Budget{}, false, nil
		}
		return models.Budget{}, false, err
	}
	return toAPIBudget(row), true, nil
}

func (r *GormBudgetRepo) Upsert(ctx context.Context, b models.Budget) (models.Budget, error) {
	row, err := toDBBudget(b)
	if err != nil {
		return models.Budget{}, err
	}

	err = r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "month"}, {Name: "category_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount_cents", "updated_at"}),
		}).
		Create(&row).Error
	if err != nil {
		if isForeignKeyViolation(err) {
			return models.Budget{}, errs.ErrValidation
		}
		return models.Budget{}, err
	}

	out, ok, err := r.FindByMonthCategory(ctx, b.Month, b.CategoryID)
	if err != nil {
		return models.Budget{}, err
	}
	if !ok {
		return models.Budget{}, errs.ErrConflict
	}
	return out, nil
}

func (r *GormBudgetRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&dbmodel.Budget{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *GormBudgetRepo) Reset() {
	_ = r.db.WithContext(context.Background()).Exec("DELETE FROM budgets").Error
}

func (r *GormBudgetRepo) BulkUpsert(items []models.Budget) {
	ctx := context.Background()
	for _, b := range items {
		row, err := toDBBudget(b)
		if err != nil {
			continue
		}
		_ = r.db.WithContext(ctx).Save(&row).Error
	}
}

func toAPIBudget(b dbmodel.Budget) models.Budget {
	return models.Budget{
		ID:          b.ID,
		Month:       b.Month,
		CategoryID:  b.CategoryID,
		AmountCents: b.AmountCents,
		CreatedAt:   b.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toDBBudget(b models.Budget) (dbmodel.Budget, error) {
	createdAt, err := time.Parse(time.RFC3339, b.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}
	updatedAt, err := time.Parse(time.RFC3339, b.UpdatedAt)
	if err != nil {
		updatedAt = createdAt
	}
	return dbmodel.Budget{
		ID:          b.ID,
		Month:       b.Month,
		CategoryID:  b.CategoryID,
		AmountCents: b.AmountCents,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}


