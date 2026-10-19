<script setup lang="ts">
import type { GameMap, MapNode } from '@/api'
import StarIcon from './StarIcon.vue'

defineProps<{ map: GameMap; node: MapNode | null }>()
const emit = defineEmits<{ start: []; close: [] }>()
</script>

<template>
  <div class="sheet" :class="{ open: node }" aria-live="polite">
    <div v-if="node" class="inner" :style="{ '--c': `var(--map-${map.id})` }">
      <button class="x" aria-label="Закрыть" @click="emit('close')">×</button>
      <h3>Узел {{ node.order }}: {{ node.title }}</h3>
      <div class="meta">{{ map.title }} - сложность {{ node.difficulty }} из 5</div>
      <ul class="th">
        <li v-for="(t, k) in node.star_thresholds" :key="k">
          <span class="st"><StarIcon v-for="n in k + 1" :key="n" :filled="true" :size="20" /></span>
          от {{ t.min_wpm }} WPM и точность от {{ t.min_accuracy }}%
        </li>
      </ul>
      <button class="btn btn-primary go" @click="emit('start')">Начать</button>
    </div>
  </div>
</template>

<style scoped>
.sheet {
  position: fixed;
  left: var(--nav-w);
  right: 0;
  bottom: var(--nav-b);
  background: var(--surface);
  border-top: 3px solid var(--border);
  padding: 16px;
  transform: translateY(calc(100% + var(--nav-b)));
  transition: transform 0.2s;
  z-index: 10;
}
.sheet.open {
  transform: none;
  border-top-color: var(--c, var(--border));
}
.inner {
  max-width: 520px;
  margin: 0 auto;
  position: relative;
}
h3 {
  margin: 0;
  font-size: 18px;
}
.meta {
  color: var(--muted);
  font-size: 14px;
  margin: 2px 0 10px;
}
.th {
  list-style: none;
  margin: 0 0 14px;
  padding: 0;
  font-size: 14px;
  display: grid;
  gap: 4px;
}
.th li {
  display: flex;
  align-items: center;
  gap: 10px;
}
.st {
  display: inline-flex;
  min-width: 66px;
}
.go {
  width: 100%;
}
.x {
  position: absolute;
  right: 0;
  top: -4px;
  background: none;
  border: 0;
  color: var(--muted);
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
}
</style>
