package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"personal-budgeting/be/internal/dbmodel"
	"personal-budgeting/be/internal/errs"
	"personal-budgeting/be/internal/models"
)

type GormCategoryRepo struct {
	db *gorm.DB
}

func NewGormCategoryRepo(db *gorm.DB) *GormCategoryRepo {
	return &GormCategoryRepo{db: db}
}

var _ CategoryRepository = (*GormCategoryRepo)(nil)

func (r *GormCategoryRepo) List(ctx context.Context) ([]models.Category, error) {
	var rows []dbmodel.Category
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.Category, 0, len(rows))
	for _, c := range rows {
		out = append(out, toAPICategory(c))
	}
	return out, nil
}

func (r *GormCategoryRepo) Get(ctx context.Context, id string) (models.Category, error) {
	var row dbmodel.Category
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if isNotFound(err) {
			return models.Category{}, errs.ErrNotFound
		}
		return models.Category{}, err
	}
	return toAPICategory(row), nil
}

func (r *GormCategoryRepo) Create(ctx context.Context, c models.Category) (models.Category, error) {
	row, err := toDBCategory(c)
	if err != nil {
		return models.Category{}, err
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		if isUniqueViolation(err) {
			return models.Category{}, errs.ErrConflict
		}
		return models.Category{}, err
	}
	return r.Get(ctx, c.ID)
}

func (r *GormCategoryRepo) Update(ctx context.Context, id string, patch CategoryPatch) (models.Category, error) {
	updates := map[string]any{}
	if patch.Name != nil {
		updates["name"] = *patch.Name
	}
	if patch.Description != nil {
		updates["description"] = *patch.Description
	}
	if patch.UpdatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *patch.UpdatedAt); err == nil {
			updates["updated_at"] = t
		}
	}
	if len(updates) == 0 {
		return r.Get(ctx, id)
	}
	tx := r.db.WithContext(ctx).Model(&dbmodel.Category{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return models.Category{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return models.Category{}, errs.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *GormCategoryRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&dbmodel.Category{}, "id = ?", id)
	if tx.Error != nil {
		if isForeignKeyViolation(tx.Error) {
			return errs.ErrConflict
		}
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// Reset/BulkUpsert are used by StateService.Replace via type assertion.
func (r *GormCategoryRepo) Reset() {
	_ = r.db.WithContext(context.Background()).Exec("DELETE FROM categories").Error
}

func (r *GormCategoryRepo) BulkUpsert(items []models.Category) {
	ctx := context.Background()
	for _, c := range items {
		row, err := toDBCategory(c)
		if err != nil {
			continue
		}
		_ = r.db.WithContext(ctx).Save(&row).Error
	}
}

func toAPICategory(c dbmodel.Category) models.Category {
	return models.Category{
		ID:          c.ID,
		Type:        models.CategoryType(c.Type),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toDBCategory(c models.Category) (dbmodel.Category, error) {
	createdAt, err := time.Parse(time.RFC3339, c.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}
	updatedAt, err := time.Parse(time.RFC3339, c.UpdatedAt)
	if err != nil {
		updatedAt = createdAt
	}
	return dbmodel.Category{
		ID:          c.ID,
		Type:        string(c.Type),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}


