import type { AppStateV1, Budget, Category, Id, Txn } from './types'

type ApiError = {
  status: number
  code: string
}

const API_BASE = (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? ''

async function readJson<T>(res: Response): Promise<T> {
  const text = await res.text()
  if (!text) return {} as T
  return JSON.parse(text) as T
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  })

  if (!res.ok) {
    const body = await readJson<{ error?: string }>(res).catch(() => ({}) as { error?: string })
    const code = body.error ?? `http_${res.status}`
    throw { status: res.status, code } satisfies ApiError
  }
  return await readJson<T>(res)
}

export function getState(): Promise<AppStateV1> {
  return request<AppStateV1>('/api/v1/state')
}

export function createCategory(input: { type: 'income' | 'expense'; name: string; description?: string }): Promise<Category> {
  return request<Category>('/api/v1/categories', { method: 'POST', body: JSON.stringify(input) })
}

export function updateCategory(id: Id, patch: { name?: string; description?: string }): Promise<Category> {
  return request<Category>(`/api/v1/categories/${id}`, { method: 'PATCH', body: JSON.stringify(patch) })
}

export async function deleteCategory(id: Id): Promise<void> {
  await request<unknown>(`/api/v1/categories/${id}`, { method: 'DELETE' })
}

export function upsertBudget(input: { month: string; categoryId: Id; amountCents: number }): Promise<Budget> {
  return request<Budget>('/api/v1/budgets', { method: 'PUT', body: JSON.stringify(input) })
}

export async function deleteBudget(id: Id): Promise<void> {
  await request<unknown>(`/api/v1/budgets/${id}`, { method: 'DELETE' })
}

export function createTxn(input: {
  kind: 'income' | 'expense'
  date: string
  categoryId: Id
  amountCents: number
  note?: string
}): Promise<Txn> {
  return request<Txn>('/api/v1/transactions', { method: 'POST', body: JSON.stringify(input) })
}

export function updateTxn(
  id: Id,
  patch: { kind?: 'income' | 'expense'; date?: string; categoryId?: Id; amountCents?: number; note?: string },
): Promise<Txn> {
  return request<Txn>(`/api/v1/transactions/${id}`, { method: 'PATCH', body: JSON.stringify(patch) })
}

export async function deleteTxn(id: Id): Promise<void> {
  await request<unknown>(`/api/v1/transactions/${id}`, { method: 'DELETE' })
}

export function prettyApiError(e: unknown): string {
  const err = e as Partial<ApiError>
  const code = err.code ?? 'unknown'
  switch (code) {
    case 'validation':
      return 'Validation error. Please check your input.'
    case 'conflict':
      return 'Conflict. This item is in use or violates a rule.'
    case 'not_found':
      return 'Not found. It may have been deleted already.'
    default:
      return 'Request failed. Please try again.'
  }
}


