import type { Achievement } from './types'
import type { MockState } from './mockTypes'
import { streakOf } from './mockUtil'

interface Def {
  code: Achievement['code']
  title: string
  description: string
  icon: Achievement['icon']
  metric: string
  target: number
}

const CATALOG: Def[] = [
  { code: 'first_steps', title: 'Первый шаг', description: 'Заверши первую тренировку', icon: 'keyboard', metric: 'sessions', target: 1 },
  { code: 'speed_200', title: 'Скорость 200', description: 'Достигни 200 символов в минуту', icon: 'bolt', metric: 'best_cpm', target: 200 },
  { code: 'speed_300', title: 'Скорострел', description: 'Достигни 300 символов в минуту', icon: 'bolt', metric: 'best_cpm', target: 300 },
  { code: 'genius', title: 'Печатный гений', description: 'Достигни 100 слов в минуту', icon: 'trophy', metric: 'best_wpm', target: 100 },
  { code: 'clean_100', title: 'Без единой ошибки', description: 'Напечатай 100 символов подряд без ошибок', icon: 'shield', metric: 'clean_streak', target: 100 },
  { code: 'steel_finger', title: 'Стальной палец', description: 'Напечатай 250 символов подряд без ошибок', icon: 'shield', metric: 'clean_streak', target: 250 },
  { code: 'perfect', title: 'Безупречно', description: 'Заверши тренировку со 100% точностью', icon: 'target', metric: 'perfect', target: 1 },
  { code: 'thousand', title: 'Тысяча символов', description: 'Напечатай 1000 символов', icon: 'keyboard', metric: 'total_chars', target: 1000 },
  { code: 'marathon', title: 'Марафонец', description: 'Напечатай 10 000 символов', icon: 'clock', metric: 'total_chars', target: 10000 },
  { code: 'streak_7', title: 'Неделя подряд', description: 'Тренируйся 7 дней подряд', icon: 'flame', metric: 'streak', target: 7 },
  { code: 'stars_10', title: 'Звёздный путь', description: 'Собери 10 звёзд на картах', icon: 'star', metric: 'stars', target: 10 },
  { code: 'stars_30', title: 'Созвездие', description: 'Собери 30 звёзд на картах', icon: 'star', metric: 'stars', target: 30 },
  { code: 'combo_50', title: 'Комбо-мастер', description: 'Набери комбо ×5 (50 подряд)', icon: 'bolt', metric: 'max_combo', target: 50 },
  { code: 'coder', title: 'Кодер', description: 'Заверши тренировку в режиме «Код»', icon: 'code', metric: 'code_sessions', target: 1 },
  { code: 'polyglot', title: 'Полиглот', description: 'Потренируйся на русском и английском', icon: 'globe', metric: 'langs', target: 2 },
]

function metricsOf(state: MockState): Record<string, number> {
  const s = state.sessions
  const max = (f: (x: (typeof s)[number]) => number) => s.reduce((m, x) => Math.max(m, f(x)), 0)
  return {
    sessions: s.length,
    best_cpm: max((x) => x.cpm),
    best_wpm: max((x) => x.wpm),
    total_chars: s.reduce((n, x) => n + x.chars, 0),
    streak: streakOf(s),
    stars: Object.values(state.stars).reduce((n, v) => n + v, 0),
    max_combo: max((x) => x.max_combo ?? 0),
    clean_streak: max((x) => x.clean_streak ?? 0),
    perfect: s.filter((x) => x.accuracy >= 100 && x.chars >= 20).length,
    code_sessions: s.filter((x) => x.mode === 'code').length,
    langs: new Set(s.map((x) => x.lang)).size,
  }
}

function toAchievement(d: Def, value: number, unlockedAt?: string): Achievement {
  return {
    code: d.code,
    title: d.title,
    description: d.description,
    icon: d.icon,
    target: d.target,
    progress: Math.round(value * 10) / 10,
    unlocked_at: unlockedAt ?? null,
  }
}

export function listAchievements(state: MockState): Achievement[] {
  const m = metricsOf(state)
  return CATALOG.map((d) => toAchievement(d, m[d.metric] ?? 0, state.achievements?.[d.code]))
}

export function unlockNew(state: MockState): Achievement[] {
  state.achievements ??= {}
  const m = metricsOf(state)
  const fresh: Achievement[] = []
  for (const d of CATALOG) {
    if (state.achievements[d.code] || (m[d.metric] ?? 0) < d.target) continue
    const at = new Date().toISOString()
    state.achievements[d.code] = at
    fresh.push(toAchievement(d, m[d.metric], at))
  }
  return fresh
}
