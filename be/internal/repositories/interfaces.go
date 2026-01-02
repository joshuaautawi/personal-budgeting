package repositories

import (
	"context"

	"personal-budgeting/be/internal/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
	Get(ctx context.Context, id string) (models.Category, error)
	Create(ctx context.Context, c models.Category) (models.Category, error)
	Update(ctx context.Context, id string, patch CategoryPatch) (models.Category, error)
	Delete(ctx context.Context, id string) error
}

type CategoryPatch struct {
	Name        *string
	Description *string
	UpdatedAt   *string
}

type BudgetRepository interface {
	List(ctx context.Context) ([]models.Budget, error)
	Get(ctx context.Context, id string) (models.Budget, error)
	Upsert(ctx context.Context, b models.Budget) (models.Budget, error)
	Delete(ctx context.Context, id string) error
	FindByMonthCategory(ctx context.Context, month string, categoryID string) (models.Budget, bool, error)
}

type TxnRepository interface {
	List(ctx context.Context) ([]models.Txn, error)
	Get(ctx context.Context, id string) (models.Txn, error)
	Create(ctx context.Context, t models.Txn) (models.Txn, error)
	Update(ctx context.Context, id string, patch TxnPatch) (models.Txn, error)
	Delete(ctx context.Context, id string) error
	CountByCategory(ctx context.Context, categoryID string) (int, error)
}

type TxnPatch struct {
	Kind        *models.TransactionKind
	Date        *string
	CategoryID  *string
	AmountCents *int64
	Note        *string
	UpdatedAt   *string
}
