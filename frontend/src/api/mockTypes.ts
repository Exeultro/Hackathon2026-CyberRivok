import type { SessionSummary } from './types'

export interface StoredSession extends SessionSummary {
  chars: number
  max_combo?: number
  clean_streak?: number
}

export interface Tally {
  hits: number
  errors: number
}

export interface MockState {
  xp_total: number
  stars: Record<string, number>
  best_wpm: number
  sessions: StoredSession[]
  keys: Record<string, Tally>
  bigrams: Record<string, Tally>
  demo: boolean
  nickname: string
  created_at: string
  achievements?: Record<string, string>
}
