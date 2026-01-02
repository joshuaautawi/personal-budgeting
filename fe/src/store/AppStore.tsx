import React, { createContext, useContext, useEffect, useMemo, useReducer } from 'react'
import type { AppStateV1, Budget, Category, CategoryType, Id, Txn, TransactionKind } from '../lib/types'
import * as apiClient from '../lib/api'

type StoreError = { message: string }

type Result = { ok: true } | { ok: false; message: string }

type AppStore = {
  state: AppStateV1
  lastError: StoreError | null
  clearError: () => void

  addCategory: (input: { type: CategoryType; name: string; description?: string }) => Promise<Result>
  updateCategory: (id: Id, patch: { name: string; description?: string }) => Promise<Result>
  deleteCategory: (id: Id) => Promise<Result>

  upsertBudget: (input: { month: string; categoryId: Id; amountCents: number }) => Promise<Result>
  deleteBudget: (id: Id) => Promise<Result>

  addTxn: (input: { kind: TransactionKind; date: string; categoryId: Id; amountCents: number; note?: string }) => Promise<Result>
  updateTxn: (
    id: Id,
    patch: { date: string; categoryId: Id; amountCents: number; note?: string; kind: TransactionKind },
  ) => Promise<Result>
  deleteTxn: (id: Id) => Promise<Result>
}

type Action =
  | { type: 'error'; message: string }
  | { type: 'clearError' }
  | { type: 'setState'; state: AppStateV1 }
  | { type: 'addCategory'; category: Category }
  | { type: 'updateCategory'; id: Id; patch: { name: string; description?: string; updatedAt: string } }
  | { type: 'deleteCategory'; id: Id }
  | { type: 'upsertBudget'; budget: Budget }
  | { type: 'deleteBudget'; id: Id }
  | { type: 'addTxn'; txn: Txn }
  | { type: 'updateTxn'; id: Id; patch: Omit<Txn, 'id' | 'createdAt'> }
  | { type: 'deleteTxn'; id: Id }

type ReducerState = {
  app: AppStateV1
  lastError: StoreError | null
}

function emptyState(): AppStateV1 {
  return { version: 1, categories: [], budgets: [], transactions: [] }
}

function reducer(s: ReducerState, a: Action): ReducerState {
  switch (a.type) {
    case 'error':
      return { ...s, lastError: { message: a.message } }
    case 'clearError':
      return { ...s, lastError: null }
    case 'setState':
      return { app: a.state, lastError: null }
    case 'addCategory':
      return { ...s, app: { ...s.app, categories: [a.category, ...s.app.categories] } }
    case 'updateCategory':
      return {
        ...s,
        app: {
          ...s.app,
          categories: s.app.categories.map((c) => (c.id === a.id ? { ...c, ...a.patch } : c)),
        },
      }
    case 'deleteCategory':
      return { ...s, app: { ...s.app, categories: s.app.categories.filter((c) => c.id !== a.id) } }
    case 'upsertBudget': {
      const idx = s.app.budgets.findIndex((b) => b.id === a.budget.id)
      if (idx >= 0) {
        const next = [...s.app.budgets]
        next[idx] = a.budget
        return { ...s, app: { ...s.app, budgets: next } }
      }
      return { ...s, app: { ...s.app, budgets: [a.budget, ...s.app.budgets] } }
    }
    case 'deleteBudget':
      return { ...s, app: { ...s.app, budgets: s.app.budgets.filter((b) => b.id !== a.id) } }
    case 'addTxn':
      return { ...s, app: { ...s.app, transactions: [a.txn, ...s.app.transactions] } }
    case 'updateTxn':
      return {
        ...s,
        app: {
          ...s.app,
          transactions: s.app.transactions.map((t) => (t.id === a.id ? { ...t, ...a.patch } : t)),
        },
      }
    case 'deleteTxn':
      return { ...s, app: { ...s.app, transactions: s.app.transactions.filter((t) => t.id !== a.id) } }
    default:
      return s
  }
}

const Ctx = createContext<AppStore | null>(null)

export function AppStoreProvider(props: { children: React.ReactNode }) {
  const [s, dispatch] = useReducer(reducer, { app: emptyState(), lastError: null })

  // Load initial state from backend.
  useEffect(() => {
    let cancelled = false
    apiClient
      .getState()
      .then((state: AppStateV1) => {
        if (cancelled) return
        dispatch({ type: 'setState', state })
      })
      .catch((e: unknown) => {
        if (cancelled) return
        dispatch({ type: 'error', message: apiClient.prettyApiError(e) })
      })
    return () => {
      cancelled = true
    }
  }, [])

  const store: AppStore = useMemo(() => {
    const clearError = () => dispatch({ type: 'clearError' })

    const setError = (message: string) => dispatch({ type: 'error', message })

    const ok: Result = { ok: true }
    const fail = (message: string): Result => ({ ok: false, message })

    const addCategory: AppStore['addCategory'] = async ({ type, name, description }) => {
      try {
        const created = await apiClient.createCategory({ type, name, description })
        dispatch({ type: 'addCategory', category: created })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const updateCategory: AppStore['updateCategory'] = async (id, patch) => {
      try {
        const updated = await apiClient.updateCategory(id, { name: patch.name, description: patch.description })
        dispatch({
          type: 'updateCategory',
          id,
          patch: { name: updated.name, description: updated.description || undefined, updatedAt: updated.updatedAt },
        })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const deleteCategory: AppStore['deleteCategory'] = async (id) => {
      try {
        await apiClient.deleteCategory(id)
        dispatch({ type: 'deleteCategory', id })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const upsertBudget: AppStore['upsertBudget'] = async ({ month, categoryId, amountCents }) => {
      try {
        const budget = await apiClient.upsertBudget({ month, categoryId, amountCents })
        dispatch({ type: 'upsertBudget', budget })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const deleteBudget: AppStore['deleteBudget'] = async (id) => {
      try {
        await apiClient.deleteBudget(id)
        dispatch({ type: 'deleteBudget', id })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const addTxn: AppStore['addTxn'] = async ({ kind, date, categoryId, amountCents, note }) => {
      try {
        const txn = await apiClient.createTxn({ kind, date, categoryId, amountCents, note })
        dispatch({ type: 'addTxn', txn })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const updateTxn: AppStore['updateTxn'] = async (id, patch) => {
      try {
        const updated = await apiClient.updateTxn(id, patch)
        dispatch({
          type: 'updateTxn',
          id,
          patch: {
            kind: updated.kind,
            date: updated.date,
            categoryId: updated.categoryId,
            amountCents: updated.amountCents,
            note: updated.note || undefined,
            updatedAt: updated.updatedAt,
          },
        })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    const deleteTxn: AppStore['deleteTxn'] = async (id) => {
      try {
        await apiClient.deleteTxn(id)
        dispatch({ type: 'deleteTxn', id })
        return ok
      } catch (e) {
        const msg = apiClient.prettyApiError(e)
        setError(msg)
        return fail(msg)
      }
    }

    return {
      state: s.app,
      lastError: s.lastError,
      clearError,
      addCategory,
      updateCategory,
      deleteCategory,
      upsertBudget,
      deleteBudget,
      addTxn,
      updateTxn,
      deleteTxn,
    }
  }, [s.app, s.lastError])

  return <Ctx.Provider value={store}>{props.children}</Ctx.Provider>
}

export function useAppStore(): AppStore {
  const ctx = useContext(Ctx)
  if (!ctx) throw new Error('useAppStore must be used within AppStoreProvider.')
  return ctx
}


