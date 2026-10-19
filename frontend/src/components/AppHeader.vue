<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import MapTabs from './MapTabs.vue'
import { useGameStore } from '@/stores/game'
import { useTheme } from '@/useTheme'

const game = useGameStore()
const { theme, toggle } = useTheme()
const route = useRoute()

onMounted(() => {
  if (!game.loaded) game.load()
})
</script>

<template>
  <header class="head">
    <RouterLink to="/" class="logo">Печать героя</RouterLink>
    <div class="center">
      <MapTabs v-if="route.name === 'map' && game.loaded" v-model="game.currentMapId" :maps="game.maps" />
    </div>
    <div class="lang" role="group" aria-label="Язык набора">
      <button
        v-for="l in (['ru', 'en'] as const)"
        :key="l"
        :aria-pressed="game.lang === l"
        @click="game.lang = l"
      >
        {{ l.toUpperCase() }}
      </button>
    </div>
    <div v-if="game.summary" class="chip" title="Дней подряд">
      <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
        <path
          d="M12 2c1 4 5 6 5 11a5 5 0 0 1-10 0c0-2 1-3 2-4 0 2 1 3 2 3 0-4-1-6 1-10z"
          fill="var(--map-chase)"
        />
      </svg>
      {{ game.summary.streak_days }}
    </div>
    <div v-if="game.progress" class="level">
      <span>ур. {{ game.progress.level }}</span>
      <div
        class="bar"
        role="progressbar"
        :aria-valuenow="game.progress.xp_in_level"
        :aria-valuemax="game.progress.xp_to_next"
      >
        <i :style="{ width: game.xpPercent + '%' }" />
      </div>
      <small>{{ game.progress.xp_in_level }} / {{ game.progress.xp_to_next }}</small>
    </div>
    <button class="theme" @click="toggle">{{ theme === 'dark' ? 'Светлая' : 'Тёмная' }}</button>
    <RouterLink to="/profile" class="avatar" aria-label="Профиль">{{ (game.user?.nickname ?? 'Г').slice(0, 1).toUpperCase() }}</RouterLink>
  </header>
</template>

<style scoped>
.head {
  position: fixed;
  inset: 0 0 auto 0;
  height: var(--head-h);
  z-index: 20;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 0 12px 0 16px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
}
.logo {
  font-weight: 800;
  font-size: 18px;
  color: var(--accent);
  text-decoration: none;
  white-space: nowrap;
}
.center {
  flex: 1;
  min-width: 0;
  overflow-x: auto;
  scrollbar-width: none;
}
.center :deep(.tabs) {
  width: max-content;
  margin: 0 auto;
}
.chip,
.level,
.avatar,
.lang {
  flex-shrink: 0;
}
.lang {
  display: flex;
  background: var(--surface-2);
  border-radius: 10px;
  padding: 2px;
}
.lang button {
  border: 0;
  background: transparent;
  color: var(--muted);
  font-weight: 800;
  font-size: 12px;
  padding: 4px 9px;
  border-radius: 8px;
  cursor: pointer;
}
.lang button[aria-pressed='true'] {
  background: var(--accent);
  color: var(--accent-ink);
}
.chip {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 800;
}
.level {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 700;
  font-size: 14px;
  width: 240px;
}
.level small {
  color: var(--muted);
  font-weight: 400;
  white-space: nowrap;
}
.bar {
  flex: 1;
  height: 10px;
  min-width: 40px;
  border-radius: 5px;
  background: var(--surface-2);
  overflow: hidden;
}
.bar i {
  display: block;
  height: 100%;
  background: var(--accent);
  border-radius: 5px;
  transition: width 0.4s;
}
.theme {
  background: none;
  border: 1px solid var(--border);
  color: var(--muted);
  border-radius: 10px;
  padding: 4px 10px;
  cursor: pointer;
  font-size: 13px;
}
.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--accent);
  color: var(--accent-ink);
  font-weight: 800;
  text-decoration: none;
}
@media (max-width: 900px) {
  .level small,
  .theme {
    display: none;
  }
  .level {
    width: 90px;
  }
}
@media (max-width: 760px) {
  .logo,
  .level .bar {
    display: none;
  }
}
@media (max-width: 600px) {
  .chip,
  .level {
    display: none;
  }
}
</style>
