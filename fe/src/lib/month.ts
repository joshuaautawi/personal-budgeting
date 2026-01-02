import { format, isValid, parse } from 'date-fns'
import type { DateKey, MonthKey } from './types'

export function nowIso(): string {
  return new Date().toISOString()
}

export function getCurrentMonthKey(): MonthKey {
  return format(new Date(), 'yyyy-MM')
}

export function monthKeyToLabel(month: MonthKey): string {
  // month is yyyy-MM
  const dt = parse(month + '-01', 'yyyy-MM-dd', new Date())
  if (!isValid(dt)) return month
  return format(dt, 'MMM yyyy')
}

export function isDateKey(value: string): value is DateKey {
  // Basic shape check; input type="date" will produce yyyy-MM-dd
  return /^\d{4}-\d{2}-\d{2}$/.test(value)
}

export function isInMonth(date: DateKey, month: MonthKey): boolean {
  return date.startsWith(month + '-')
}


