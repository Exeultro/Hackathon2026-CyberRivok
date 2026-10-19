<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
const sampleText = 'Пар зима хлеб вы пар зима сад лето жук луг шар'
const chars = sampleText.split('')
const typed = ref([])
const started = ref(false)
const startedAt = ref(null)
const now = ref(Date.now())
let timer = null
onMounted(() => { timer = setInterval(() => { now.value = Date.now() }, 200) })
onUnmounted(() => clearInterval(timer))
const currentIndex = computed(() => typed.value.length)
const correctCount = computed(() => typed.value.reduce((acc, ch, i) => acc + (ch === chars[i] ? 1 : 0), 0))
const accuracy = computed(() => !typed.value.length ? 100 : Math.round((correctCount.value / typed.value.length) * 100))
const elapsedMs = computed(() => !startedAt.value ? 0 : now.value - startedAt.value)
const wpm = computed(() => !elapsedMs.value ? 0 : Math.round((correctCount.value / 5) / (elapsedMs.value / 60000)) || 0)
function onKey(e) {
  if (e.key === 'Backspace') { typed.value.pop(); return }
  if (e.key.length !== 1 || typed.value.length >= chars.length) return
  if (!started.value) { started.value = true; startedAt.value = Date.now() }
  typed.value.push(e.key)
}
onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))
function charClass(i) {
  if (i < typed.value.length) return typed.value[i] === chars[i] ? 'is-correct' : 'is-error'
  if (i === currentIndex.value) return 'is-current'
  return ''
}
</script>

<template>
  <div class="demo">
    <div class="demo__stats">
      <div class="stat"><span class="stat__value">{{ wpm }}</span><span class="stat__label">WPM</span></div>
      <div class="stat"><span class="stat__value">{{ accuracy }}%</span><span class="stat__label">точность</span></div>
      <div class="stat"><span class="stat__value">{{ correctCount * 10 }}</span><span class="stat__label">очки</span></div>
      <div class="stat stat--right"><span class="stat__value">×1</span><span class="stat__label">комбо</span></div>
    </div>
    <div class="demo__text" :class="{ 'is-active': started }">
      <span v-for="(ch, i) in chars" :key="i" class="char" :class="charClass(i)">{{ ch }}</span>
    </div>
    <p v-if="!started" class="demo__hint">Нажми любую клавишу, чтобы начать</p>
    <p v-else class="demo__hint">Скорость {{ wpm }} WPM · Точность {{ accuracy }}%</p>
  </div>
</template>

<style scoped>
.demo { background: var(--bg-card); border-radius: var(--radius-lg); padding: 24px 28px 28px; box-shadow: var(--shadow-card); }
.demo__stats { display: flex; align-items: flex-end; gap: 32px; padding-bottom: 16px; margin-bottom: 20px; border-bottom: 1px solid rgba(255,255,255,.06); }
.stat { display: flex; flex-direction: column; gap: 2px; }
.stat--right { margin-left: auto; text-align: right; }
.stat__value { font-size: 26px; font-weight: 700; line-height: 1; }
.stat__label { font-size: 12px; color: var(--text-muted); text-transform: uppercase; letter-spacing: .05em; }
.demo__text { font-family: var(--font-mono); font-size: 26px; line-height: 1.7; letter-spacing: .02em; user-select: none; }
.char { color: var(--text-muted); transition: color .1s; }
.char.is-correct { color: var(--success); }
.char.is-error { color: var(--error); text-decoration: underline; }
.char.is-current { color: var(--accent); background: rgba(245,165,36,.12); border-radius: 4px; box-shadow: inset 0 -2px 0 var(--accent); }
.demo__hint { margin-top: 20px; text-align: center; font-size: 13px; color: var(--text-muted); }
</style>