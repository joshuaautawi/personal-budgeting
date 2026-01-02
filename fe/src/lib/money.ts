import type { Cents } from './types'

export function formatCents(cents: Cents): string {
  const amount = cents / 100
  return amount.toLocaleString('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 2,
    minimumFractionDigits: 2,
  })
}

export function parseMoneyToCents(raw: string): { ok: true; cents: Cents } | { ok: false; error: string } {
  const trimmed = raw.trim()
  if (!trimmed) return { ok: false, error: 'Amount is required.' }
  // Accept common inputs:
  // - "10000"
  // - "10.000" (thousands separators)
  // - "10.000,50" (id-ID style decimal)
  // - "10000.50" (en style decimal)
  const normalized = trimmed
    .replace(/[Rr][Pp]/g, '') // optional "Rp"
    .replace(/\s/g, '')
    .replace(/\./g, '') // remove thousand separators
    .replace(/,/g, '.') // allow comma decimal

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


