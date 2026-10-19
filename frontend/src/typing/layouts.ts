import type { Lang } from '@/api'

export interface KeyDef {
  n: string
  s: string
  f: number
}

export const FINGERS = [
  'левый мизинец',
  'левый безымянный',
  'левый средний',
  'левый указательный',
  'правый указательный',
  'правый средний',
  'правый безымянный',
  'правый мизинец',
  'большой палец',
]

export const FINGER_COLORS = [
  '#a78bfa',
  '#60a5fa',
  '#34d399',
  '#fbbf24',
  '#fbbf24',
  '#34d399',
  '#60a5fa',
  '#a78bfa',
  '#f472b6',
]

function row(normal: string, shifted: string, fingers: number[]): KeyDef[] {
  return [...normal].map((n, i) => ({ n, s: shifted[i], f: fingers[i] }))
}

const R1 = [0, 0, 1, 2, 3, 3, 4, 4, 5, 6, 7, 7, 7]
const R2 = [0, 1, 2, 3, 3, 4, 4, 5, 6, 7, 7, 7, 7]
const R3 = [0, 1, 2, 3, 3, 4, 4, 5, 6, 7, 7]
const R4 = [0, 1, 2, 3, 3, 4, 4, 5, 6, 7]

export const LAYOUTS: Record<Lang, KeyDef[][]> = {
  en: [
    row('`1234567890-=', '~!@#$%^&*()_+', R1),
    row('qwertyuiop[]\\', 'QWERTYUIOP{}|', R2),
    row("asdfghjkl;'", 'ASDFGHJKL:"', R3),
    row('zxcvbnm,./', 'ZXCVBNM<>?', R4),
  ],
  ru: [
    row('ё1234567890-=', 'Ё!"№;%:?*()_+', R1),
    row('йцукенгшщзхъ\\', 'ЙЦУКЕНГШЩЗХЪ/', R2),
    row('фывапролджэ', 'ФЫВАПРОЛДЖЭ', R3),
    row('ячсмитьбю.', 'ЯЧСМИТЬБЮ,', R4),
  ],
}

export interface Target {
  row: number
  col: number
  finger: number
  shift: boolean
}

export function findTarget(lang: Lang, ch: string | null): Target | null {
  if (ch === null) return null
  if (ch === ' ') return { row: 4, col: 0, finger: 8, shift: false }
  if (ch === '\n') return { row: 2, col: LAYOUTS[lang][2].length, finger: 7, shift: false }
  const rows = LAYOUTS[lang]
  for (let r = 0; r < rows.length; r++) {
    for (let c = 0; c < rows[r].length; c++) {
      const k = rows[r][c]
      if (k.n === ch) return { row: r, col: c, finger: k.f, shift: false }
      if (k.s === ch) return { row: r, col: c, finger: k.f, shift: true }
    }
  }
  return null
}
