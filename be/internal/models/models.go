package models

type CategoryType string

const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

type TransactionKind string

const (
	KindIncome  TransactionKind = "income"
	KindExpense TransactionKind = "expense"
)

type Category struct {
	ID          string       `json:"id"`
	Type        CategoryType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
}

type Budget struct {
	ID          string `json:"id"`
	Month       string `json:"month"` // YYYY-MM
	CategoryID  string `json:"categoryId"`
	AmountCents int64  `json:"amountCents"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Txn struct {
	ID          string          `json:"id"`
	Kind        TransactionKind `json:"kind"`
	Date        string          `json:"date"` // YYYY-MM-DD
	CategoryID  string          `json:"categoryId"`
	AmountCents int64           `json:"amountCents"`
	Note        string          `json:"note,omitempty"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type AppStateV1 struct {
	Version      int        `json:"version"`
	Categories   []Category `json:"categories"`
	Budgets      []Budget   `json:"budgets"`
	Transactions []Txn      `json:"transactions"`
}
