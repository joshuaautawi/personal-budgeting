import type { Cents } from './types'

export function formatCents(cents: Cents): string {
  const dollars = cents / 100
  return dollars.toLocaleString(undefined, {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2,
    minimumFractionDigits: 2,
  })
}

export function parseMoneyToCents(raw: string): { ok: true; cents: Cents } | { ok: false; error: string } {
  const trimmed = raw.trim()
  if (!trimmed) return { ok: false, error: 'Amount is required.' }
  const normalized = trimmed.replace(/[$,\s]/g, '')

  // Allow: 10, 10.5, 10.50, .5
  if (!/^(\d+(\.\d{0,2})?|\.\d{1,2})$/.test(normalized)) {
    return { ok: false, error: 'Enter a valid amount (up to 2 decimals).' }
  }

  const num = Number(normalized)
  if (!Number.isFinite(num)) return { ok: false, error: 'Enter a valid amount.' }
  if (num < 0) return { ok: false, error: 'Amount cannot be negative.' }

  const cents = Math.round(num * 100)
  return { ok: true, cents }
}


