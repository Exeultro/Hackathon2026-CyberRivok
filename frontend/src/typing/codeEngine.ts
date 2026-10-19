import type { KeyStat, PressResult } from './engine'
import { multiplierFor } from './engine'

export type CellState = 'ok' | 'bad' | 'pending' | 'skip'

export interface Cell {
  ch: string
  state: CellState
  idx: number
}

const OPENERS: Record<string, string> = { '(': ')', '[': ']', '{': '}', '"': '"', "'": "'", '`': '`' }
const CLOSERS = new Set(Object.values(OPENERS))
const BRACKET_OPEN: Record<string, string> = { ')': '(', ']': '[', '}': '{' }
const QUOTES = new Set(["'", '"', '`'])
const AHEAD = new Set(['(', ')', '[', ']', '{', '}', '<', '>', '"', "'", '`'])
const NAV = new Set(['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'Home', 'End'])
const WORD = /[\p{L}\p{N}_]/u
const INDENT = '    '

export class CodeEngine {
  readonly mode = 'code' as const
  buf: string[] = []
  cur = 0
  correct = 0
  errors = 0
  combo = 0
  maxCombo = 0
  streak = 0
  longestStreak = 0
  score = 0
  startedAt: number | null = null
  finishedAt: number | null = null
  private cleanKeys = 0
  private pairOf = new Map<number, number>()
  private nextId = 1
  private ids: number[] = []
  private bound = new Map<number, number>()
  private awards: number[] = []
  private typedFlag: boolean[] = []
  private stats = new Map<string, KeyStat>()
  private pairs = new Map<string, KeyStat>()

  constructor(
    readonly text: string,
    readonly autoPair = true,
  ) {
    this.pairOf = pairTable(text)
  }

  get done() {
    return this.buf.length === this.text.length && this.buf.join('') === this.text
  }

  get expected(): string | null {
    const { items, cursor } = this.cells()
    const c = items[cursor]
    return c && c.state !== 'bad' ? c.ch : null
  }

  get hasWrong() {
    return this.align().wrong > 0
  }

  progress() {
    const ok = this.cells().items.filter((c) => c.state === 'ok').length
    return this.text.length ? ok / this.text.length : 0
  }

  get skipped(): Cell[] {
    return this.cells().items.filter((c) => c.state === 'skip')
  }

  charsCorrect() {
    return Math.max(0, this.correct)
  }

  charsTyped() {
    return this.cleanKeys + this.errors
  }

  expectsIndent() {
    const c = this.cells().items[this.cells().cursor]
    if (!c || c.state === 'bad' || c.ch !== ' ') return false
    const atLineStart = c.idx === 0 || this.text[c.idx - 1] === '\n' || this.text[c.idx - 1] === ' '
    return atLineStart && this.text.slice(c.idx, c.idx + INDENT.length) === INDENT
  }

  cells(): { items: Cell[]; cursor: number } {
    const { map } = this.align()
    const items: Cell[] = []
    let next = 0
    let cursor = -1
    const gap = (upto: number) => {
      for (; next < upto; next++) items.push({ ch: this.text[next], state: 'pending', idx: next })
    }
    for (let i = 0; i <= this.buf.length; i++) {
      if (i === this.cur) cursor = items.length
      if (i === this.buf.length) break
      const m = map[i]
      if (m >= 0) {
        gap(m)
        const accepted = this.typedFlag[i] || i < this.cur
        items.push({ ch: this.text[m], state: accepted ? 'ok' : 'pending', idx: m })
        next = m + 1
      } else {
        items.push({ ch: this.buf[i], state: 'bad', idx: -1 })
      }
    }
    gap(this.text.length)
    if (cursor === -1) cursor = items.length
    const matched = new Set(map.filter((m) => m >= 0))
    for (let k = 0; k < cursor; k++) {
      if (items[k].state === 'pending' && !matched.has(items[k].idx)) items[k].state = 'skip'
    }
    return { items, cursor }
  }

  press(key: string, now: number): PressResult {
    if (this.finishedAt !== null) return 'ignored'
    if (NAV.has(key)) return this.move(key)

    let result: PressResult
    if (key === 'Backspace') result = this.erase(-1)
    else if (key === 'Delete') result = this.erase(0)
    else if (key === 'Tab') result = this.type(INDENT, now)
    else if (key === 'Enter' || key === '\n') result = this.type('\n', now)
    else if (key.length === 1) result = this.type(key, now)
    else return 'ignored'

    if (result !== 'ignored' && this.startedAt === null) this.startedAt = now
    if (this.done) this.finishedAt = now
    return result
  }

  durationMs(now: number) {
    if (this.startedAt === null) return 0
    return (this.finishedAt ?? now) - this.startedAt
  }

  accuracy() {
    const total = this.cleanKeys + this.errors
    return total ? (this.cleanKeys / total) * 100 : 100
  }

  cpm(now: number) {
    const minutes = this.durationMs(now) / 60000
    return minutes > 0 ? this.charsCorrect() / minutes : 0
  }

  wpm(now: number) {
    return this.cpm(now) / 5
  }

  keyStats(): KeyStat[] {
    return [...this.stats.values()]
  }

  bigramStats(): KeyStat[] {
    return [...this.pairs.values()]
  }

  private align() {
    const { map, wrong } = this.scan(this.buf.length)
    return { map, wrong, matched: map.filter((m) => m >= 0).length }
  }

  private scan(limit: number) {
    const buf = this.buf
    const n = this.text.length
    const map: number[] = new Array(buf.length).fill(-1)
    const open: { ch: string; at: number }[] = []
    const quote: Record<string, number> = {}
    const where = new Map<number, number>()
    this.ids.forEach((id, i) => where.set(id, i))
    let wrong = 0
    let j = 0

    for (let i = 0; i < limit; i++) {
      const c = buf[i]
      let partner = -1
      if (BRACKET_OPEN[c]) {
        const owner = where.get(this.bound.get(this.ids[i]) ?? -1)
        if (owner !== undefined && owner < i) {
          partner = owner
          const s = open.findIndex((o) => o.at === owner)
          if (s >= 0) open.splice(s, 1)
        } else {
          for (let s = open.length - 1; s >= 0; s--) {
            if (open[s].ch === BRACKET_OPEN[c]) {
              partner = open[s].at
              open.splice(s, 1)
              break
            }
          }
        }
      } else if (QUOTES.has(c)) {
        const owner = where.get(this.bound.get(this.ids[i]) ?? -1)
        if (owner !== undefined && owner < i) {
          partner = owner
          if (quote[c] === owner) delete quote[c]
        } else if (quote[c] !== undefined) {
          partner = quote[c]
          delete quote[c]
        } else {
          quote[c] = i
        }
      }

      let k = -1
      if (partner >= 0 && map[partner] >= 0) {
        const closer = this.pairOf.get(map[partner])
        if (closer !== undefined && closer >= j && this.text[closer] === c) k = closer
      }
      if (k < 0) {
        if (AHEAD.has(c)) {
          k = j
          while (k < n && this.text[k] !== c) k++
          if (k >= n) k = -1
        } else if (this.text[j] === c) {
          k = j
        }
      }

      if (k >= 0) {
        map[i] = k
        j = k + 1
        if (c === '(' || c === '[' || c === '{') open.push({ ch: c, at: i })
      } else {
        wrong += 1
      }
    }
    return { map, wrong, open, quote }
  }

  private bindCloser(at: number) {
    const c = this.buf[at]
    const { open, quote } = this.scan(at)
    let owner = -1
    if (BRACKET_OPEN[c]) {
      for (let s = open.length - 1; s >= 0; s--) {
        if (open[s].ch === BRACKET_OPEN[c]) {
          owner = open[s].at
          break
        }
      }
    } else if (QUOTES.has(c) && quote[c] !== undefined) {
      owner = quote[c]
    }
    if (owner >= 0) this.bound.set(this.ids[at], this.ids[owner])
  }

  private type(chars: string, now: number): PressResult {
    void now
    const before = this.align()
    const { items, cursor } = this.cells()
    const wanted = items[cursor] && items[cursor].state !== 'bad' ? items[cursor].ch : null
    const prev = this.cur > 0 ? this.buf[this.cur - 1] : ''
    const single = chars.length === 1 ? chars : null

    if (single && this.autoPair && CLOSERS.has(single) && this.buf[this.cur] === single && !this.typedFlag[this.cur]) {
      return this.award(single, prev, this.cur, true)
    }

    this.buf.splice(this.cur, 0, ...chars)
    this.awards.splice(this.cur, 0, ...Array.from({ length: chars.length }, () => 0))
    this.typedFlag.splice(this.cur, 0, ...Array.from({ length: chars.length }, () => false))
    this.ids.splice(this.cur, 0, ...Array.from({ length: chars.length }, () => this.nextId++))
    const at = this.cur
    this.cur += chars.length
    if (single && (BRACKET_OPEN[single] || QUOTES.has(single))) this.bindCloser(at)

    if (single && this.autoPair && OPENERS[single] && this.shouldPair(single, prev)) {
      this.buf.splice(this.cur, 0, OPENERS[single])
      this.awards.splice(this.cur, 0, 0)
      this.typedFlag.splice(this.cur, 0, false)
      this.ids.splice(this.cur, 0, this.nextId++)
      this.bound.set(this.ids[this.cur], this.ids[this.cur - 1])
    }

    if (this.align().wrong > before.wrong) {
      this.fail(wanted ?? chars[0], prev)
      return 'wrong'
    }
    const done = this.award(chars === INDENT ? ' ' : chars, prev, at, false)
    for (let i = 0; i < chars.length; i++) this.typedFlag[at + i] = true
    return done
  }

  private shouldPair(c: string, prev: string) {
    const next = this.buf[this.cur]
    if (c === '"' || c === "'" || c === '`') {
      return !(prev && WORD.test(prev)) && !(next && WORD.test(next))
    }
    return true
  }

  private award(key: string, prev: string, index: number, overtype: boolean): PressResult {
    this.combo += 1
    this.streak += 1
    this.maxCombo = Math.max(this.maxCombo, this.combo)
    this.longestStreak = Math.max(this.longestStreak, this.streak)
    const points = 10 * multiplierFor(this.combo)
    this.score += points
    this.awards[index] = points
    this.typedFlag[index] = true
    this.cleanKeys += 1
    this.correct += 1
    if (overtype) this.cur += 1
    this.tally(this.stats, key, false)
    if (prev && /\S/.test(prev) && /\S/.test(key)) this.tally(this.pairs, prev + key, false)
    return 'correct'
  }

  private fail(expected: string, prev: string) {
    this.errors += 1
    this.combo = 0
    this.streak = 0
    this.tally(this.stats, expected, true)
    if (prev && /\S/.test(prev) && /\S/.test(expected)) this.tally(this.pairs, prev + expected, true)
  }

  private erase(offset: -1 | 0): PressResult {
    let at = this.cur + offset
    if (at < 0 || at >= this.buf.length) return 'ignored'
    let count = 1
    if (offset === -1 && this.autoPair && OPENERS[this.buf[at]] === this.buf[at + 1] && !this.typedFlag[at + 1]) count = 2
    for (let i = 0; i < count; i++) {
      const award = this.awards[at]
      if (award > 0) {
        this.correct -= 1
        this.score -= award
      }
      this.buf.splice(at, 1)
      this.awards.splice(at, 1)
      this.typedFlag.splice(at, 1)
      this.ids.splice(at, 1)
    }
    if (offset === -1) this.cur = at
    this.combo = 0
    this.streak = 0
    return 'back'
  }

  private move(key: string): PressResult {
    const lineStart = (i: number) => this.buf.lastIndexOf('\n', i - 1) + 1
    const lineEnd = (i: number) => {
      const k = this.buf.indexOf('\n', i)
      return k === -1 ? this.buf.length : k
    }
    if (key === 'ArrowLeft') this.cur = Math.max(0, this.cur - 1)
    else if (key === 'ArrowRight') this.cur = Math.min(this.buf.length, this.cur + 1)
    else if (key === 'Home') this.cur = lineStart(this.cur)
    else if (key === 'End') this.cur = lineEnd(this.cur)
    else {
      const start = lineStart(this.cur)
      const col = this.cur - start
      if (key === 'ArrowUp') {
        if (start === 0) this.cur = 0
        else {
          const prevStart = lineStart(start - 1)
          this.cur = prevStart + Math.min(col, start - 1 - prevStart)
        }
      } else {
        const end = lineEnd(this.cur)
        if (end >= this.buf.length) this.cur = this.buf.length
        else this.cur = end + 1 + Math.min(col, lineEnd(end + 1) - (end + 1))
      }
    }
    return 'move'
  }

  private tally(map: Map<string, KeyStat>, key: string, error: boolean) {
    let s = map.get(key)
    if (!s) {
      s = { key, hits: 0, errors: 0 }
      map.set(key, s)
    }
    s.hits += 1
    if (error) s.errors += 1
  }
}

function pairTable(text: string) {
  const pairs = new Map<number, number>()
  const stack: number[] = []
  for (let i = 0; i < text.length; i++) {
    const c = text[i]
    if (c === '(' || c === '[' || c === '{') stack.push(i)
    else if (BRACKET_OPEN[c]) {
      for (let s = stack.length - 1; s >= 0; s--) {
        if (text[stack[s]] === BRACKET_OPEN[c]) {
          pairs.set(stack[s], i)
          stack.length = s
          break
        }
      }
    }
  }
  for (const q of QUOTES) {
    const at: number[] = []
    for (let i = 0; i < text.length; i++) if (text[i] === q) at.push(i)
    for (let a = 0; a + 1 < at.length; a += 2) pairs.set(at[a], at[a + 1])
  }
  return pairs
}
