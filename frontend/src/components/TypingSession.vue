<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { api } from '@/api'
import type { CodeLanguage, Lang, MapNode, Mode, SessionResult } from '@/api'
import { useGameStore } from '@/stores/game'
import { BACKSPACE, Engine, multiplierFor } from '@/typing/engine'
import { CodeEngine } from '@/typing/codeEngine'
import { sfx } from '@/useSound'
import TypingText from '@/components/TypingText.vue'
import KeyboardHint from '@/components/KeyboardHint.vue'
import ResultCard from '@/components/ResultCard.vue'

export interface SessionConfig {
  mode: Mode
  difficulty: number
  count: number
  code_language?: CodeLanguage
  nodeId?: string
  title: string
}

const props = defineProps<{ config: SessionConfig; node: MapNode | null; backLabel: string }>()
const emit = defineEmits<{ exit: []; next: [id: string] }>()
const game = useGameStore()

const EDIT_KEYS = ['Tab', 'Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'Home', 'End']

type Status = 'loading' | 'ready' | 'running' | 'finished' | 'error'
const status = ref<Status>('loading')
const engine = shallowRef<Engine | CodeEngine | null>(null)
const tick = ref(0)
const nowMs = ref(performance.now())
const wrong = ref(false)
const shake = ref(false)
const pop = ref(0)
const sparks = ref<{ id: number; x: number; y: number }[]>([])
const result = ref<SessionResult | null>(null)
const sources = ref<string[]>([])
const textRef = ref<InstanceType<typeof TypingText> | null>(null)
const rootRef = ref<HTMLElement | null>(null)

let timer: number | undefined
let sparkId = 0

const playLang = computed<Lang>(() => (props.config.mode === 'code' ? 'en' : game.lang))

const state = computed(() => {
  void tick.value
  const e = engine.value
  if (!e) return null
  const common = {
    progress: e.progress(),
    combo: e.combo,
    mult: multiplierFor(e.combo),
    score: e.score,
    wpm: Math.round(e.wpm(nowMs.value)),
    acc: Math.round(e.accuracy()),
    expected: e.expected,
  }
  if (e instanceof CodeEngine) {
    const { items, cursor } = e.cells()
    const skipped = items.filter((c) => c.state === 'skip').length
    const message = e.hasWrong
      ? 'Ошибка. Нажми Backspace, чтобы исправить'
      : skipped
        ? 'Пропущен символ (оранжевый). Нажми ←, впиши его'
        : ''
    return { ...common, pos: 0, typed: [] as string[], cells: items, cursor, message, indent: e.expectsIndent() }
  }
  return {
    ...common,
    pos: e.pos,
    typed: [...e.typed],
    cells: undefined,
    cursor: undefined,
    message:
      game.realMode && e.pos > 0 && (e.typed[e.pos - 1] !== e.text[e.pos - 1] || (e.pos >= e.text.length && e.uncorrected > 0))
        ? 'Ошибка. Нажми Backspace, чтобы исправить'
        : '',
    indent: false,
  }
})

async function start() {
  const c = props.config
  stopTimer()
  result.value = null
  wrong.value = false
  status.value = 'loading'
  try {
    const items = await api.getTexts({
      mode: c.mode,
      lang: playLang.value,
      difficulty: c.difficulty,
      count: c.mode === 'code' ? 1 : c.count,
      code_language: c.code_language,
    })
    sources.value = items.map((i) => i.source).filter((s): s is string => !!s)
    const text = items.map((i) => i.text).join(c.mode === 'code' || c.mode === 'sentences' ? '\n' : ' ')
    engine.value =
      c.mode === 'code' ? new CodeEngine(text, game.autoPair) : new Engine(text, game.realMode ? 'real' : 'strict')
    tick.value += 1
    status.value = 'ready'
  } catch {
    status.value = 'error'
  }
}

function fit() {
  const el = rootRef.value
  if (!el) return
  let scale = 1
  el.style.setProperty('--fit', '1')
  while (document.documentElement.scrollHeight > window.innerHeight && scale > 0.62) {
    scale -= 0.04
    el.style.setProperty('--fit', scale.toFixed(2))
  }
}

function stopTimer() {
  if (timer !== undefined) window.clearInterval(timer)
  timer = undefined
}

function burst() {
  nextTick(() => {
    const p = textRef.value?.cursorPoint()
    if (!p) return
    const id = ++sparkId
    sparks.value.push({ id, ...p })
    window.setTimeout(() => {
      sparks.value = sparks.value.filter((s) => s.id !== id)
    }, 750)
  })
}

function flashWrong() {
  shake.value = true
  window.setTimeout(() => (shake.value = false), 220)
  if (game.sound) sfx.error()
}

async function finish(e: Engine | CodeEngine) {
  stopTimer()
  status.value = 'finished'
  const c = props.config
  const duration = e.durationMs(performance.now())
  try {
    const res = await api.createSession({
      mode: c.mode,
      lang: playLang.value,
      code_language: c.code_language,
      difficulty: c.difficulty,
      node_id: c.nodeId,
      started_at: new Date(Date.now() - duration).toISOString(),
      duration_ms: Math.max(1000, Math.round(duration)),
      chars_typed: e.charsTyped(),
      chars_correct: e.charsCorrect(),
      errors: e.errors,
      max_combo: e.maxCombo,
      longest_clean_streak: e.longestStreak,
      client_score: e.score,
      key_stats: e.keyStats(),
      bigram_stats: e.bigramStats(),
    })
    result.value = res
    celebrate(res)
    await game.load()
  } catch {
    status.value = 'error'
  }
}

function celebrate(res: SessionResult) {
  if (!game.sound) return
  const stars = res.node_result?.stars ?? 0
  if (res.node_result && !res.node_result.passed) {
    sfx.fail()
    return
  }
  sfx.win()
  for (let i = 0; i < stars; i++) window.setTimeout(() => sfx.star(i), 450 + i * 220)
}

function onKey(ev: KeyboardEvent) {
  if (ev.key === 'Escape') {
    emit('exit')
    return
  }
  const code = props.config.mode === 'code'
  if (ev.key === 'Tab' && !code) {
    ev.preventDefault()
    start()
    return
  }
  const e = engine.value
  if (!e || (status.value !== 'ready' && status.value !== 'running')) return
  if (ev.ctrlKey || ev.metaKey || ev.altKey) return
  let ch: string | null = null
  if (ev.key === 'Enter') ch = '\n'
  else if (code && EDIT_KEYS.includes(ev.key)) ch = ev.key
  else if (ev.key === 'Backspace') ch = BACKSPACE
  else if (ev.key.length === 1) ch = ev.key
  if (ch === null) return
  ev.preventDefault()

  const now = performance.now()
  const before = multiplierFor(e.combo)
  const res = e.press(ch, now)
  nowMs.value = now
  tick.value += 1

  if (status.value === 'ready' && res !== 'move' && res !== 'ignored') {
    status.value = 'running'
    timer = window.setInterval(() => (nowMs.value = performance.now()), 250)
  }

  if (res === 'wrong') {
    wrong.value = !game.realMode && !code
    flashWrong()
  } else if (res === 'back') {
    wrong.value = false
    if (game.sound) sfx.back()
  } else if (res === 'correct') {
    wrong.value = false
    if (game.sound) sfx.key()
    const after = multiplierFor(e.combo)
    if (after > before || (e.combo >= 50 && e.combo % 25 === 0)) {
      pop.value += 1
      burst()
      if (game.sound) sfx.combo(after)
    }
    if (e.done) finish(e)
  }
}

function goNext() {
  const next = result.value?.node_result?.next_node_id
  if (next) emit('next', next)
}

onMounted(() => {
  window.addEventListener('keydown', onKey)
  window.addEventListener('resize', fit)
  start()
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
  window.removeEventListener('resize', fit)
  stopTimer()
})

watch(
  () => [engine.value?.text, game.showKeyboard, sources.value.length],
  () => nextTick(fit),
  { flush: 'post' },
)

watch(() => JSON.stringify(props.config), start)
watch(
  () => [game.lang, game.realMode, game.autoPair],
  () => {
    if (status.value !== 'finished') start()
  },
)
</script>

<template>
  <div ref="rootRef" class="play">
    <div class="bar">
      <button class="back" @click="emit('exit')">{{ backLabel }}</button>
      <div class="title">{{ config.title }}</div>
      <div class="spacer" />
      <button class="tgl" @click="start">Заново</button>
      <button
        v-if="config.mode === 'code'"
        class="tgl"
        :aria-pressed="game.autoPair"
        title="Скобки и кавычки закрываются сами, как в редакторе"
        @click="game.autoPair = !game.autoPair"
      >
        Автоскобки
      </button>
      <button
        v-else
        class="tgl"
        :aria-pressed="game.realMode"
        title="Ошибочная буква впечатывается красным, исправляй её Backspace"
        @click="game.realMode = !game.realMode"
      >
        Реальный режим
      </button>
      <button class="tgl" :aria-pressed="game.sound" @click="game.sound = !game.sound">Звук</button>
      <button class="tgl" :aria-pressed="game.showKeyboard" @click="game.showKeyboard = !game.showKeyboard">
        Клавиатура
      </button>
    </div>

    <template v-if="state && engine">
      <div class="hud">
        <div class="stat"><b>{{ state.wpm }}</b><small>WPM</small></div>
        <div class="stat"><b>{{ state.acc }}%</b><small>точность</small></div>
        <div class="stat"><b>{{ state.score }}</b><small>очки</small></div>
        <div :key="pop" class="combo" :class="'m' + state.mult">
          <b>×{{ state.mult }}</b>
          <small>комбо {{ state.combo }}</small>
        </div>
      </div>
      <div class="progress"><i :style="{ width: state.progress * 100 + '%' }" /></div>

      <TypingText
        ref="textRef"
        :text="engine.text"
        :typed="state.typed"
        :cells="state.cells"
        :cursor="state.cursor"
        :pos="state.pos"
        :wrong="wrong"
        :unit="game.unit"
        :shake="shake"
        :sparks="sparks"
      />
      <p class="start" :class="{ hidden: status !== 'ready' }">Нажми любую клавишу, чтобы начать</p>
      <p v-if="sources.length" class="source">{{ sources.join(' · ') }}</p>
      <p v-if="config.mode === 'code'" class="note">Tab = 4 пробела · стрелки двигают курсор · наперёд только скобки</p>
      <KeyboardHint v-if="game.showKeyboard" :lang="playLang" :expected="state.expected" :message="state.message" :tab="state.indent" />
    </template>
    <p v-else-if="status === 'error'" class="muted">Не удалось загрузить текст.</p>
    <p v-else class="muted">Загрузка...</p>

    <ResultCard
      v-if="result"
      :result="result"
      :node="node"
      :back-label="backLabel"
      @retry="start"
      @map="emit('exit')"
      @next="goNext"
    />
  </div>
</template>

<style scoped>
.play {
  max-width: 920px;
  margin: 0 auto;
  padding: 8px 16px 12px;
}
.bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.back {
  color: var(--muted);
  background: none;
  border: 0;
  cursor: pointer;
  font-weight: 700;
  font-size: 16px;
  padding: 0;
}
.title {
  font-weight: 800;
}
.spacer {
  flex: 1;
}
.tgl {
  border: 1px solid var(--border);
  background: transparent;
  color: var(--muted);
  border-radius: 10px;
  padding: 4px 10px;
  font-size: 13px;
  cursor: pointer;
}
.tgl[aria-pressed='true'] {
  color: var(--accent);
  border-color: var(--accent);
}
.hud {
  display: flex;
  gap: 14px;
  align-items: center;
  margin-bottom: 4px;
}
.stat b,
.combo b {
  display: block;
  font-size: 24px;
  line-height: 1.1;
}
.stat small,
.combo small {
  color: var(--muted);
}
.combo {
  margin-left: auto;
  text-align: right;
  animation: pop 0.3s;
}
.combo.m2 b {
  color: var(--accent);
}
.combo.m3 b {
  color: var(--map-chase);
}
.combo.m5 b {
  color: var(--err);
}
@keyframes pop {
  40% {
    transform: scale(1.3) rotate(-4deg);
  }
}
.progress {
  height: 6px;
  border-radius: 3px;
  background: var(--surface-2);
  overflow: hidden;
  margin-bottom: 10px;
}
.progress i {
  display: block;
  height: 100%;
  background: var(--accent);
  transition: width 0.1s;
}
.start {
  text-align: center;
  color: var(--muted);
  margin: 6px 0 0;
}
.start.hidden {
  visibility: hidden;
}
.source {
  text-align: right;
  color: var(--muted);
  font-size: 14px;
  font-style: italic;
  margin: 8px 4px 0;
}
.note {
  text-align: center;
  color: var(--accent);
  margin: 2px 0 0;
  font-size: 13px;
}
.muted {
  color: var(--muted);
}
@media (prefers-reduced-motion: reduce) {
  .combo {
    animation: none;
  }
}
</style>
