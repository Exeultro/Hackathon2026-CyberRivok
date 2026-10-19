import type { StoredSession } from './mockTypes'

export function localDay(iso: string | Date) {
  const d = typeof iso === 'string' ? new Date(iso) : iso
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

export function streakOf(sessions: StoredSession[]) {
  const days = new Set(sessions.map((s) => localDay(s.started_at)))
  const cursor = new Date()
  if (!days.has(localDay(cursor))) cursor.setDate(cursor.getDate() - 1)
  let n = 0
  while (days.has(localDay(cursor))) {
    n += 1
    cursor.setDate(cursor.getDate() - 1)
  }
  return n
}
