<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { CodeLanguage, Mode } from '@/api'
import { useGameStore } from '@/stores/game'
import TypingSession from '@/components/TypingSession.vue'
import type { SessionConfig } from '@/components/TypingSession.vue'

type Length = 'short' | 'medium' | 'long'

const MODES: { id: Mode; label: string; hint: string }[] = [
  { id: 'words', label: 'Слова', hint: 'Случайные слова из словаря' },
  { id: 'sentences', label: 'Тексты', hint: 'Настоящая литература: от справочников до стихов' },
  { id: 'code', label: 'Код', hint: 'Фрагменты кода' },
  { id: 'symbols', label: 'Символы', hint: 'Буквы, цифры и спецсимволы' },
]
const CODE: { id: CodeLanguage; label: string }[] = [
  { id: 'python', label: 'Python' },
  { id: 'javascript', label: 'JavaScript' },
  { id: 'html', label: 'HTML' },
  { id: 'css', label: 'CSS' },
]
const LENGTHS: { id: Length; label: string }[] = [
  { id: 'short', label: 'Короткая' },
  { id: 'medium', label: 'Средняя' },
  { id: 'long', label: 'Длинная' },
]
const COUNTS: Record<Mode, Record<Length, number>> = {
  words: { short: 1, medium: 3, long: 5 },
  sentences: { short: 1, medium: 2, long: 3 },
  symbols: { short: 1, medium: 3, long: 5 },
  code: { short: 1, medium: 1, long: 1 },
}

interface Saved {
  mode: Mode
  difficulty: number
  length: Length
  code: CodeLanguage
}
const KEY = 'hero.practice'

function readSaved(): Saved {
  const fallback: Saved = { mode: 'words', difficulty: 2, length: 'medium', code: 'python' }
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? { ...fallback, ...(JSON.parse(raw) as Partial<Saved>) } : fallback
  } catch {
    return fallback
  }
}

const game = useGameStore()
const saved = readSaved()
const mode = ref<Mode>(saved.mode)
const difficulty = ref(saved.difficulty)
const length = ref<Length>(saved.length)
const code = ref<CodeLanguage>(saved.code)
const running = ref(false)
const runId = ref(0)

watch([mode, difficulty, length, code], () => {
  try {
    localStorage.setItem(
      KEY,
      JSON.stringify({ mode: mode.value, difficulty: difficulty.value, length: length.value, code: code.value }),
    )
  } catch {
    // settings just will not persist
  }
})

const TEXT_LEVELS = ['Справочная и техническая литература', 'Научная проза', 'Классическая проза', 'Стихи с ясным ритмом', 'Сложные стихи']

const hint = computed(() => MODES.find((m) => m.id === mode.value)?.hint ?? '')

const config = computed<SessionConfig>(() => ({
  mode: mode.value,
  difficulty: difficulty.value,
  count: COUNTS[mode.value][length.value],
  code_language: mode.value === 'code' ? code.value : undefined,
  title: `Свободная тренировка - ${MODES.find((m) => m.id === mode.value)?.label}`,
}))

function start() {
  runId.value += 1
  running.value = true
}
</script>

<template>
  <main>
    <TypingSession
      v-if="running"
      :key="runId"
      :config="config"
      :node="null"
      back-label="← Настройки"
      @exit="running = false"
    />

    <section v-else class="setup">
      <h1>Свободная тренировка</h1>
      <p class="muted">Без карты и звёзд: выбирай что печатать, XP и статистика всё равно копятся.</p>

      <div class="group">
        <h2>Режим</h2>
        <div class="seg" role="group" aria-label="Режим">
          <button v-for="m in MODES" :key="m.id" :aria-pressed="mode === m.id" @click="mode = m.id">
            {{ m.label }}
          </button>
        </div>
        <p class="muted small">{{ hint }}</p>
      </div>

      <div v-if="mode === 'code'" class="group">
        <h2>Язык программирования</h2>
        <div class="seg" role="group" aria-label="Язык программирования">
          <button v-for="c in CODE" :key="c.id" :aria-pressed="code === c.id" @click="code = c.id">
            {{ c.label }}
          </button>
        </div>
      </div>

      <div v-if="mode !== 'code'" class="group">
        <h2>Сложность</h2>
        <div class="seg" role="group" aria-label="Сложность">
          <button v-for="n in 5" :key="n" :aria-pressed="difficulty === n" @click="difficulty = n">{{ n }}</button>
        </div>
        <p class="muted small">
          {{ mode === 'sentences' ? TEXT_LEVELS[difficulty - 1] : 'Влияет на длину слов и набор символов.' }}
        </p>
      </div>

      <div v-if="mode !== 'code'" class="group">
        <h2>Длина</h2>
        <div class="seg" role="group" aria-label="Длина">
          <button v-for="l in LENGTHS" :key="l.id" :aria-pressed="length === l.id" @click="length = l.id">
            {{ l.label }}
          </button>
        </div>
      </div>

      <p class="muted small">
        Язык набора: <b>{{ game.lang === 'ru' ? 'русский' : 'английский' }}</b>, меняется переключателем RU/EN вверху.
      </p>

      <button class="btn btn-primary start" @click="start">Начать</button>
    </section>
  </main>
</template>

<style scoped>
.setup {
  max-width: 620px;
  margin: 0 auto;
  padding: 16px;
}
h1 {
  margin: 0 0 4px;
}
h2 {
  margin: 0 0 8px;
  font-size: 14px;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--muted);
}
.muted {
  color: var(--muted);
}
.small {
  font-size: 13px;
  margin: 6px 0 0;
}
.group {
  margin: 18px 0;
}
.seg {
  display: inline-flex;
  flex-wrap: wrap;
  background: var(--surface-2);
  border-radius: 12px;
  padding: 3px;
  gap: 3px;
}
.seg button {
  border: 0;
  background: transparent;
  color: var(--muted);
  padding: 8px 16px;
  border-radius: 9px;
  cursor: pointer;
  font-weight: 700;
}
.seg button[aria-pressed='true'] {
  background: var(--accent);
  color: var(--accent-ink);
}
.start {
  margin-top: 12px;
  width: 100%;
  font-size: 18px;
  padding: 14px;
}
</style>
