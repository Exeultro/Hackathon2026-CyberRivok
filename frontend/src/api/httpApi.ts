import type { components } from './schema'
import type { Api, LeaderboardQuery, SessionCreate, TextQuery } from './types'

type AuthResponse = components['schemas']['AuthResponse']

const BASE = (import.meta.env.VITE_API_URL as string).replace(/\/$/, '')
const KEY = 'hero.http.v1'

interface Stored {
  email: string
  password: string
  access?: string
  refresh?: string
}

function load(): Stored | null {
  try {
    return JSON.parse(localStorage.getItem(KEY) ?? 'null')
  } catch {
    return null
  }
}

function save(s: Stored) {
  localStorage.setItem(KEY, JSON.stringify(s))
}

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
    public details?: Record<string, unknown>,
  ) {
    super(message)
  }
}

function rnd(n: number): string {
  const a = new Uint8Array(n)
  crypto.getRandomValues(a)
  return Array.from(a, (b) => b.toString(16).padStart(2, '0')).join('')
}

async function raw(method: string, path: string, body?: unknown, token?: string): Promise<Response> {
  const headers: Record<string, string> = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  if (token) headers.Authorization = `Bearer ${token}`
  return fetch(BASE + path, { method, headers, body: body === undefined ? undefined : JSON.stringify(body) })
}

async function fail(res: Response): Promise<never> {
  const e = await res.json().catch(() => ({}))
  throw new ApiError(res.status, e.code ?? 'error', e.message ?? res.statusText, e.details)
}

function applyAuth(s: Stored, r: AuthResponse): Stored {
  const next = { ...s, access: r.tokens.access_token, refresh: r.tokens.refresh_token }
  save(next)
  return next
}

// No login screen yet: the browser gets a guest account on first run.
async function register(): Promise<Stored> {
  const id = rnd(8)
  const s: Stored = { email: `guest-${id}@hero.local`, password: rnd(16) }
  const res = await raw('POST', '/auth/register', { email: s.email, password: s.password, nickname: `Герой-${id.slice(0, 6)}` })
  if (!res.ok) return fail(res)
  return applyAuth(s, await res.json())
}

let session: Promise<Stored> | null = null

function ensureSession(): Promise<Stored> {
  session ??= (async () => load() ?? register())()
  return session
}

async function reauth(s: Stored): Promise<Stored> {
  if (s.refresh) {
    const res = await raw('POST', '/auth/refresh', { refresh_token: s.refresh })
    if (res.ok) {
      const t = await res.json()
      const next = { ...s, access: t.access_token, refresh: t.refresh_token }
      save(next)
      return next
    }
  }
  const res = await raw('POST', '/auth/login', { email: s.email, password: s.password })
  if (!res.ok) return fail(res)
  return applyAuth(s, await res.json())
}

async function call<T>(method: string, path: string, body?: unknown): Promise<T> {
  let s = await ensureSession()
  let res = await raw(method, path, body, s.access)
  if (res.status === 401) {
    session = reauth(s)
    s = await session
    res = await raw(method, path, body, s.access)
  }
  if (!res.ok) return fail(res)
  return res.status === 204 ? (undefined as T) : res.json()
}

function qs(params: Record<string, string | number | undefined>): string {
  const u = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) if (v !== undefined) u.set(k, String(v))
  return `?${u}`
}

export const httpApi: Api = {
  async getMaps() {
    return (await call<{ items: Awaited<ReturnType<Api['getMaps']>> }>('GET', '/maps')).items
  },
  getProgress: () => call('GET', '/me/progress'),
  getSummary: () => call('GET', '/me/stats/summary'),
  async getSessions(limit = 20) {
    return (await call<{ items: Awaited<ReturnType<Api['getSessions']>> }>('GET', `/sessions${qs({ limit })}`)).items
  },
  async getTimeline(range) {
    return (await call<{ points: Awaited<ReturnType<Api['getTimeline']>> }>('GET', `/me/stats/timeline${qs({ range })}`)).points
  },
  async getKeys(kind, limit = 10) {
    return (await call<{ items: Awaited<ReturnType<Api['getKeys']>> }>('GET', `/me/stats/keys${qs({ kind, limit })}`)).items
  },
  async getAchievements() {
    return (await call<{ items: Awaited<ReturnType<Api['getAchievements']>> }>('GET', '/achievements')).items
  },
  getLeaderboard(q: LeaderboardQuery) {
    return call('GET', `/leaderboard${qs({ metric: q.metric, period: q.period, limit: q.limit })}`)
  },
  getMe: () => call('GET', '/me'),
  updateMe: (patch) => call('PATCH', '/me', patch),
  async reset() {
    await call('DELETE', '/me/progress')
  },
  async getTexts(q: TextQuery) {
    const r = await call<{ items: Awaited<ReturnType<Api['getTexts']>> }>(
      'GET',
      `/texts${qs({ mode: q.mode, lang: q.lang, difficulty: q.difficulty, count: q.count, code_language: q.code_language })}`,
    )
    return r.items
  },
  createSession: (body: SessionCreate) => call('POST', '/sessions', body),
}
