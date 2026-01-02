package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// These are the *database* models (GORM structs).
// API models that match the frontend live in `internal/models`.

type Category struct {
	ID          string `gorm:"primaryKey;type:text"`
	Type        string `gorm:"type:text;not null"`
	Name        string `gorm:"type:text;not null"`
	Description string `gorm:"type:text;not null;default:''"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Category) TableName() string { return "categories" }

type Budget struct {
	ID         string `gorm:"primaryKey;type:text"`
	Month      string `gorm:"type:text;not null;index:budgets_month_category_uq,unique"`
	CategoryID string `gorm:"type:text;not null;index:budgets_month_category_uq,unique;index"`
	AmountCents int64 `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Category Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (Budget) TableName() string { return "budgets" }

type Transaction struct {
	ID         string `gorm:"primaryKey;type:text"`
	Kind       string `gorm:"type:text;not null"`
	Date       string `gorm:"type:text;not null;index"`
	CategoryID string `gorm:"type:text;not null;index"`
	AmountCents int64 `gorm:"not null"`
	Note       string `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Category Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (Transaction) TableName() string { return "transactions" }

// Ensure GORM sees these models even if only referenced indirectly.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Category{}, &Budget{}, &Transaction{})
}


