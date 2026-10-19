<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '@/api'
import type { Leaderboard } from '@/api'

type Metric = 'speed' | 'accuracy' | 'xp'
type Period = 'today' | 'week' | 'all'

const METRICS: { id: Metric; label: string }[] = [
  { id: 'speed', label: 'Скорость' },
  { id: 'accuracy', label: 'Точность' },
  { id: 'xp', label: 'Опыт' },
]
const PERIODS: { id: Period; label: string }[] = [
  { id: 'today', label: 'Сегодня' },
  { id: 'week', label: 'Неделя' },
  { id: 'all', label: 'Всё время' },
]

const metric = ref<Metric>('speed')
const period = ref<Period>('week')
const data = ref<Leaderboard | null>(null)
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    data.value = await api.getLeaderboard({ metric: metric.value, period: period.value, limit: 20 })
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch([metric, period], load)

const meOutside = computed(() => {
  const d = data.value
  return !!d?.me && !d.items.some((i) => i.user.id === d.me?.user.id)
})

function format(v: number) {
  if (metric.value === 'speed') return `${v.toFixed(1)} WPM`
  if (metric.value === 'accuracy') return `${v.toFixed(1)}%`
  return `${Math.round(v)} XP`
}
function medal(rank: number) {
  return rank === 1 ? 'gold' : rank === 2 ? 'silver' : rank === 3 ? 'bronze' : ''
}
</script>

<template>
  <main class="board">
    <h1>Таблица лидеров</h1>

    <div class="toolbar">
      <div class="seg" role="group" aria-label="Рейтинг">
        <button v-for="m in METRICS" :key="m.id" :aria-pressed="metric === m.id" @click="metric = m.id">
          {{ m.label }}
        </button>
      </div>
      <div class="seg" role="group" aria-label="Период">
        <button v-for="p in PERIODS" :key="p.id" :aria-pressed="period === p.id" @click="period = p.id">
          {{ p.label }}
        </button>
      </div>
    </div>

    <section class="panel" :aria-busy="loading">
      <ol v-if="data" class="list">
        <li
          v-for="r in data.items"
          :key="r.user.id"
          :class="{ me: r.user.id === data.me?.user.id, [medal(r.rank)]: true }"
        >
          <span class="rank">{{ r.rank }}</span>
          <span class="ava">{{ r.user.nickname.slice(0, 1) }}</span>
          <span class="name">{{ r.user.nickname }}<small>ур. {{ r.user.level }}</small></span>
          <b class="val">{{ format(r.value) }}</b>
        </li>
        <template v-if="meOutside && data.me">
          <li class="gap" aria-hidden="true">...</li>
          <li class="me">
            <span class="rank">{{ data.me.rank }}</span>
            <span class="ava">{{ data.me.user.nickname.slice(0, 1) }}</span>
            <span class="name">{{ data.me.user.nickname }}<small>ур. {{ data.me.user.level }}</small></span>
            <b class="val">{{ format(data.me.value) }}</b>
          </li>
        </template>
      </ol>
      <p v-if="data && !data.me" class="muted">
        {{ metric === 'accuracy' ? 'В рейтинг точности попадают игроки минимум с 3 тренировками за период.' : 'За этот период у тебя пока нет тренировок.' }}
        Потренируйся, чтобы попасть в таблицу.
      </p>
    </section>
  </main>
</template>

<style scoped>
.board {
  max-width: 760px;
  margin: 0 auto;
  padding: 16px;
}
h1 {
  margin: 0 0 14px;
}
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: space-between;
  margin-bottom: 12px;
}
.seg {
  display: inline-flex;
  background: var(--surface-2);
  border-radius: 12px;
  padding: 3px;
  gap: 3px;
}
.seg button {
  border: 0;
  background: transparent;
  color: var(--muted);
  padding: 7px 14px;
  border-radius: 9px;
  cursor: pointer;
  font-weight: 700;
}
.seg button[aria-pressed='true'] {
  background: var(--accent);
  color: var(--accent-ink);
}
.panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 10px;
}
.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 6px;
}
.list li {
  display: grid;
  grid-template-columns: 38px 40px 1fr auto;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: 12px;
  border: 2px solid transparent;
}
.list li.me {
  border-color: var(--accent);
  background: color-mix(in srgb, var(--accent) 10%, transparent);
}
.list li.gap {
  display: block;
  text-align: center;
  color: var(--muted);
  padding: 0;
}
.rank {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-weight: 800;
  color: var(--muted);
}
.gold .rank,
.silver .rank,
.bronze .rank {
  color: #1b1d24;
}
.gold .rank {
  background: #f5c542;
}
.silver .rank {
  background: #cbd5e1;
}
.bronze .rank {
  background: #d99058;
}
.ava {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--surface-2);
  font-weight: 800;
}
.name {
  display: grid;
  font-weight: 700;
}
.name small {
  color: var(--muted);
  font-weight: 400;
}
.val {
  font-size: 17px;
}
.muted {
  color: var(--muted);
  margin: 12px 8px 4px;
}
@media (max-width: 520px) {
  .list li {
    grid-template-columns: 32px 1fr auto;
  }
  .ava {
    display: none;
  }
}
</style>
