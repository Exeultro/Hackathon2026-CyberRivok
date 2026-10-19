import mapsContent from '@content/maps.json'
import textsContent from '@content/texts.json'
import type { MockState, StoredSession, Tally } from './mockTypes'
import { listAchievements, unlockNew } from './achievements'
import { buildLeaderboard } from './leaderboard'
import { localDay, streakOf } from './mockUtil'
import type {
  Achievement,
  Api,
  GameMap,
  KeyStat,
  Lang,
  Leaderboard,
  LeaderboardQuery,
  Mode,
  Progress,
  SessionCreate,
  SessionResult,
  SessionSummary,
  StatsRange,
  StatsSummary,
  TextItem,
  TextQuery,
  TimelinePoint,
  User,
} from './types'

const STORAGE_KEY = 'hero.mock.v1'
const MODES: Mode[] = ['words', 'sentences', 'code', 'symbols']
const WORD_LENGTHS: Record<number, [number, number]> = {
  1: [2, 4],
  2: [3, 5],
  3: [4, 7],
  4: [5, 9],
  5: [6, 10],
}

const SEED: MockState = {
  xp_total: 260,
  stars: { 'journey-1': 3, 'journey-2': 3, 'journey-3': 2, 'journey-4': 1 },
  best_wpm: 0,
  sessions: [],
  keys: {},
  bigrams: {},
  demo: false,
  nickname: 'Герой',
  created_at: new Date().toISOString(),
}

const DAY = 86_400_000

function demoHistory(): Pick<MockState, 'sessions' | 'keys' | 'bigrams'> {
  let a = 12345
  const rnd = () => {
    a = (a * 1664525 + 1013904223) >>> 0
    return a / 4294967296
  }
  const sessions: StoredSession[] = []
  for (let d = 13; d >= 0; d--) {
    if (d === 6) continue
    const runs = d % 3 === 0 ? 2 : 1
    for (let k = 0; k < runs; k++) {
      const progress = 13 - d
      const wpm = 20 + progress * 1.7 + rnd() * 4
      const accuracy = Math.min(99, 88 + progress * 0.55 + rnd() * 2)
      const duration = 60_000 + rnd() * 120_000
      const chars = Math.round(((wpm * 5) / 60_000) * duration)
      const when = new Date()
      when.setDate(when.getDate() - d)
      when.setHours(10 + k * 6, Math.floor(rnd() * 50), 0, 0)
      sessions.push({
        id: crypto.randomUUID(),
        mode: MODES[(progress + k) % 4],
        lang: progress % 2 ? 'en' : 'ru',
        difficulty: 1 + Math.floor(progress / 5),
        started_at: when.toISOString(),
        duration_ms: Math.round(duration),
        wpm: +wpm.toFixed(1),
        cpm: +(wpm * 5).toFixed(1),
        accuracy: +accuracy.toFixed(1),
        errors: Math.round((chars * (100 - accuracy)) / accuracy),
        score: chars * 12,
        xp_gained: Math.round(chars * (accuracy / 100) * 1.2),
        chars,
        max_combo: 14 + progress * 3 + Math.floor(rnd() * 8),
        clean_streak: 28 + progress * 6 + Math.floor(rnd() * 20),
      })
    }
  }
  const keys: Record<string, Tally> = {
    р: { hits: 120, errors: 18 },
    щ: { hits: 40, errors: 9 },
    ъ: { hits: 30, errors: 7 },
    ж: { hits: 60, errors: 8 },
    ь: { hits: 90, errors: 9 },
    ф: { hits: 50, errors: 5 },
    ц: { hits: 45, errors: 6 },
    q: { hits: 35, errors: 7 },
    z: { hits: 28, errors: 5 },
    p: { hits: 60, errors: 6 },
    ';': { hits: 40, errors: 6 },
    '[': { hits: 25, errors: 5 },
    x: { hits: 30, errors: 4 },
  }
  const bigrams: Record<string, Tally> = {
    ть: { hits: 50, errors: 9 },
    ст: { hits: 80, errors: 9 },
    ро: { hits: 70, errors: 8 },
    ск: { hits: 45, errors: 6 },
    пр: { hits: 60, errors: 7 },
    ие: { hits: 55, errors: 5 },
    th: { hits: 90, errors: 10 },
    er: { hits: 70, errors: 6 },
    ng: { hits: 40, errors: 6 },
  }
  return { sessions, keys, bigrams }
}

