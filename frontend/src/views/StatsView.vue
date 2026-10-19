<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '@/api'
import type { KeyStat, SessionSummary, StatsRange, StatsSummary, TimelinePoint } from '@/api'
import LineChart from '@/components/LineChart.vue'

const MODE_NAME: Record<string, string> = {
  words: 'Слова',
  sentences: 'Тексты',
  code: 'Код',
  symbols: 'Символы',
}
const RANGES: { id: StatsRange; label: string }[] = [
  { id: 'week', label: 'Неделя' },
  { id: 'month', label: 'Месяц' },
  { id: 'year', label: 'Год' },
]

const summary = ref<StatsSummary | null>(null)
const timeline = ref<TimelinePoint[]>([])
const keys = ref<KeyStat[]>([])
const sessions = ref<SessionSummary[]>([])
const range = ref<StatsRange>('week')
const kind = ref<'key' | 'bigram'>('key')
const asTable = ref(false)

async function loadTimeline() {
  timeline.value = await api.getTimeline(range.value)
}
async function loadKeys() {
  keys.value = await api.getKeys(kind.value, 10)
}

onMounted(async () => {
  const [s, list] = await Promise.all([api.getSummary(), api.getSessions(8)])
  summary.value = s
  sessions.value = list
  await Promise.all([loadTimeline(), loadKeys()])
})
watch(range, loadTimeline)
watch(kind, loadKeys)

const labels = computed(() =>
  timeline.value.map((p) =>
    new Date(p.date + 'T00:00:00').toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' }),
  ),
)
const speed = computed(() => timeline.value.map((p) => p.avg_wpm))
const accuracy = computed(() => timeline.value.map((p) => p.avg_accuracy))
const worst = computed(() => Math.max(0.01, ...keys.value.map((k) => k.error_rate)))

