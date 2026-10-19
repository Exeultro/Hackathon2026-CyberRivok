export const BACKSPACE = '\b'

export type PressResult = 'correct' | 'wrong' | 'back' | 'move' | 'ignored'
export type ErrorMode = 'strict' | 'real'

export interface KeyStat {
  key: string
  hits: number
  errors: number
}

export function multiplierFor(combo: number) {
  if (combo >= 50) return 5
  if (combo >= 25) return 3
  if (combo >= 10) return 2
  return 1
}

export class Engine {
  pos = 0
  correct = 0
  errors = 0
  combo = 0
  maxCombo = 0
  streak = 0
  longestStreak = 0
  score = 0
  typed: string[] = []
  startedAt: number | null = null
  finishedAt: number | null = null
  private awarded: number[] = []
  private stats = new Map<string, KeyStat>()
  private pairs = new Map<string, KeyStat>()

  constructor(
    readonly text: string,
    readonly mode: ErrorMode = 'strict',
  ) {}

  get done() {
    return this.pos >= this.text.length && this.uncorrected === 0
  }

  get expected(): string | null {
    return this.pos >= this.text.length ? null : this.text[this.pos]
  }

  progress() {
    return this.text.length ? this.pos / this.text.length : 0
  }

  charsCorrect() {
    return this.correct
  }

  charsTyped() {
    return this.correct + this.errors
  }

  get uncorrected() {
    return this.typed.reduce((n, ch, i) => n + (ch === this.text[i] ? 0 : 1), 0)
  }

  press(ch: string, now: number): PressResult {
    if (this.finishedAt !== null) return 'ignored'
    if (ch === BACKSPACE) return this.back()
    if (this.pos >= this.text.length) return 'ignored'
    if (this.startedAt === null) this.startedAt = now

    const want = this.text[this.pos]
    const stat = this.stat(want)
    const pair = this.pairStat()
    stat.hits += 1
    if (pair) pair.hits += 1

    if (ch === want) {
      this.correct += 1
      this.combo += 1
      this.streak += 1
      this.maxCombo = Math.max(this.maxCombo, this.combo)
      this.longestStreak = Math.max(this.longestStreak, this.streak)
      const points = 10 * multiplierFor(this.combo)
      this.score += points
      this.typed[this.pos] = ch
      this.awarded[this.pos] = points
      this.advance(now)
      return 'correct'
    }

    stat.errors += 1
    if (pair) pair.errors += 1
    this.errors += 1
    this.combo = 0
    this.streak = 0
    if (this.mode === 'real') {
      this.typed[this.pos] = ch
      this.awarded[this.pos] = 0
      this.advance(now)
    }
    return 'wrong'
  }

  durationMs(now: number) {
    if (this.startedAt === null) return 0
    return (this.finishedAt ?? now) - this.startedAt
  }

  accuracy() {
    const total = this.correct + this.errors
    return total ? (this.correct / total) * 100 : 100
  }

  cpm(now: number) {
    const minutes = this.durationMs(now) / 60000
    return minutes > 0 ? this.correct / minutes : 0
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

  private pairStat(): KeyStat | null {
    if (this.pos === 0) return null
    const key = this.text[this.pos - 1] + this.text[this.pos]
    if (/\s/.test(key)) return null
    let s = this.pairs.get(key)
    if (!s) {
      s = { key, hits: 0, errors: 0 }
      this.pairs.set(key, s)
    }
    return s
  }

  private advance(now: number) {
    this.pos += 1
    if (this.done) this.finishedAt = now
  }

  private back(): PressResult {
    if (this.mode !== 'real' || this.pos === 0) return 'ignored'
    const i = this.pos - 1
    if (this.typed[i] === this.text[i]) {
      this.correct -= 1
      this.score -= this.awarded[i]
      this.combo = 0
      this.streak = 0
    }
    this.typed.pop()
    this.awarded.pop()
    this.pos = i
    return 'back'
  }

  private stat(key: string): KeyStat {
    let s = this.stats.get(key)
    if (!s) {
      s = { key, hits: 0, errors: 0 }
      this.stats.set(key, s)
    }
    return s
  }
}