function load(): MockState {
  let state = structuredClone(SEED)
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) state = { ...state, ...(JSON.parse(raw) as Partial<MockState>) }
  } catch {
    // storage is unavailable, fall back to the seed
  }
  if (!state.demo) {
    const demo = demoHistory()
    state.sessions = [...state.sessions, ...demo.sessions].sort((x, y) => y.started_at.localeCompare(x.started_at))
    for (const [k, v] of Object.entries(demo.keys)) state.keys[k] = merge(state.keys[k], v)
    for (const [k, v] of Object.entries(demo.bigrams)) state.bigrams[k] = merge(state.bigrams[k], v)
    state.demo = true
    save(state)
  }
  if (!state.achievements) {
    unlockNew(state)
    save(state)
  }
  return state
}

function merge(a: Tally | undefined, b: Tally): Tally {
  return { hits: (a?.hits ?? 0) + b.hits, errors: (a?.errors ?? 0) + b.errors }
}

function save(state: MockState) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  } catch {
    // progress just will not persist
  }
}

function levelInfo(xpTotal: number) {
  let level = 1
  let need = 100
  let rest = xpTotal
  while (rest >= need) {
    rest -= need
    level += 1
    need = 100 + 50 * (level - 1)
  }
  return { level, xp_in_level: rest, xp_to_next: need }
}

function progressOf(xpTotal: number): Progress {
  const info = levelInfo(xpTotal)
  return {
    ...info,
    xp_total: xpTotal,
    unlocked_modes: MODES,
    max_difficulty: Math.min(5, 1 + Math.floor(info.level / 3)),
  }
}

function allMaps(): GameMap[] {
  return (mapsContent as unknown as { maps: GameMap[] }).maps
}

function pick<T>(list: T[]): T {
  return list[Math.floor(Math.random() * list.length)]
}

function randomWords(lang: Lang, difficulty: number, n: number): string {
  const pool = textsContent.words[lang]
  const [lo, hi] = WORD_LENGTHS[difficulty] ?? WORD_LENGTHS[1]
  let words = pool.filter((w) => w.length >= lo && w.length <= hi)
  if (words.length < 8) words = pool.filter((w) => w.length >= lo - 1)
  return Array.from({ length: n }, () => pick(words)).join(' ')
}

function pickPassages(lang: Lang, level: number, n: number) {
  const all = textsContent.passages.filter((p) => p.lang === lang)
  const levels = [...new Set(all.map((p) => p.level))]
  const nearest = levels.reduce((best, l) => (Math.abs(l - level) < Math.abs(best - level) ? l : best), levels[0])
  const pool = all.filter((p) => p.level === nearest).sort(() => Math.random() - 0.5)
  return Array.from({ length: n }, (_, i) => pool[i % pool.length])
}

function randomSymbols(lang: Lang, difficulty: number, groups: number): string {
  const sets = textsContent.symbols[lang]
  const chars = sets.slice(0, Math.min(sets.length, difficulty)).join('')
  return Array.from({ length: groups }, () => {
    const len = 3 + Math.floor(Math.random() * 4)
    return Array.from({ length: len }, () => chars[Math.floor(Math.random() * chars.length)]).join('')
  }).join(' ')
}

