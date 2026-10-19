<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import type { GameMap, MapNode, Mode } from '@/api'
import NodeStars from './NodeStars.vue'
import { PATH_W, controls, nodePoint, pathHeight } from '@/geometry'
import type { Pt } from '@/geometry'
import { SCENES } from '@/scenes'

const props = defineProps<{ map: GameMap; selectedId: string | null }>()
const emit = defineEmits<{ select: [node: MapNode] }>()

const GLYPH: Record<Mode, string> = {
  words: 'Аа',
  sentences: '¶',
  code: '</>',
  symbols: '#@',
}

const root = ref<HTMLElement | null>(null)
const width = ref(PATH_W)
let observer: ResizeObserver | null = null

onMounted(() => {
  if (!root.value) return
  width.value = root.value.clientWidth || PATH_W
  observer = new ResizeObserver(([entry]) => {
    width.value = entry.contentRect.width
  })
  observer.observe(root.value)
})
onBeforeUnmount(() => observer?.disconnect())

const shape = computed(() => SCENES[props.map.id].path)
const scale = computed(() => width.value / PATH_W)
const points = computed<Pt[]>(() =>
  props.map.nodes.map((_, i) => {
    const p = nodePoint(i, shape.value)
    return { x: p.x * scale.value, y: p.y }
  }),
)
const height = computed(() => pathHeight(props.map.nodes.length))
const currentIndex = computed(() => props.map.nodes.findIndex((n) => n.stars === 0))

const linePath = computed(() =>
  points.value
    .map((p, i, all) => {
      if (i === 0) return `M${p.x},${p.y}`
      const [c1, c2] = controls(all[i - 1], p)
      return `C${c1.x},${c1.y} ${c2.x},${c2.y} ${p.x},${p.y}`
    })
    .join(' '),
)

function nodeClass(i: number, node: MapNode) {
  return {
    cur: i === currentIndex.value,
    todo: node.stars === 0 && i !== currentIndex.value,
    sel: node.id === props.selectedId,
  }
}
</script>

<template>
  <div ref="root" class="path" :style="{ height: height + 'px', '--c': `var(--map-${map.id})` }">
    <svg class="line" :width="width" :height="height" :viewBox="`0 0 ${width} ${height}`" aria-hidden="true">
      <path :d="linePath" class="trail" />
    </svg>
    <div
      v-for="(node, i) in map.nodes"
      :key="node.id"
      class="n"
      :style="{ left: points[i].x + 'px', top: points[i].y + 'px' }"
    >
      <NodeStars v-if="node.stars > 0" :count="node.stars" />
      <div v-if="i === currentIndex" class="bubble">Старт</div>
      <button
        class="node"
        :class="nodeClass(i, node)"
        :aria-label="`Узел ${node.order}: ${node.title}`"
        @click="emit('select', node)"
      >
        {{ GLYPH[node.mode] }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.path {
  position: relative;
}
.line {
  position: absolute;
  inset: 0;
}
.trail {
  fill: none;
  stroke: var(--trail);
  stroke-width: 42;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.n {
  position: absolute;
  width: 78px;
  height: 66px;
  transform: translate(-50%, -50%);
}
.n::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 78px;
  height: 66px;
  transform: translateY(18px);
  border-radius: 50%;
  background: rgb(0 0 0 / 0.25);
}
.node {
  position: relative;
  width: 78px;
  height: 66px;
  padding: 0;
  border-radius: 50%;
  border: 0;
  cursor: pointer;
  font-weight: 800;
  font-size: 22px;
  letter-spacing: -0.02em;
  background: var(--c);
  color: #fff;
  box-shadow: 0 10px 0 color-mix(in srgb, var(--c) 58%, #000);
  transition: transform 0.08s;
}
.node:active {
  transform: translateY(5px);
  box-shadow: 0 5px 0 color-mix(in srgb, var(--c) 58%, #000);
}
.node.cur {
  animation: pulse 1.6s infinite;
}
.node.todo {
  background: #b9bcc6;
  color: #8a8e9b;
  box-shadow: 0 10px 0 #8d909c;
}
.bubble {
  position: absolute;
  left: 50%;
  bottom: calc(100% + 18px);
  transform: translateX(-50%);
  padding: 4px 12px;
  border-radius: 10px;
  background: #fff;
  color: var(--c);
  font-weight: 800;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  box-shadow: 0 3px 0 rgb(0 0 0 / 0.2);
  animation: bob 1.6s infinite;
  z-index: 2;
}
.bubble::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 100%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-top-color: #fff;
}
@keyframes bob {
  50% {
    transform: translateX(-50%) translateY(-4px);
  }
}
.node.sel {
  outline: 3px solid var(--accent);
  outline-offset: 5px;
}
@keyframes pulse {
  50% {
    transform: scale(1.07);
  }
}
@media (prefers-reduced-motion: reduce) {
  .node.cur,
  .bubble {
    animation: none;
  }
}
</style>
