<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import TypingSession from '@/components/TypingSession.vue'
import type { SessionConfig } from '@/components/TypingSession.vue'

const props = defineProps<{ nodeId: string }>()
const game = useGameStore()
const router = useRouter()

onMounted(() => {
  if (!game.loaded) game.load()
})

const found = computed(() => {
  for (const map of game.maps) {
    const node = map.nodes.find((n) => n.id === props.nodeId)
    if (node) return { map, node }
  }
  return null
})

const config = computed<SessionConfig | null>(() => {
  const f = found.value
  if (!f) return null
  return {
    mode: f.node.mode,
    difficulty: f.node.difficulty,
    count: f.node.text_count ?? 3,
    code_language: f.node.code_language,
    nodeId: f.node.id,
    title: `${f.map.title} - ${f.node.title}`,
  }
})
</script>

<template>
  <main>
    <TypingSession
      v-if="config"
      :key="nodeId"
      :config="config"
      :node="found?.node ?? null"
      back-label="← Карта"
      @exit="router.push({ name: 'map' })"
      @next="router.replace({ name: 'play', params: { nodeId: $event } })"
    />
    <p v-else class="muted">Загрузка...</p>
  </main>
</template>

<style scoped>
.muted {
  color: var(--muted);
  padding: 16px;
}
</style>
