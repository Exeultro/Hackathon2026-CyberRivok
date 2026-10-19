import type { Leaderboard, LeaderboardQuery } from './types'

export interface MeStats {
  nickname: string
  level: number
  speed: number | null
  accuracy: number | null
  xp: number
}

const NAMES = [
  'Рунный', 'Ворон', 'Мираж', 'Сокол', 'Атаман', 'Тень', 'Гром', 'Лис', 'Ветер', 'Кремень',
  'Буря', 'Странник', 'Искра', 'Бродяга', 'Орион', 'Дозор', 'Полынь', 'Клинок', 'Сталь', 'Заря',
  'Рысь', 'Дымка', 'Лазутчик', 'Янтарь',
]

interface Rival {
  id: string
  nickname: string
  level: number
  speed: number
  accuracy: number
  xp: number
  todayF: number
  weekF: number
}

function rivals(): Rival[] {
  let a = 987654
  const rnd = () => {
    a = (a * 1664525 + 1013904223) >>> 0
    return a / 4294967296
  }
  return NAMES.map((nickname, i) => {
    const level = 3 + Math.floor(rnd() * 18)
    return {
      id: `rival-${i}`,
      nickname,
      level,
      speed: 34 + rnd() * 62,
      accuracy: 91 + rnd() * 8.5,
      xp: level * level * 55 + rnd() * 400,
      todayF: 0.78 + rnd() * 0.2,
      weekF: 0.9 + rnd() * 0.1,
    }
  })
}

const RIVALS = rivals()

function valueOf(r: Rival, q: LeaderboardQuery): number {
  const f = q.period === 'today' ? r.todayF : q.period === 'week' ? r.weekF : 1
  if (q.metric === 'speed') return r.speed * f
  if (q.metric === 'accuracy') return Math.min(100, r.accuracy - (q.period === 'today' ? 0.8 : 0))
  const share = q.period === 'today' ? 0.07 * (0.6 + r.todayF) : q.period === 'week' ? 0.3 * (0.8 + r.weekF * 0.3) : 1
  return r.xp * share
}

export function buildLeaderboard(q: LeaderboardQuery, me: MeStats): Leaderboard {
  const rows = RIVALS.map((r) => ({
    id: r.id,
    nickname: r.nickname,
    level: r.level,
    value: +valueOf(r, q).toFixed(1),
  }))
  const mine = q.metric === 'speed' ? me.speed : q.metric === 'accuracy' ? me.accuracy : me.xp
  if (mine !== null) rows.push({ id: 'me', nickname: me.nickname, level: me.level, value: +mine.toFixed(1) })
  rows.sort((x, y) => y.value - x.value)

  const ranked = rows.map((r, i) => ({
    rank: i + 1,
    user: { id: r.id, nickname: r.nickname, level: r.level },
    value: r.value,
  }))
  const limit = q.limit ?? 20
  return {
    metric: q.metric,
    period: q.period,
    items: ranked.slice(0, limit),
    total: ranked.length,
    me: ranked.find((r) => r.user.id === 'me') ?? null,
  }
}
