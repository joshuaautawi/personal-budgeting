import { useEffect, useMemo, useState } from 'react'
import { formatCents, parseMoneyToCents } from '../lib/money'
import { isInMonth } from '../lib/month'
import type { Id, TransactionKind } from '../lib/types'
import { Modal } from '../components/Modal'
import { useAppStore } from '../store/AppStore'

type KindFilter = 'all' | TransactionKind

export function Transactions(props: { month: string }) {
  const { state, addTxn, deleteTxn, updateTxn } = useAppStore()

  const [addOpen, setAddOpen] = useState(false)
  const [addError, setAddError] = useState<string | null>(null)

  const [kind, setKind] = useState<TransactionKind>('expense')
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [amount, setAmount] = useState('')
  const [categoryId, setCategoryId] = useState<Id>('')
  const [note, setNote] = useState('')

  const [filterKind, setFilterKind] = useState<KindFilter>('all')
  const [filterCategoryId, setFilterCategoryId] = useState<Id>('all' as Id)
  const [fromDate, setFromDate] = useState('')
  const [toDate, setToDate] = useState('')

  const [editingId, setEditingId] = useState<Id | null>(null)
  const [editKind, setEditKind] = useState<TransactionKind>('expense')
  const [editDate, setEditDate] = useState('')
  const [editAmount, setEditAmount] = useState('')
  const [editCategoryId, setEditCategoryId] = useState<Id>('')
  const [editNote, setEditNote] = useState('')

  const categoriesForKind = useMemo(
    () => state.categories.filter((c) => c.type === kind),
    [state.categories, kind],
  )

  useEffect(() => {
    // Keep default aligned to "today" (even if user is viewing a different month).
    const today = new Date().toISOString().slice(0, 10)
    if (date !== today) setDate(today)
  }, [props.month, date])

  const allCategories = state.categories
  const monthTxns = useMemo(() => state.transactions.filter((t) => isInMonth(t.date, props.month)), [state.transactions, props.month])

  const filteredTxns = useMemo(() => {
    return monthTxns
      .filter((t) => (filterKind === 'all' ? true : t.kind === filterKind))
      .filter((t) => (filterCategoryId === ('all' as Id) ? true : t.categoryId === filterCategoryId))
      .filter((t) => (fromDate ? t.date >= fromDate : true))
      .filter((t) => (toDate ? t.date <= toDate : true))
  }, [monthTxns, filterKind, filterCategoryId, fromDate, toDate])

  const totals = useMemo(() => {
    const income = filteredTxns.filter((t) => t.kind === 'income').reduce((s, t) => s + t.amountCents, 0)
    const expense = filteredTxns.filter((t) => t.kind === 'expense').reduce((s, t) => s + t.amountCents, 0)
    return { income, expense, net: income - expense }
  }, [filteredTxns])

  const beginEdit = (id: Id) => {
    const t = state.transactions.find((x) => x.id === id)
    if (!t) return
    setEditingId(id)
    setEditKind(t.kind)
    setEditDate(t.date)
    setEditAmount((t.amountCents / 100).toFixed(2))
    setEditCategoryId(t.categoryId)
    setEditNote(t.note ?? '')
  }

  const saveEdit = () => {
    if (!editingId) return
    const parsed = parseMoneyToCents(editAmount)
    if (!parsed.ok) return
    void updateTxn(editingId, {
      kind: editKind,
      date: editDate,
      categoryId: editCategoryId,
      amountCents: parsed.cents,
      note: editNote,
    }).then((res) => {
      if (!res.ok) return
      setEditingId(null)
    })
  }

  return (
    <div className="grid">
      <section className="card card--hero">
        <div className="card__header">
          <div className="card__title">Transactions</div>
          <div className="card__hint">Filter by month (top), date range, and category.</div>
          <div style={{ flex: 1 }} />
          <button
            className="btn"
            onClick={() => {
              setAddError(null)
              setKind('expense')
              setDate(new Date().toISOString().slice(0, 10))
              setAmount('')
              setCategoryId('' as Id)
              setNote('')
              setAddOpen(true)
            }}
          >
            Add transaction
          </button>
        </div>

        <div className="filters">
          <label className="field">
            <div className="field__label">Type</div>
            <select className="select" value={filterKind} onChange={(e) => setFilterKind(e.target.value as KindFilter)}>
              <option value="all">All</option>
              <option value="income">Income</option>
              <option value="expense">Expense</option>
            </select>
          </label>

          <label className="field field--full">
            <div className="field__label">Category</div>
            <select className="select" value={filterCategoryId} onChange={(e) => setFilterCategoryId(e.target.value as Id)}>
              <option value="all">All categories</option>
              {allCategories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.type === 'income' ? 'Income' : 'Expense'} — {c.name}
                </option>
              ))}
            </select>
          </label>

          <label className="field">
            <div className="field__label">From</div>
            <input className="input" type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)} />
          </label>

          <label className="field">
            <div className="field__label">To</div>
            <input className="input" type="date" value={toDate} onChange={(e) => setToDate(e.target.value)} />
          </label>
        </div>

        <div className="totals">
          <div className="pill pill--good">Income: {formatCents(totals.income)}</div>
          <div className="pill pill--bad">Expenses: {formatCents(totals.expense)}</div>
          <div className="pill">Net: {formatCents(totals.net)}</div>
        </div>

        {filteredTxns.length === 0 ? (
          <div className="empty">No transactions match these filters for this month.</div>
        ) : (
          <div className="list">
            {filteredTxns.map((t) => {
              const cat = state.categories.find((c) => c.id === t.categoryId)
              const isEditing = editingId === t.id
              return (
                <div key={t.id} className="row">
                  <div className="row__main">
                    <div className="row__title">
                      <span className={t.kind === 'income' ? 'badge badge--good' : 'badge badge--bad'}>
                        {t.kind === 'income' ? 'Income' : 'Expense'}
                      </span>
                      <span className="mono">{t.date}</span>
                      <span className="row__amount">{formatCents(t.amountCents)}</span>
                    </div>
                    <div className="row__sub">
                      {cat ? cat.name : 'Unknown category'}
                      {t.note ? ` — ${t.note}` : ''}
                    </div>

                    {isEditing ? (
                      <div className="row__edit">
                        <div className="form form--compact">
                          <label className="field">
                            <div className="field__label">Type</div>
                            <select
                              className="select"
                              value={editKind}
                              onChange={(e) => {
                                const next = e.target.value as TransactionKind
                                setEditKind(next)
                                setEditCategoryId('')
                              }}
                            >
                              <option value="expense">Expense</option>
                              <option value="income">Income</option>
                            </select>
                          </label>

                          <label className="field">
                            <div className="field__label">Date</div>
                            <input className="input" type="date" value={editDate} onChange={(e) => setEditDate(e.target.value)} />
                          </label>

                          <label className="field">
                            <div className="field__label">Amount</div>
                            <input className="input" inputMode="decimal" value={editAmount} onChange={(e) => setEditAmount(e.target.value)} />
                          </label>

                          <label className="field field--full">
                            <div className="field__label">Category</div>
                            <select className="select" value={editCategoryId} onChange={(e) => setEditCategoryId(e.target.value)}>
                              <option value="" disabled>
                                Select…
                              </option>
                              {state.categories
                                .filter((c) => c.type === editKind)
                                .map((c) => (
                                  <option key={c.id} value={c.id}>
                                    {c.name}
                                  </option>
                                ))}
                            </select>
                          </label>

                          <label className="field field--full">
                            <div className="field__label">Note</div>
                            <input className="input" value={editNote} onChange={(e) => setEditNote(e.target.value)} />
                          </label>

                          <div className="form__actions">
                            <button className="btn btn--ghost" onClick={() => setEditingId(null)}>
                              Cancel
                            </button>
                            <button className="btn" onClick={saveEdit}>
                              Save
                            </button>
                          </div>
                        </div>
                      </div>
                    ) : null}
                  </div>

                  <div className="row__actions">
                    <button className="btn btn--ghost" onClick={() => beginEdit(t.id)}>
                      Edit
                    </button>
                    <button className="btn btn--danger" onClick={() => void deleteTxn(t.id)}>
                      Delete
                    </button>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </section>

      <Modal
        open={addOpen}
        title="Add transaction"
        onClose={() => {
          setAddOpen(false)
          setAddError(null)
        }}
      >
        {addError ? <div className="inlineAlert">{addError}</div> : null}

        {state.categories.length === 0 ? (
          <div className="empty">Create at least one category first (Categories tab).</div>
        ) : null}

        <div className="form">
          <label className="field">
            <div className="field__label">Type</div>
            <select
              className="select"
              value={kind}
              onChange={(e) => {
                const next = e.target.value as TransactionKind
                setKind(next)
                setCategoryId('' as Id)
              }}
            >
              <option value="expense">Expense</option>
              <option value="income">Income</option>
            </select>
          </label>

          <label className="field">
            <div className="field__label">Date</div>
            <input className="input" type="date" value={date} onChange={(e) => setDate(e.target.value)} />
          </label>

          <label className="field">
            <div className="field__label">Amount</div>
            <input
              className="input"
              inputMode="decimal"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="e.g. 23.50"
            />
          </label>

          <label className="field field--full">
            <div className="field__label">Category</div>
            <select className="select" value={categoryId} onChange={(e) => setCategoryId(e.target.value as Id)}>
              <option value="" disabled>
                Select…
              </option>
              {categoriesForKind.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
            {categoriesForKind.length === 0 ? (
              <div className="field__help">No {kind} categories yet. Create one in Categories.</div>
            ) : null}
          </label>

          <label className="field field--full">
            <div className="field__label">Note (optional)</div>
            <input className="input" value={note} onChange={(e) => setNote(e.target.value)} placeholder="Optional description" />
          </label>

          <div className="form__actions">
            <button
              className="btn btn--ghost"
              onClick={() => {
                setAddOpen(false)
                setAddError(null)
              }}
            >
              Cancel
            </button>
            <button
              className="btn"
              onClick={async () => {
                setAddError(null)
                const parsed = parseMoneyToCents(amount)
                if (!parsed.ok) return setAddError(parsed.error)
                if (!categoryId) return setAddError('Select a category.')
                const res = await addTxn({ kind, date, categoryId, amountCents: parsed.cents, note })
                if (!res.ok) return
                setAddOpen(false)
              }}
            >
              Add
            </button>
          </div>
        </div>
      </Modal>
    </div>
  )
}


