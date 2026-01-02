import { useMemo, useState } from 'react'
import { getCurrentMonthKey, monthKeyToLabel } from './lib/month'
import { useAppStore } from './store/AppStore'
import { Dashboard } from './screens/Dashboard'
import { Transactions } from './screens/Transactions'
import { Budgets } from './screens/Budgets'
import { Categories } from './screens/Categories'

type TabKey = 'dashboard' | 'transactions' | 'budgets' | 'categories'

export default function App() {
  const { lastError, clearError } = useAppStore()
  const [tab, setTab] = useState<TabKey>('dashboard')
  const [month, setMonth] = useState(getCurrentMonthKey())

  const title = useMemo(() => monthKeyToLabel(month), [month])

  return (
    <div className="app">
      <header className="topbar">
        <div className="topbar__left">
          <div className="brand">
            <div className="brand__name">Personal Budget</div>
            <div className="brand__sub">Minimal budgeting & tracking</div>
          </div>
        </div>

        <div className="topbar__right">
          <label className="field">
            <div className="field__label">Month</div>
            <input
              className="input"
              type="month"
              value={month}
              onChange={(e) => setMonth(e.target.value)}
              aria-label="Select month"
            />
          </label>
        </div>
      </header>

      <nav className="tabs" aria-label="Primary navigation">
        <button className={tab === 'dashboard' ? 'tab tab--active' : 'tab'} onClick={() => setTab('dashboard')}>
          Dashboard
        </button>
        <button className={tab === 'transactions' ? 'tab tab--active' : 'tab'} onClick={() => setTab('transactions')}>
          Transactions
        </button>
        <button className={tab === 'budgets' ? 'tab tab--active' : 'tab'} onClick={() => setTab('budgets')}>
          Budgets
        </button>
        <button className={tab === 'categories' ? 'tab tab--active' : 'tab'} onClick={() => setTab('categories')}>
          Categories
        </button>

        <div className="tabs__spacer" />
        <div className="tabs__meta">{title}</div>
      </nav>

      {lastError ? (
        <div className="toast" role="alert">
          <div className="toast__msg">{lastError.message}</div>
          <button className="btn btn--ghost" onClick={clearError}>
            Dismiss
          </button>
        </div>
      ) : null}

      <main className="main">
        {tab === 'dashboard' ? <Dashboard month={month} /> : null}
        {tab === 'transactions' ? <Transactions month={month} /> : null}
        {tab === 'budgets' ? <Budgets month={month} /> : null}
        {tab === 'categories' ? <Categories /> : null}
      </main>

      <footer className="footer">
        <div className="footer__hint">
          Data is stored in your backend database (Postgres via the Go API).
        </div>
      </footer>
    </div>
  )
}
