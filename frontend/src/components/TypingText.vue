<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import type { Cell } from '@/typing/codeEngine'

const props = defineProps<{
  text: string
  typed: string[]
  pos: number
  wrong: boolean
  unit?: 'char' | 'word'
  shake: boolean
  sparks: { id: number; x: number; y: number }[]
  cells?: Cell[]
  cursor?: number
}>()

interface View {
  text: string
  cls: string
  br: boolean
}

const root = ref<HTMLElement | null>(null)
const inner = ref<HTMLElement | null>(null)
const offset = ref(0)
const VISIBLE_LINES = 5
const chars = computed(() => [...props.text])
const RAYS = Array.from({ length: 8 }, (_, k) => {
  const a = (k * Math.PI) / 4
  return { dx: Math.cos(a) * 34, dy: Math.sin(a) * 34 }
})

const wordRange = computed(() => {
  if (props.unit !== 'word' || props.cells) return null
  const text = chars.value
  let a = Math.min(props.pos, text.length - 1)
  let b = a
  while (a > 0 && !/\s/.test(text[a - 1])) a--
  while (b < text.length && !/\s/.test(text[b])) b++
  return { a, b }
})

function inWord(i: number) {
  const r = wordRange.value
  return !!r && i >= r.a && i < r.b
}

function fromCells(cells: Cell[], cursor: number): View[] {
  return cells.map((c, i) => {
    const bad = c.state === 'bad'
    const state = c.state === 'ok' ? 'ok' : bad ? 'bad' : c.state === 'skip' ? 'skip' : ''
    const here = i === cursor ? (bad ? ' cur wrong' : ' cur') : ''
    const glyph = c.ch === '\n' ? '↵' : c.ch === ' ' && (bad || c.state === 'skip') ? '·' : c.ch
    return { text: glyph, cls: (state + here).trim(), br: c.ch === '\n' && !bad }
  })
}

function fromText(): View[] {
  return chars.value.map((ch, i) => {
    let cls = ''
    if (i < props.pos) cls = (props.typed[i] === ch ? 'ok' : 'bad') + (inWord(i) ? ' word' : '')
    else if (i === props.pos) cls = (props.wrong ? 'cur wrong' : 'cur') + (inWord(i) ? ' word' : '')
    else if (inWord(i)) cls = 'word'
    const shown = i < props.pos ? props.typed[i] : ch
    const glyph =
      shown === '\n' ? '↵' : shown === ' ' && i < props.pos && props.typed[i] !== ch ? '·' : shown
    return { text: glyph, cls, br: ch === '\n' }
  })
}

const view = computed<View[]>(() =>
  props.cells ? fromCells(props.cells, props.cursor ?? 0) : fromText(),
)
const caretAtEnd = computed(() => !!props.cells && (props.cursor ?? 0) >= props.cells.length)

function follow() {
  const box = inner.value
  const el = box?.querySelector<HTMLElement>('.cur, .caret')
  if (!box || !el) return
  const lh = parseFloat(getComputedStyle(el).lineHeight)
  if (!lh) return
  const line = Math.round(el.offsetTop / lh)
  const maxOffset = Math.max(0, box.scrollHeight - VISIBLE_LINES * lh)
  offset.value = Math.min(Math.max(0, line - 1) * lh, maxOffset)
}

watch(
  () => [props.pos, props.cursor, props.text, props.cells],
  () => nextTick(follow),
  { flush: 'post' },
)
onMounted(follow)

function cursorPoint() {
  const el = root.value?.querySelector<HTMLElement>('.cur, .caret')
  if (!el || !root.value) return { x: 0, y: 0 }
  const a = el.getBoundingClientRect()
  const b = root.value.getBoundingClientRect()
  return { x: a.left - b.left + a.width / 2, y: a.top - b.top + a.height / 2 }
}

defineExpose({ cursorPoint })
</script>

<template>
  <div ref="root" class="text" :class="{ shake }">
    <div class="viewport">
      <div ref="inner" class="inner" :style="{ transform: `translateY(${-offset}px)` }">
        <template v-for="(v, i) in view" :key="i">
          <span class="ch" :class="v.cls">{{ v.text }}</span><br v-if="v.br" />
        </template>
        <span v-if="caretAtEnd" class="caret">&nbsp;</span>
      </div>
    </div>
    <span v-for="s in sparks" :key="s.id" class="burst" :style="{ left: s.x + 'px', top: s.y + 'px' }">
      <i v-for="(r, k) in RAYS" :key="k" :style="{ '--dx': r.dx + 'px', '--dy': r.dy + 'px' }" />
    </span>
  </div>
</template>

<style scoped>
.text {
  position: relative;
  --fs: calc(clamp(20px, 3.7vh, 28px) * var(--fit, 1));
  font: 500 var(--fs) / 1.6 ui-monospace, 'JetBrains Mono', Consolas, monospace;
  white-space: pre-wrap;
  word-break: normal;
  color: var(--muted);
  padding: 16px 24px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 18px;
  min-height: 120px;
}
.viewport {
  overflow: hidden;
  max-height: calc(var(--fs) * 1.6 * 5);
}
.inner {
  position: relative;
  transition: transform 0.15s ease-out;
}
.ch.word {
  background: color-mix(in srgb, var(--accent) 10%, transparent);
}
.ch.ok {
  color: var(--ok);
}
.ch.skip {
  color: var(--map-chase);
  background: color-mix(in srgb, var(--map-chase) 22%, transparent);
  box-shadow: inset 0 -3px 0 var(--map-chase);
  border-radius: 3px;
}
.ch.bad {
  color: var(--err);
  background: color-mix(in srgb, var(--err) 22%, transparent);
  border-radius: 3px;
}
.ch.cur {
  color: var(--text);
  background: color-mix(in srgb, var(--accent) 22%, transparent);
  box-shadow: inset 0 -3px 0 var(--accent);
  border-radius: 3px;
}
.ch.cur.wrong {
  background: color-mix(in srgb, var(--err) 35%, transparent);
  box-shadow: inset 0 -3px 0 var(--err);
}
.caret {
  background: color-mix(in srgb, var(--accent) 22%, transparent);
  box-shadow: inset 0 -3px 0 var(--accent);
  border-radius: 3px;
}
.shake {
  animation: shake 0.2s;
}
@keyframes shake {
  25% {
    transform: translateX(-5px);
  }
  75% {
    transform: translateX(5px);
  }
}
.burst {
  position: absolute;
  width: 0;
  height: 0;
  pointer-events: none;
}
.burst i {
  position: absolute;
  left: -4px;
  top: -4px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
  animation: fly 0.7s ease-out forwards;
}
.burst i:nth-child(odd) {
  background: #fff;
}
@keyframes fly {
  to {
    transform: translate(var(--dx), var(--dy)) scale(0);
    opacity: 0;
  }
}
@media (max-width: 600px) {
  .text {
    font-size: 20px;
    padding: 16px;
  }
}
@media (prefers-reduced-motion: reduce) {
  .shake,
  .burst i {
    animation: none;
  }
}
</style>
