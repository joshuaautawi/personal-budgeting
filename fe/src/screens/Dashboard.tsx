import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { formatCents } from '../lib/money'
import { isInMonth, monthKeyToLabel } from '../lib/month'
import { useAppStore } from '../store/AppStore'

export function Dashboard(props: { month: string }) {
  const { state } = useAppStore()

  const txns = state.transactions.filter((t) => isInMonth(t.date, props.month))
  const incomeCents = txns.filter((t) => t.kind === 'income').reduce((sum, t) => sum + t.amountCents, 0)
  const expenseCents = txns.filter((t) => t.kind === 'expense').reduce((sum, t) => sum + t.amountCents, 0)
  const netCents = incomeCents - expenseCents

  const chartData = [
    {
      name: monthKeyToLabel(props.month),
      Income: incomeCents / 100,
      Expenses: expenseCents / 100,
    },
  ]

  return (
    <div className="grid">
      <section className="card card--hero">
        <div className="hero__title">Monthly Summary</div>
        <div className="hero__subtitle">{monthKeyToLabel(props.month)}</div>

        <div className="stats">
          <div className="stat">
            <div className="stat__label">Income</div>
            <div className="stat__value stat__value--good">{formatCents(incomeCents)}</div>
          </div>
          <div className="stat">
            <div className="stat__label">Expenses</div>
            <div className="stat__value stat__value--bad">{formatCents(expenseCents)}</div>
          </div>
          <div className="stat">
            <div className="stat__label">Net</div>
            <div className="stat__value">{formatCents(netCents)}</div>
          </div>
        </div>
      </section>

      <section className="card">
        <div className="card__header">
          <div className="card__title">Income vs Expenses</div>
          <div className="card__hint">Updates automatically when you add/edit/delete transactions.</div>
        </div>

        <div className="chart">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={chartData} margin={{ left: 8, right: 16, top: 8, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" opacity={0.25} />
              <XAxis dataKey="name" tickLine={false} axisLine={false} />
              <YAxis tickLine={false} axisLine={false} />
              <Tooltip formatter={(value) => formatCents(Math.round(Number(value) * 100))} />
              <Bar dataKey="Income" fill="var(--good)" radius={[10, 10, 10, 10]} />
              <Bar dataKey="Expenses" fill="var(--bad)" radius={[10, 10, 10, 10]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </section>
    </div>
  )
}


