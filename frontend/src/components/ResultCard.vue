<script setup lang="ts">
import { computed } from 'vue'
import type { MapNode, SessionResult } from '@/api'
import StarIcon from './StarIcon.vue'
import AchievementBadge from './AchievementBadge.vue'

const props = defineProps<{ result: SessionResult; node: MapNode | null; backLabel?: string }>()
const emit = defineEmits<{ retry: []; map: []; next: [] }>()

const nr = computed(() => props.result.node_result)
const stars = computed(() => nr.value?.stars ?? 0)
const hasStars = computed(() => !!nr.value)
const passed = computed(() => nr.value?.passed ?? true)
const levelUp = computed(() => props.result.progress.level > props.result.level_before)
const need = computed(() => props.node?.star_thresholds[0])
</script>

<template>
  <div class="overlay">
    <div class="card" role="dialog" aria-modal="true">
      <h2>{{ !hasStars ? 'Тренировка завершена' : passed ? 'Уровень пройден!' : 'Не хватило немного' }}</h2>
      <div v-if="hasStars" class="stars">
        <span v-for="i in 3" :key="i" class="st" :style="{ animationDelay: i * 0.18 + 's' }">
          <StarIcon :filled="i <= stars" :size="i === 2 ? 64 : 52" />
        </span>
      </div>
      <p v-if="hasStars && !passed && need" class="need">Нужно: от {{ need.min_wpm }} WPM и точность от {{ need.min_accuracy }}%</p>
      <div class="grid">
        <div><b>{{ Math.round(result.wpm) }}</b><small>WPM</small></div>
        <div><b>{{ Math.round(result.cpm) }}</b><small>симв/мин</small></div>
        <div><b>{{ Math.round(result.accuracy) }}%</b><small>точность</small></div>
        <div><b>{{ result.errors }}</b><small>ошибок</small></div>
        <div><b>{{ result.score }}</b><small>очки</small></div>
        <div><b>+{{ result.xp_gained }}</b><small>XP</small></div>
      </div>
      <p v-if="result.is_personal_best" class="flash">Новый личный рекорд!</p>
      <p v-if="levelUp" class="flash up">Новый уровень героя: {{ result.progress.level }}</p>
      <ul v-if="result.new_achievements.length" class="achv">
        <li v-for="a in result.new_achievements" :key="a.code">
          <AchievementBadge :icon="a.icon" :unlocked="true" :size="40" />
          <span><b>{{ a.title }}</b><small>{{ a.description }}</small></span>
        </li>
      </ul>
      <div class="actions">
        <button class="btn btn-ghost" @click="emit('map')">{{ (backLabel ?? '← Карта').replace('← ', '') }}</button>
        <button class="btn" @click="emit('retry')">Ещё раз</button>
        <button v-if="passed && nr?.next_node_id" class="btn btn-primary" @click="emit('next')">Дальше</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 30;
  display: grid;
  place-items: center;
  background: rgb(0 0 0 / 0.6);
  padding: 16px;
}
.card {
  width: min(460px, 100%);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 24px;
  text-align: center;
}
h2 {
  margin: 0 0 8px;
}
.stars {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  gap: 6px;
  margin: 8px 0 16px;
}
.st {
  display: inline-flex;
  animation: pop 0.5s both;
}
@keyframes pop {
  from {
    transform: scale(0) rotate(-40deg);
  }
  60% {
    transform: scale(1.25) rotate(8deg);
  }
  to {
    transform: scale(1) rotate(0);
  }
}
.need {
  color: var(--muted);
  margin: 0 0 12px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 14px;
}
.grid div {
  background: var(--surface-2);
  border-radius: 12px;
  padding: 10px 4px;
}
.grid b {
  display: block;
  font-size: 22px;
}
.grid small {
  color: var(--muted);
}
.flash {
  margin: 6px 0;
  font-weight: 800;
  color: var(--accent);
}
.flash.up {
  color: var(--ok);
}
.achv {
  list-style: none;
  padding: 0;
  margin: 10px 0 0;
  display: grid;
  gap: 8px;
  text-align: left;
}
.achv li {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--surface-2);
  border-radius: 12px;
  padding: 8px 10px;
}
.achv span {
  display: grid;
}
.achv small {
  color: var(--muted);
}
.actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 14px;
}
@media (prefers-reduced-motion: reduce) {
  .st {
    animation: none;
  }
}
</style>
