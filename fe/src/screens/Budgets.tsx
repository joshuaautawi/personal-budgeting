import { useMemo, useState } from 'react'
import { formatCents, parseMoneyToCents } from '../lib/money'
import { isInMonth } from '../lib/month'
import type { Id } from '../lib/types'
import { useAppStore } from '../store/AppStore'

export function Budgets(props: { month: string }) {
  const { state, upsertBudget, deleteBudget } = useAppStore()
  const expenseCategories = useMemo(() => state.categories.filter((c) => c.type === 'expense'), [state.categories])

  const [categoryId, setCategoryId] = useState<Id>('')
  const [amount, setAmount] = useState('')

  const monthExpensesByCategory = useMemo(() => {
    const map = new Map<Id, number>()
    for (const t of state.transactions) {
      if (t.kind !== 'expense') continue
      if (!isInMonth(t.date, props.month)) continue
      map.set(t.categoryId, (map.get(t.categoryId) ?? 0) + t.amountCents)
    }
    return map
  }, [state.transactions, props.month])

  const monthBudgets = useMemo(() => state.budgets.filter((b) => b.month === props.month), [state.budgets, props.month])

  const totalBudgeted = useMemo(() => monthBudgets.reduce((sum, b) => sum + b.amountCents, 0), [monthBudgets])

  const rows = useMemo(() => {
    const byCat = new Map<Id, typeof monthBudgets[number]>()
    for (const b of monthBudgets) byCat.set(b.categoryId, b)
    return expenseCategories.map((c) => {
      const b = byCat.get(c.id)
      const actual = monthExpensesByCategory.get(c.id) ?? 0
      const budgeted = b?.amountCents ?? 0
      const pct = budgeted > 0 ? Math.min(1, actual / budgeted) : 0
      return { category: c, budget: b, actual, budgeted, pct }
    })
  }, [expenseCategories, monthBudgets, monthExpensesByCategory])

  return (
    <div className="grid">
      <section className="card">
        <div className="card__header">
          <div className="card__title">Set Monthly Budget</div>
          <div className="card__hint">Budgets are only allowed for expense categories.</div>
        </div>

        <div className="form">
          <label className="field field--full">
            <div className="field__label">Expense Category</div>
            <select className="select" value={categoryId} onChange={(e) => setCategoryId(e.target.value)}>
              <option value="" disabled>
                Select…
              </option>
              {expenseCategories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>

          <label className="field">
            <div className="field__label">Amount</div>
            <input className="input" inputMode="decimal" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="e.g. 400" />
          </label>

          <div className="form__actions">
            <button
              className="btn"
              onClick={async () => {
                if (!categoryId) return
                const parsed = parseMoneyToCents(amount)
                if (!parsed.ok) return
                const res = await upsertBudget({ month: props.month, categoryId, amountCents: parsed.cents })
                if (!res.ok) return
                setAmount('')
                setCategoryId('')
              }}
            >
              Save budget
            </button>
          </div>
        </div>
      </section>

      <section className="card">
        <div className="card__header">
          <div className="card__title">Budget Utilization</div>
          <div className="card__hint">Budget vs actual expense for the selected month.</div>
        </div>

        <div className="totals" style={{ marginBottom: 12 }}>
          <div className="pill">Total budgeted: {formatCents(totalBudgeted)}</div>
        </div>

        {rows.length === 0 ? (
          <div className="empty">Create an expense category to start budgeting.</div>
        ) : (
          <div className="list">
            {rows.map((r) => {
              const over = r.budgeted > 0 && r.actual > r.budgeted
              return (
                <div key={r.category.id} className="row row--stack">
                  <div className="row__main">
                    <div className="row__title">
                      <span>{r.category.name}</span>
                      <span className={over ? 'badge badge--bad' : 'badge badge--good'}>
                        {r.budgeted === 0 ? 'No budget' : over ? 'Over budget' : 'On track'}
                      </span>
                    </div>

                    <div className="row__sub">
                      Budget: <span className="mono">{formatCents(r.budgeted)}</span> · Actual:{' '}
                      <span className="mono">{formatCents(r.actual)}</span>
                    </div>

                    <div className="progress">
                      <div
                        className={over ? 'progress__bar progress__bar--bad' : 'progress__bar progress__bar--good'}
                        style={{ width: `${Math.min(100, (r.budgeted > 0 ? (r.actual / r.budgeted) * 100 : 0) || 0)}%` }}
                      />
                    </div>
                  </div>

                  <div className="row__actions">
                    {r.budget ? (
                      <button className="btn btn--danger" onClick={() => void deleteBudget(r.budget!.id)}>
                        Delete budget
                      </button>
                    ) : (
                      <button
                        className="btn btn--ghost"
                        onClick={() => {
                          setCategoryId(r.category.id)
                          setAmount(r.budget ? (r.budget.amountCents / 100).toFixed(2) : '')
                        }}
                      >
                        Set budget
                      </button>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </section>
    </div>
  )
}


