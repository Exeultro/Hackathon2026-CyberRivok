<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'
import type { GameMap, MapId } from '@/api'

const props = defineProps<{ maps: GameMap[]; modelValue: MapId }>()
const emit = defineEmits<{ 'update:modelValue': [id: MapId] }>()

const root = ref<HTMLElement | null>(null)

function revealActive() {
  root.value?.querySelector('[aria-selected="true"]')?.scrollIntoView({ inline: 'center', block: 'nearest' })
}

onMounted(revealActive)
watch(() => props.modelValue, () => nextTick(revealActive))
</script>

<template>
  <div ref="root" class="tabs" role="tablist">
    <button
      v-for="m in maps"
      :key="m.id"
      class="tab"
      role="tab"
      :aria-selected="m.id === modelValue"
      :style="{ '--c': `var(--map-${m.id})` }"
      @click="emit('update:modelValue', m.id)"
    >
      {{ m.title }}
    </button>
  </div>
</template>

<style scoped>
.tabs {
  display: flex;
  gap: 6px;
}
.tab {
  height: 32px;
  padding: 0 14px;
  font-weight: 700;
  font-size: 13px;
  border-radius: 999px;
  cursor: pointer;
  background: transparent;
  color: var(--muted);
  border: 2px solid var(--border);
  white-space: nowrap;
}
.tab:hover {
  color: var(--text);
}
.tab[aria-selected='true'] {
  color: #fff;
  background: var(--c);
  border-color: var(--c);
}
@media (max-width: 700px) {
  .tab {
    padding: 0 8px;
    font-size: 12px;
  }
}
</style>
