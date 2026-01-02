export type CategoryType = 'income' | 'expense'
export type TransactionKind = 'income' | 'expense'

export type MonthKey = string // YYYY-MM
export type DateKey = string // YYYY-MM-DD

export type Id = string

export type Cents = number

export type AppStateV1 = {
  version: 1
  categories: Category[]
  budgets: Budget[]
  transactions: Txn[]
}

export type Category = {
  id: Id
  type: CategoryType
  name: string
  description?: string
  createdAt: string
  updatedAt: string
}

export type Budget = {
  id: Id
  month: MonthKey
  categoryId: Id // must be expense category
  amountCents: Cents
  createdAt: string
  updatedAt: string
}

export type Txn = {
  id: Id
  kind: TransactionKind
  date: DateKey
  categoryId: Id // must match kind
  amountCents: Cents
  note?: string
  createdAt: string
  updatedAt: string
}