export const mockApi: Api = {
  async getProgress(): Promise<Progress> {
    return progressOf(load().xp_total)
  },

  async getSummary(): Promise<StatsSummary> {
    const list = load().sessions
    const n = list.length
    const sum = (f: (s: StoredSession) => number) => list.reduce((acc, s) => acc + f(s), 0)
    return {
      sessions_count: n,
      total_time_ms: sum((s) => s.duration_ms),
      total_chars: sum((s) => s.chars),
      best_wpm: Math.max(0, ...list.map((s) => s.wpm)),
      best_cpm: Math.max(0, ...list.map((s) => s.cpm)),
      avg_wpm: n ? +(sum((s) => s.wpm) / n).toFixed(1) : 0,
      avg_accuracy: n ? +(sum((s) => s.accuracy) / n).toFixed(1) : 0,
      streak_days: streakOf(list),
    }
  },

  async getSessions(limit = 10): Promise<SessionSummary[]> {
    return load().sessions.slice(0, limit)
  },

  async getTimeline(range: StatsRange): Promise<TimelinePoint[]> {
    const days = range === 'week' ? 7 : range === 'month' ? 30 : 365
    const from = new Date()
    from.setHours(0, 0, 0, 0)
    from.setDate(from.getDate() - (days - 1))
    const byDay = new Map<string, StoredSession[]>()
    for (const s of load().sessions) {
      if (new Date(s.started_at) < from) continue
      const key = localDay(s.started_at)
      byDay.set(key, [...(byDay.get(key) ?? []), s])
    }
    return [...byDay.entries()]
      .sort(([x], [y]) => x.localeCompare(y))
      .map(([date, list]) => {
        const avg = (f: (s: StoredSession) => number) => +(list.reduce((a, s) => a + f(s), 0) / list.length).toFixed(1)
        return {
          date,
          sessions: list.length,
          avg_wpm: avg((s) => s.wpm),
          avg_cpm: avg((s) => s.cpm),
          avg_accuracy: avg((s) => s.accuracy),
          errors: list.reduce((a, s) => a + s.errors, 0),
        }
      })
  },

  async getKeys(kind, limit = 10): Promise<KeyStat[]> {
    const state = load()
    const source = kind === 'key' ? state.keys : state.bigrams
    return Object.entries(source)
      .filter(([, v]) => v.hits >= 3 && v.errors > 0)
      .map(([key, v]) => ({ key, hits: v.hits, errors: v.errors, error_rate: +(v.errors / v.hits).toFixed(3) }))
      .sort((x, y) => y.error_rate - x.error_rate || y.hits - x.hits)
      .slice(0, limit)
  },

  async getMaps(): Promise<GameMap[]> {
    const state = load()
    const { level } = levelInfo(state.xp_total)
    return allMaps().map((m) => {
      const nodes = m.nodes.map((n) => ({ ...n, stars: state.stars[n.id] ?? 0 }))
      return {
        ...m,
        ready: level >= m.recommended_level,
        stars_total: nodes.reduce((s, n) => s + n.stars, 0),
        stars_max: nodes.length * 3,
        nodes,
      }
    })
  },

  async getAchievements(): Promise<Achievement[]> {
    return listAchievements(load())
  },

  async getLeaderboard(q: LeaderboardQuery): Promise<Leaderboard> {
    const state = load()
    const from = new Date()
    from.setHours(0, 0, 0, 0)
    if (q.period === 'week') from.setDate(from.getDate() - 6)
    const list = q.period === 'all' ? state.sessions : state.sessions.filter((s) => new Date(s.started_at) >= from)
    const mean = list.length >= 3 ? list.reduce((a, s) => a + s.accuracy, 0) / list.length : null
    return buildLeaderboard(q, {
      nickname: state.nickname,
      level: levelInfo(state.xp_total).level,
      speed: list.length ? Math.max(...list.map((s) => s.wpm)) : null,
      accuracy: mean,
      xp: q.period === 'all' ? state.xp_total : list.reduce((a, s) => a + s.xp_gained, 0),
    })
  },

  async getMe(): Promise<User> {
    const state = load()
    return {
      id: 'me',
      email: 'hero@example.com',
      nickname: state.nickname,
      created_at: state.created_at,
      settings: {},
    }
  },

  async updateMe(patch): Promise<User> {
    const state = load()
    const nick = patch.nickname.trim()
    if (nick.length >= 2 && nick.length <= 24) state.nickname = nick
    save(state)
    return mockApi.getMe()
  },

  async reset(): Promise<void> {
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch {
      // nothing to reset
    }
  },

  async getTexts(q: TextQuery): Promise<TextItem[]> {
    if (q.mode === 'sentences') {
      return pickPassages(q.lang, q.difficulty, q.count).map((p, i) => ({
        id: `sentences-${q.lang}-${Date.now()}-${i}`,
        mode: q.mode,
        lang: q.lang,
        difficulty: p.level,
        text: p.text,
        source: p.source,
      }))
    }
    return Array.from({ length: q.count }, (_, i) => {
      let text: string
      if (q.mode === 'words') text = randomWords(q.lang, q.difficulty, 6)
      else if (q.mode === 'code') text = pick(textsContent.code[q.code_language ?? 'python'])
      else text = randomSymbols(q.lang, q.difficulty, 6)
      return {
        id: `${q.mode}-${q.lang}-${Date.now()}-${i}`,
        mode: q.mode,
        lang: q.lang,
        difficulty: q.difficulty,
        code_language: q.mode === 'code' ? (q.code_language ?? 'python') : undefined,
        text,
      }
    })
  },

  async createSession(body): Promise<SessionResult> {
    return finishSession(body)
  },
}

