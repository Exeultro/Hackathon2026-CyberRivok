<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import type { MapId, MapNode } from '@/api'
import MapPath from '@/components/MapPath.vue'
import NodeSheet from '@/components/NodeSheet.vue'
import NotReadyModal from '@/components/NotReadyModal.vue'
import Terrain from '@/components/Terrain.vue'
import Wind from '@/components/Wind.vue'
import { pathHeight } from '@/geometry'
import { SCENES, sceneVars } from '@/scenes'

const MAP_IDS: MapId[] = ['journey', 'chase', 'search', 'siege']

const game = useGameStore()
const router = useRouter()
const route = useRoute()

onMounted(() => document.documentElement.classList.add('no-scrollbar'))
onBeforeUnmount(() => document.documentElement.classList.remove('no-scrollbar'))

const queryMap = MAP_IDS.find((id) => id === route.query.map)
if (queryMap) game.currentMapId = queryMap

const newSeed = () => Math.random().toString(36).slice(2, 10)
const seed = ref(newSeed())
const selected = ref<MapNode | null>(null)
const warning = ref(false)

watch(
  () => game.currentMapId,
  () => {
    selected.value = null
    seed.value = newSeed()
  },
)

const map = computed(() => game.maps.find((m) => m.id === game.currentMapId) ?? null)
const scene = computed(() => SCENES[game.currentMapId])
const sceneHeight = computed(() => pathHeight(map.value?.nodes.length ?? 0) + 160)

function go() {
  if (!selected.value) return
  router.push({ name: 'play', params: { nodeId: selected.value.id } })
}

function start() {
  if (map.value && map.value.ready === false) warning.value = true
  else go()
}
</script>

<template>
  <main v-if="game.loaded && map">
    <section class="scene" :style="{ ...sceneVars(scene), height: sceneHeight + 'px' }">
      <Terrain :key="game.currentMapId" :scene="scene" :height="sceneHeight" :seed="seed" :node-count="map.nodes.length" />
      <div class="col">
        <MapPath :map="map" :selected-id="selected?.id ?? null" @select="selected = $event" />
      </div>
      <Wind :key="'w' + seed" :kind="scene.wind.kind" :count="scene.wind.count" :height="sceneHeight" :seed="seed" />
    </section>
    <NodeSheet :map="map" :node="selected" @start="start" @close="selected = null" />
    <NotReadyModal
      :open="warning"
      :recommended-level="map.recommended_level"
      :level="game.progress?.level ?? 1"
      @cancel="warning = false"
      @confirm="((warning = false), go())"
    />
  </main>
  <p v-else-if="game.loadError" class="muted">
    Нет связи с сервером. <button @click="game.load()">Повторить</button>
  </p>
  <p v-else class="muted">Загрузка...</p>
</template>

<style scoped>
.scene {
  position: relative;
  overflow: hidden;
}
.col {
  position: relative;
  z-index: 1;
  max-width: 640px;
  margin: 0 auto;
  padding: 0 16px;
}
.muted {
  color: var(--muted);
  padding: 16px;
}
</style>
