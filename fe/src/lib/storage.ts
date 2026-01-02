import type { AppStateV1 } from './types'

const STORAGE_KEY = 'pb.appState.v1'

export function loadState(): AppStateV1 | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as unknown
    if (!parsed || typeof parsed !== 'object') return null
    // Very lightweight runtime check; we keep it forgiving for local storage.
    const state = parsed as AppStateV1
    if (state.version !== 1) return null
    if (!Array.isArray(state.categories) || !Array.isArray(state.budgets) || !Array.isArray(state.transactions)) {
      return null
    }
    return state
  } catch {
    return null
  }
}

export function saveState(state: AppStateV1): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
}