function finishSession(body: SessionCreate): SessionResult {
  const state = load()
  const before = levelInfo(state.xp_total).level

  const minutes = body.duration_ms / 60000
  const cpm = body.chars_correct / minutes
  const wpm = cpm / 5
  const accuracy = body.chars_typed ? (body.chars_correct / body.chars_typed) * 100 : 0
  const xp = Math.round(body.chars_correct * (accuracy / 100) * (1 + body.difficulty * 0.1))
  const score = body.client_score ?? body.chars_correct * 10

  const summary: StoredSession = {
    id: crypto.randomUUID(),
    mode: body.mode,
    lang: body.lang,
    difficulty: body.difficulty,
    started_at: body.started_at,
    duration_ms: body.duration_ms,
    wpm: +wpm.toFixed(1),
    cpm: +cpm.toFixed(1),
    accuracy: +accuracy.toFixed(1),
    errors: body.errors,
    score,
    xp_gained: xp,
    chars: body.chars_correct,
    max_combo: body.max_combo,
    clean_streak: body.longest_clean_streak ?? 0,
  }

  const isBest = wpm > Math.max(state.best_wpm, ...state.sessions.map((s) => s.wpm))
  state.best_wpm = Math.max(state.best_wpm, wpm)
  state.xp_total += xp
  state.sessions = [summary, ...state.sessions].slice(0, 400)
  for (const k of body.key_stats) state.keys[k.key] = merge(state.keys[k.key], k)
  for (const k of body.bigram_stats ?? []) state.bigrams[k.key] = merge(state.bigrams[k.key], k)

  let nodeResult: SessionResult['node_result']
  if (body.node_id) {
    const nodes = allMaps().flatMap((m) => m.nodes.map((n, i) => ({ n, next: m.nodes[i + 1]?.id ?? null })))
    const found = nodes.find((e) => e.n.id === body.node_id)
    if (found) {
      const stars = found.n.star_thresholds.filter((t) => wpm >= t.min_wpm && accuracy >= t.min_accuracy).length
      const best = Math.max(state.stars[found.n.id] ?? 0, stars)
      state.stars[found.n.id] = best
      nodeResult = {
        node_id: found.n.id,
        passed: stars >= 1,
        stars,
        best_stars: best,
        next_node_id: found.next,
      }
    }
  }

  const fresh = unlockNew(state)
  save(state)

  const { chars: _chars, max_combo: _combo, clean_streak: _streak, ...publicSummary } = summary
  void [_chars, _combo, _streak]

  return {
    ...publicSummary,
    level_before: before,
    is_personal_best: isBest,
    progress: progressOf(state.xp_total),
    new_achievements: fresh,
    node_result: nodeResult,
  }
}
