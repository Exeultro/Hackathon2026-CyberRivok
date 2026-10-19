<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api'
import type { Achievement } from '@/api'
import { useGameStore } from '@/stores/game'
import AchievementBadge from '@/components/AchievementBadge.vue'

const game = useGameStore()
const achievements = ref<Achievement[]>([])

onMounted(async () => {
  if (!game.loaded) await game.load()
  achievements.value = await api.getAchievements()
})

const sorted = computed(() =>
  [...achievements.value].sort((a, b) => Number(!!b.unlocked_at) - Number(!!a.unlocked_at)),
)
const done = computed(() => achievements.value.filter((a) => a.unlocked_at).length)
const stars = computed(() => ({
  total: game.maps.reduce((n, m) => n + (m.stars_total ?? 0), 0),
  max: game.maps.reduce((n, m) => n + (m.stars_max ?? 0), 0),
}))

function pct(a: Achievement) {
  return Math.min(100, (a.progress / a.target) * 100)
}
function shown(v: number) {
  return Number.isInteger(v) ? v : v.toFixed(1)
}
function date(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
}
</script>

<template>
  <main class="profile">
    <section v-if="game.user && game.progress" class="hero">
      <span class="ava">{{ game.user.nickname.slice(0, 1).toUpperCase() }}</span>
      <div class="info">
        <h1>{{ game.user.nickname }}</h1>
        <div class="lvl">
          <b>Уровень {{ game.progress.level }}</b>
          <small>{{ game.progress.xp_in_level }} / {{ game.progress.xp_to_next }} XP</small>
        </div>
        <div class="bar" role="progressbar" :aria-valuenow="game.xpPercent" aria-valuemin="0" aria-valuemax="100">
          <i :style="{ width: game.xpPercent + '%' }" />
        </div>
      </div>
    </section>

    <section v-if="game.summary" class="facts">
      <div><b>{{ game.summary.streak_days }}</b><small>дней подряд</small></div>
      <div><b>{{ Math.round(game.summary.best_wpm) }}</b><small>лучший WPM</small></div>
      <div><b>{{ stars.total }} / {{ stars.max }}</b><small>звёзд на картах</small></div>
      <div><b>{{ game.summary.sessions_count }}</b><small>тренировок</small></div>
    </section>

    <section class="panel">
      <h2>Карты</h2>
      <ul class="maps">
        <li v-for="m in game.maps" :key="m.id" :style="{ '--c': `var(--map-${m.id})` }">
          <span class="mname">{{ m.title }}</span>
          <div class="mbar"><i :style="{ width: ((m.stars_total ?? 0) / (m.stars_max || 1)) * 100 + '%' }" /></div>
          <small>{{ m.stars_total }} / {{ m.stars_max }}</small>
        </li>
      </ul>
    </section>

    <section class="panel">
      <h2>Достижения <small>{{ done }} из {{ achievements.length }}</small></h2>
      <ul class="grid">
        <li v-for="a in sorted" :key="a.code" :class="{ off: !a.unlocked_at }">
          <AchievementBadge :icon="a.icon" :unlocked="!!a.unlocked_at" :size="52" />
          <div class="txt">
            <b>{{ a.title }}</b>
            <small>{{ a.description }}</small>
            <small v-if="a.unlocked_at" class="got">Получено {{ date(a.unlocked_at) }}</small>
            <template v-else>
              <div class="abar"><i :style="{ width: pct(a) + '%' }" /></div>
              <small>{{ shown(a.progress) }} / {{ a.target }}</small>
            </template>
          </div>
        </li>
      </ul>
    </section>
  </main>
</template>

<style scoped>
.profile {
  max-width: 880px;
  margin: 0 auto;
  padding: 16px;
}
.hero {
  display: flex;
  gap: 18px;
  align-items: center;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 18px;
  padding: 18px;
  margin-bottom: 12px;
}
.ava {
  flex: none;
  width: 84px;
  height: 84px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-size: 40px;
  font-weight: 800;
  background: var(--accent);
  color: var(--accent-ink);
  box-shadow: 0 6px 0 color-mix(in srgb, var(--accent) 58%, #000);
}
.info {
  flex: 1;
  min-width: 0;
}
h1 {
  margin: 0;
}
.lvl {
  display: flex;
  justify-content: space-between;
  margin: 6px 0 4px;
}
small {
  color: var(--muted);
}
.bar {
  height: 12px;
  border-radius: 6px;
  background: var(--surface-2);
  overflow: hidden;
}
.bar i {
  display: block;
  height: 100%;
  background: var(--accent);
  border-radius: 6px;
}
.facts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}
.facts div {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 14px;
  padding: 10px 14px;
}
.facts b {
  display: block;
  font-size: 22px;
}
.panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 12px;
}
h2 {
  margin: 0 0 12px;
  font-size: 16px;
}
h2 small {
  font-weight: 400;
  margin-left: 8px;
}
.maps {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}
.maps li {
  display: grid;
  grid-template-columns: 110px 1fr 56px;
  align-items: center;
  gap: 12px;
}
.mname {
  font-weight: 700;
  color: var(--c);
}
.mbar {
  height: 10px;
  border-radius: 5px;
  background: var(--surface-2);
  overflow: hidden;
}
.mbar i {
  display: block;
  height: 100%;
  background: var(--c);
  border-radius: 5px;
}
.grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 10px;
}
.grid li {
  display: flex;
  gap: 12px;
  align-items: center;
  background: var(--surface-2);
  border-radius: 14px;
  padding: 10px 12px;
}
.grid li.off {
  opacity: 0.85;
}
.txt {
  display: grid;
  gap: 2px;
  min-width: 0;
}
.got {
  color: var(--ok);
}
.abar {
  height: 6px;
  border-radius: 3px;
  background: var(--border);
  overflow: hidden;
  margin-top: 4px;
}
.abar i {
  display: block;
  height: 100%;
  background: var(--accent);
}
@media (max-width: 520px) {
  .maps li {
    grid-template-columns: 90px 1fr 48px;
  }
  .ava {
    width: 64px;
    height: 64px;
    font-size: 30px;
  }
}
</style>
