<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useViewBand } from '@/useViewBand'
import type { WindKind } from '@/scenes'

const props = defineProps<{ kind: WindKind; count: number; height: number; seed: string }>()

const root = ref<HTMLElement | null>(null)
const width = ref(1200)
const band = useViewBand(root)
let observer: ResizeObserver | null = null

onMounted(() => {
  if (!root.value) return
  width.value = root.value.clientWidth || width.value
  observer = new ResizeObserver(([entry]) => {
    width.value = entry.contentRect.width
  })
  observer.observe(root.value)
})
onBeforeUnmount(() => observer?.disconnect())

const PALETTE: Record<WindKind, string[]> = {
  leaves: ['#9be26b', '#d9c342', '#e08a3a', '#6fb04a'],
  petals: ['#ffffff', '#f9a8d4', '#fde047', '#ffffff'],
  dust: ['#ffffff', '#efe2c4', '#d9c9a3'],
  embers: ['#ffb347', '#ff7a2f', '#ffd27a'],
}

function rng(seed: string) {
  let a = 2166136261
  for (let i = 0; i < seed.length; i++) a = Math.imul(a ^ seed.charCodeAt(i), 16777619)
  return () => {
    a = (a + 0x6d2b79f5) | 0
    let t = Math.imul(a ^ (a >>> 15), 1 | a)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const particles = computed(() => {
  const rand = rng(props.seed + 'wind' + props.kind)
  const colors = PALETTE[props.kind]
  const edge = Math.max(90, width.value * 0.3)
  return Array.from({ length: props.count }, () => {
    const f = Math.pow(rand(), 1.4)
    const left = rand() < 0.5 ? f * edge : width.value - f * edge
    const up = rand() < 0.5 ? -1 : 1
    let size: number
    let dy: number
    let dur: number
    let sway: number
    if (props.kind === 'leaves') {
      size = 9 + rand() * 6
      dy = 140 + rand() * 180
      dur = 9 + rand() * 7
      sway = 16 + rand() * 26
    } else if (props.kind === 'petals') {
      size = 5 + rand() * 4
      dy = up * (90 + rand() * 140)
      dur = 10 + rand() * 8
      sway = 20 + rand() * 28
    } else if (props.kind === 'dust') {
      size = 14 + rand() * 16
      dy = up * (50 + rand() * 90)
      dur = 6 + rand() * 5
      sway = 12 + rand() * 16
    } else {
      size = 3 + rand() * 3
      dy = -(160 + rand() * 150)
      dur = 5 + rand() * 4
      sway = 10 + rand() * 20
    }
    return {
      left,
      top: rand() * props.height,
      size,
      dy,
      dur,
      delay: -rand() * dur,
      sway,
      swayDur: 2.2 + rand() * 2.6,
      spin: 1.4 + rand() * 1.6,
      color: colors[Math.floor(rand() * colors.length)],
    }
  })
})

const shown = computed(() =>
  particles.value.filter((p) => p.top > band.value.top - 300 && p.top < band.value.bottom + 300),
)
</script>

<template>
  <div ref="root" class="wind" :class="kind" aria-hidden="true">
    <i
      v-for="(p, n) in shown"
      :key="n"
      class="p"
      :style="{
        left: p.left + 'px',
        top: p.top + 'px',
        '--dy': p.dy + 'px',
        animationDuration: p.dur + 's',
        animationDelay: p.delay + 's',
      }"
    >
      <u
        class="q"
        :style="{
          '--sw': p.sway + 'px',
          width: p.size + 'px',
          height: (kind === 'dust' ? 2 : p.size) + 'px',
          background: p.color,
          animationDuration: p.swayDur + 's, ' + p.spin + 's',
          animationDelay: -p.swayDur * 0.5 + 's, 0s',
        }"
      />
    </i>
  </div>
</template>

<style scoped>
.wind {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  z-index: 3;
}
.p {
  position: absolute;
  display: block;
  animation: travel linear infinite;
}
.q {
  display: block;
  text-decoration: none;
  animation-name: sway, spin;
  animation-timing-function: ease-in-out, linear;
  animation-iteration-count: infinite, infinite;
  animation-direction: alternate, normal;
}
.leaves .q {
  border-radius: 0 100% 0 100%;
}
.petals .q {
  border-radius: 50% 50% 50% 0;
  animation-name: sway, wobble;
  animation-direction: alternate, alternate;
}
.dust .q {
  border-radius: 2px;
  opacity: 0.55;
  animation-name: sway, slant;
  animation-direction: alternate, alternate;
}
.embers .q {
  border-radius: 50%;
  box-shadow: 0 0 6px 2px rgb(255 150 60 / 0.8);
  animation-name: sway, flicker;
  animation-direction: alternate, alternate;
}
@keyframes travel {
  0% {
    transform: translateY(0);
    opacity: 0;
  }
  14% {
    opacity: 1;
  }
  86% {
    opacity: 1;
  }
  100% {
    transform: translateY(var(--dy));
    opacity: 0;
  }
}
@keyframes sway {
  from {
    translate: calc(var(--sw) * -1) 0;
  }
  to {
    translate: var(--sw) 0;
  }
}
@keyframes spin {
  0% {
    rotate: 0deg;
    scale: 1 1;
  }
  50% {
    rotate: 180deg;
    scale: 0.35 1;
  }
  100% {
    rotate: 360deg;
    scale: 1 1;
  }
}
@keyframes wobble {
  from {
    rotate: -35deg;
  }
  to {
    rotate: 40deg;
  }
}
@keyframes slant {
  from {
    rotate: -25deg;
  }
  to {
    rotate: 25deg;
  }
}
@keyframes flicker {
  from {
    scale: 0.7;
    opacity: 0.6;
  }
  to {
    scale: 1.15;
    opacity: 1;
  }
}
@media (prefers-reduced-motion: reduce) {
  .p,
  .q {
    animation: none;
  }
}
</style>
