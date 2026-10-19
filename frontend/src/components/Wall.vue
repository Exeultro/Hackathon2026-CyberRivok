<script setup lang="ts">
import { computed } from 'vue'
import { TOWER } from '@/geometry'
import type { Pt } from '@/geometry'

const props = defineProps<{ lines: Pt[][]; towers: Pt[]; width: number; height: number }>()

const MERLON_OFFSET = 15

function path(pts: Pt[]) {
  return pts.map((p, i) => `${i ? 'L' : 'M'}${p.x.toFixed(1)},${p.y.toFixed(1)}`).join(' ')
}

function offset(pts: Pt[], d: number): Pt[] {
  return pts.map((p, i) => {
    const a = pts[Math.max(0, i - 1)]
    const b = pts[Math.min(pts.length - 1, i + 1)]
    const len = Math.hypot(b.x - a.x, b.y - a.y) || 1
    return { x: p.x - ((b.y - a.y) / len) * d, y: p.y + ((b.x - a.x) / len) * d }
  })
}

const H = TOWER.s / 2
const SLICES = 14

const bastions = computed(() =>
  props.towers.map((p) => {
    const gx = p.x
    const gy = p.y
    const tx = gx + TOWER.rise.x
    const ty = gy + TOWER.rise.y
    const at = (k: number) => ({ x: gx + ((tx - gx) * k) / SLICES, y: gy + ((ty - gy) * k) / SLICES })
    const slices = Array.from({ length: SLICES + 1 }, (_, k) => {
      const c = at(k)
      const size = TOWER.s + (k === 0 ? 10 : k === 1 ? 5 : 0)
      return { x: c.x - size / 2, y: c.y - size / 2, size, t: k / SLICES }
    })
    const mid = at(Math.round(SLICES * 0.55))
    return {
      gx,
      gy,
      tx,
      ty,
      slices,
      shadow: { cx: gx + 12, cy: gy + H - 2, rx: TOWER.s * 0.82, ry: TOWER.s * 0.3 },
      east: `${tx + H},${ty - H} ${tx + H},${ty + H} ${gx + H},${gy + H} ${gx + H},${gy - H}`,
      door: { x: gx - 9, y: gy + H - 21, w: 18, h: 21 },
      slits: [mid.x - 15, mid.x + 12].map((x) => ({ x, y: mid.y + H - 12 })),
    }
  }),
)

const walls = computed(() =>
  props.lines.map((pts) => ({
    d: path(pts),
    left: path(offset(pts, -MERLON_OFFSET)),
    right: path(offset(pts, MERLON_OFFSET)),
  })),
)
</script>

<template>
  <svg class="wall" :width="width" :height="height" :viewBox="`0 0 ${width} ${height}`" aria-hidden="true">
    <g v-for="(w, i) in walls" :key="i">
      <path :d="w.d" class="shadow" transform="translate(16 18)" />
      <path :d="w.d" class="face" transform="translate(10 11)" />
      <path :d="w.d" class="bricks" transform="translate(10 11)" />
      <path :d="w.d" class="top" />
      <path :d="w.d" class="seams" />
      <path :d="w.d" class="walk" />
      <path :d="w.left" class="merlon-base" />
      <path :d="w.left" class="merlon" />
      <path :d="w.right" class="merlon-base" />
      <path :d="w.right" class="merlon" />
    </g>
    <g v-for="(b, i) in bastions" :key="'t' + i">
      <ellipse :cx="b.shadow.cx" :cy="b.shadow.cy" :rx="b.shadow.rx" :ry="b.shadow.ry" class="tshadow" />
      <rect
        v-for="(s, k) in b.slices"
        :key="k"
        :x="s.x"
        :y="s.y"
        :width="s.size"
        :height="s.size"
        rx="7"
        class="tslice"
        :style="{ '--t': s.t }"
      />
      <polygon :points="b.east" class="teast" />
      <path :d="`M${b.door.x},${b.door.y + b.door.h} V${b.door.y + 7} a9,7 0 0 1 18,0 V${b.door.y + b.door.h} Z`" class="tdoor" />
      <rect v-for="(sl, k) in b.slits" :key="'s' + k" :x="sl.x" :y="sl.y" width="3.4" height="9" rx="1.4" class="tdoor" />
      <rect :x="b.tx - H - 3" :y="b.ty - H - 3" :width="TOWER.s + 6" :height="TOWER.s + 6" rx="7" class="ttop" />
      <rect :x="b.tx - H + 7" :y="b.ty - H + 7" :width="TOWER.s - 14" :height="TOWER.s - 14" rx="3" class="tfloor" />
      <rect :x="b.tx - H + 1" :y="b.ty - H + 2" :width="TOWER.s - 2" :height="TOWER.s - 2" rx="5" class="tring-base" />
      <rect :x="b.tx - H" :y="b.ty - H" :width="TOWER.s" :height="TOWER.s" rx="5" class="tring" />
    </g>
  </svg>
</template>

<style scoped>
.wall {
  position: absolute;
  inset: 0;
}
path {
  fill: none;
  stroke-linejoin: round;
}
.shadow {
  stroke: rgb(0 0 0 / 0.28);
  stroke-width: 36;
}
.face {
  stroke: var(--w3);
  stroke-width: 34;
}
.bricks {
  stroke: rgb(0 0 0 / 0.35);
  stroke-width: 34;
  stroke-dasharray: 1.4 9;
}
.walk {
  stroke: var(--w2);
  stroke-width: 12;
  opacity: 0.45;
}
.top {
  stroke: var(--w1);
  stroke-width: 34;
}
.seams {
  stroke: var(--w3);
  stroke-width: 30;
  stroke-dasharray: 1.2 14;
  opacity: 0.35;
}
.tshadow {
  fill: rgb(0 0 0 / 0.3);
}
.tslice {
  fill: color-mix(in srgb, var(--w3) calc((1 - var(--t)) * 72%), var(--w2));
  stroke: rgb(0 0 0 / 0.2);
  stroke-width: 1;
}
.teast {
  fill: rgb(0 0 0 / 0.2);
}
.tdoor {
  fill: rgb(0 0 0 / 0.55);
}
.ttop {
  fill: var(--w1);
}
.tfloor {
  fill: var(--w3);
  opacity: 0.5;
}
.tring-base,
.tring {
  fill: none;
  stroke-width: 8;
  stroke-dasharray: 10 8;
}
.tring-base {
  stroke: var(--w3);
}
.tring {
  stroke: var(--w1);
}
.merlon-base {
  stroke: var(--w3);
  stroke-width: 10;
  stroke-dasharray: 12 12;
  transform: translate(3px, 4px);
}
.merlon {
  stroke: var(--w1);
  stroke-width: 9;
  stroke-dasharray: 12 12;
}
</style>