function duration(ms: number) {
  const min = Math.round(ms / 60000)
  return min >= 60 ? `${Math.floor(min / 60)} ч ${min % 60} мин` : `${min} мин`
}
function when(iso: string) {
  return new Date(iso).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
function show(key: string) {
  return key === ' ' ? '␣' : key
}
</script>

<template>
  <main class="stats">
    <h1>Статистика</h1>

    <section v-if="summary" class="cards">
      <div class="card"><b>{{ summary.sessions_count }}</b><small>тренировок</small></div>
      <div class="card"><b>{{ duration(summary.total_time_ms) }}</b><small>всего времени</small></div>
      <div class="card"><b>{{ Math.round(summary.best_wpm) }}</b><small>лучший WPM</small></div>
      <div class="card"><b>{{ Math.round(summary.avg_wpm) }}</b><small>средний WPM</small></div>
      <div class="card"><b>{{ Math.round(summary.avg_accuracy) }}%</b><small>средняя точность</small></div>
      <div class="card"><b>{{ summary.streak_days }}</b><small>дней подряд</small></div>
    </section>

    <div class="toolbar">
      <div class="seg" role="group" aria-label="Период">
        <button v-for="r in RANGES" :key="r.id" :aria-pressed="range === r.id" @click="range = r.id">
          {{ r.label }}
        </button>
      </div>
      <button class="seg-solo" :aria-pressed="asTable" @click="asTable = !asTable">Таблица</button>
    </div>

    <section v-if="!asTable" class="charts">
      <div class="panel">
        <h2>Скорость, WPM</h2>
        <LineChart :labels="labels" :values="speed" title="Скорость" unit=" WPM" color="--accent" :y-min="0" />
      </div>
      <div class="panel">
        <h2>Точность, %</h2>
        <LineChart
          :labels="labels"
          :values="accuracy"
          title="Точность"
          unit="%"
          color="--map-search"
          :y-min="80"
          :y-max="100"
        />
      </div>
    </section>

    <section v-else class="panel">
      <h2>Динамика по дням</h2>
      <table>
        <thead>
          <tr><th>День</th><th>Тренировок</th><th>WPM</th><th>Точность</th><th>Ошибок</th></tr>
        </thead>
        <tbody>
          <tr v-for="(p, i) in timeline" :key="p.date">
            <td>{{ labels[i] }}</td><td>{{ p.sessions }}</td><td>{{ p.avg_wpm }}</td>
            <td>{{ p.avg_accuracy }}%</td><td>{{ p.errors }}</td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="panel">
      <div class="head">
        <h2>Что даётся тяжелее всего</h2>
        <div class="seg" role="group" aria-label="Что показывать">
          <button :aria-pressed="kind === 'key'" @click="kind = 'key'">Клавиши</button>
          <button :aria-pressed="kind === 'bigram'" @click="kind = 'bigram'">Сочетания</button>
        </div>
      </div>
      <ul v-if="keys.length" class="keys">
        <li v-for="k in keys" :key="k.key">
          <span class="cap">{{ show(k.key) }}</span>
          <div class="track" :title="`${k.errors} ошибок из ${k.hits} нажатий`">
            <i :style="{ width: (k.error_rate / worst) * 100 + '%' }" />
          </div>
          <span class="pct">{{ Math.round(k.error_rate * 100) }}%</span>
          <small>{{ k.errors }} из {{ k.hits }}</small>
        </li>
      </ul>
      <p v-else class="muted">Пока мало данных. Сделай несколько тренировок.</p>
      <p class="muted note">Доля ошибок на клавише или сочетании от всех нажатий. Тренируй верхние.</p>
    </section>

    <section class="panel">
      <h2>Последние тренировки</h2>
      <table v-if="sessions.length">
        <thead>
          <tr><th>Когда</th><th>Режим</th><th>WPM</th><th>Точность</th><th>Ошибок</th><th>XP</th></tr>
        </thead>
        <tbody>
          <tr v-for="s in sessions" :key="s.id">
            <td>{{ when(s.started_at) }}</td><td>{{ MODE_NAME[s.mode] }}</td><td>{{ Math.round(s.wpm) }}</td>
            <td>{{ Math.round(s.accuracy) }}%</td><td>{{ s.errors }}</td><td>+{{ s.xp_gained }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </main>
</template>

<style scoped>
.stats {
  max-width: 980px;
  margin: 0 auto;
  padding: 16px;
}
h1 {
  margin: 0 0 14px;
}
h2 {
  margin: 0 0 10px;
  font-size: 16px;
}
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 10px;
  margin-bottom: 16px;
}
.card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 14px;
  padding: 12px 14px;
}
.card b {
  display: block;
  font-size: 24px;
  line-height: 1.2;
}
.card small,
.muted {
  color: var(--muted);
}
.toolbar {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 12px;
}
.seg {
  display: inline-flex;
  background: var(--surface-2);
  border-radius: 12px;
  padding: 3px;
  gap: 3px;
}
.seg button,
.seg-solo {
  border: 0;
  background: transparent;
  color: var(--muted);
  padding: 6px 14px;
  border-radius: 9px;
  cursor: pointer;
  font-weight: 600;
}
.seg-solo {
  border: 1px solid var(--border);
  margin-left: auto;
}
.seg button[aria-pressed='true'] {
  background: var(--accent);
  color: var(--accent-ink);
}
.seg-solo[aria-pressed='true'] {
  color: var(--accent);
  border-color: var(--accent);
}
.charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}
.panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 12px;
}
.charts .panel {
  margin-bottom: 0;
}
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}
.head h2 {
  margin: 0;
}
.keys {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 8px;
}
.keys li {
  display: grid;
  grid-template-columns: 44px 1fr 44px 90px;
  align-items: center;
  gap: 10px;
}
.cap {
  height: 36px;
  display: grid;
  place-items: center;
  border-radius: 8px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-bottom-width: 3px;
  font: 700 17px ui-monospace, Consolas, monospace;
}
.track {
  height: 10px;
  border-radius: 5px;
  background: var(--surface-2);
  overflow: hidden;
}
.track i {
  display: block;
  height: 100%;
  border-radius: 5px;
  background: var(--err);
}
.pct {
  font-weight: 800;
  text-align: right;
}
.keys small {
  color: var(--muted);
  white-space: nowrap;
}
.note {
  margin: 10px 0 0;
  font-size: 13px;
}
table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}
th,
td {
  text-align: left;
  padding: 8px 6px;
  border-bottom: 1px solid var(--border);
}
th {
  color: var(--muted);
  font-weight: 600;
}
@media (max-width: 560px) {
  .keys li {
    grid-template-columns: 38px 1fr 40px;
  }
  .keys small {
    display: none;
  }
}
</style>
