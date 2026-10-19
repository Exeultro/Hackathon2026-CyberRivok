<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import type { CritterKind, PropKind, PropSpec, Scene } from '@/scenes'
import { TOWER, layout, nodePoint, trailSamples } from '@/geometry'
import type { Pt } from '@/geometry'
import Prop from './Prop.vue'
import Critter from './Critter.vue'
import Wall from './Wall.vue'
import { useViewBand } from '@/useViewBand'

const props = defineProps<{ scene: Scene; height: number; seed: string; nodeCount: number }>()

interface Placed {
  id: number
  kind: PropKind
  x: number
  y: number
  size: number
  tone: number
  depth: number
  variant: number
  sy: number
  sway: boolean
}

interface CritterPlaced {
  id: number
  kind: CritterKind
  x: number
  y: number
  size: number
  dx: number
  dy: number
  dur: number
  delay: number
  face: 1 | -1
  flying: boolean
}

const DECOR: PropKind[] = ['tuft', 'flowers', 'pond']
const GLOWING: PropKind[] = ['torch', 'campfire']
const ASPECT: Record<PropKind, number> = {
  oak: 64 / 60,
  pine: 80 / 60,
  bush: 36 / 50,
  flowers: 24 / 40,
  tuft: 16 / 24,
  mushroom: 1,
  stump: 26 / 30,
  log: 22 / 50,
  pond: 40 / 70,
  fence: 24 / 60,
  signpost: 44 / 30,
  hill: 70 / 140,
  haystack: 44 / 50,
  rock: 32 / 40,
  cairn: 28 / 24,
  peak: 1,
  torch: 2,
  banner: 56 / 30,
  tent: 50 / 60,
  campfire: 28 / 30,
  barrel: 26 / 20,
  catapult: 50 / 70,
}
const baseOf = (p: Placed) => p.y + 0.15 * p.size * ASPECT[p.kind]

const BIG: PropKind[] = ['peak', 'hill']
const STRUCT: PropKind[] = ['pond', 'tent', 'catapult']
const SWAYING: PropKind[] = ['oak', 'pine', 'bush', 'tuft', 'flowers']
const MAX_SWAYING = 50
const FOOT: Partial<Record<PropKind, { r: number; lift: number }>> = {
  peak: { r: 0.42, lift: 0.35 },
  hill: { r: 0.45, lift: 0.12 },
  pond: { r: 0.45, lift: 0.1 },
  tent: { r: 0.45, lift: 0.2 },
  catapult: { r: 0.45, lift: 0.15 },
}

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

function rng(seed: string) {
  let h = 1779033703 ^ seed.length
  for (let i = 0; i < seed.length; i++) {
    h = Math.imul(h ^ seed.charCodeAt(i), 3432918353)
    h = (h << 13) | (h >>> 19)
  }
  let a = h >>> 0
  return () => {
    a = (a + 0x6d2b79f5) | 0
    let t = Math.imul(a ^ (a >>> 15), 1 | a)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const trail = computed(() => trailSamples(props.nodeCount, width.value, props.scene.path))
const nodeCenters = computed(() => {
  const { colLeft, s } = layout(width.value)
  return Array.from({ length: props.nodeCount }, (_, i) => {
    const p = nodePoint(i, props.scene.path)
    return { x: colLeft + p.x * s, y: p.y }
  })
})
const edgeWidth = computed(() => {
  const { amp, amp2 } = props.scene.path
  return Math.max(60, width.value / 2 - ((amp + amp2) * layout(width.value).s + 85))
})

const wallGeom = computed(() => {
  const W = width.value
  const amp = Math.min(46, edgeWidth.value * 0.35)
  const base = amp + Math.max(30, Math.min(70, edgeWidth.value * 0.5))
  const arg = (side: -1 | 1, y: number) => y / 260 + (side < 0 ? 0.6 : 2.3)
  return {
    amp,
    x: (side: -1 | 1, y: number) => (side < 0 ? base : W - base) + amp * Math.sin(arg(side, y)),
  }
})

const wallLines = computed<Pt[][]>(() => {
  if (props.scene.terrain !== 'fortress') return []
  const wall = wallGeom.value
  return ([-1, 1] as const).map((side) => {
    const pts: Pt[] = []
    for (let y = -40; y <= props.height + 40; y += 12) pts.push({ x: wall.x(side, y), y })
    return pts
  })
})

const towerSpots = computed<Pt[]>(() => {
  if (props.scene.terrain !== 'fortress') return []
  const wall = wallGeom.value
  const out: Pt[] = []
  for (const side of [-1, 1] as const) {
    for (let y = side < 0 ? 180 : 420; y < props.height; y += 460) out.push({ x: wall.x(side, y) - side * TOWER.gap, y })
  }
  return out
})

const placed = computed<Placed[]>(() => {
  const W = width.value
  const rand = rng(props.seed)
  const edgeW = edgeWidth.value
  const fortress = props.scene.terrain === 'fortress'
  const wall = wallGeom.value
  const out: Placed[] = []
  const solids: { x: number; y: number; r: number }[] = towerSpots.value.map((s) => ({ x: s.x, y: s.y - 15, r: 48 }))

  const rank = (k: PropKind) => (BIG.includes(k) ? 0 : STRUCT.includes(k) ? 1 : 2)
  const specs = [...props.scene.props].sort((a, b) => rank(a.kind) - rank(b.kind))

  const clear = (kind: PropKind, x: number, y: number, size: number) => {
    const need = size * 0.45 + 40
    if (!trail.value.every((p) => Math.hypot(p.x - x, p.y - y) >= need)) return false
    if (!nodeCenters.value.every((p) => Math.hypot(p.x - x, p.y - y) >= need + 50)) return false
    if (fortress) {
      const pad = 46 + size * 0.5
      if (x < wall.x(-1, y) + pad || x > wall.x(1, y) - pad) return false
    }
    if (!BIG.includes(kind)) {
      const f = FOOT[kind] ?? { r: 0.3, lift: 0.15 }
      const cx = x
      const cy = y - size * f.lift
      const r = size * f.r
      if (!solids.every((s) => Math.hypot(s.x - cx, s.y - cy) >= (s.r + r) * 0.92)) return false
    }
    return true
  }

  const add = (spec: PropSpec, side: -1 | 1 | 0) => {
    for (let attempt = 0; attempt < 14; attempt++) {
      let x: number
      let depth: number
      if (side === 0) {
        const f = Math.pow(rand(), 3)
        x = rand() < 0.5 ? f * (W / 2) : W - f * (W / 2)
        depth = 1 - f
      } else {
        const f = Math.pow(rand(), 1.8)
        x = side < 0 ? f * edgeW : W - f * edgeW
        depth = 1 - f
      }
      const y = rand() * props.height
      const size = (spec.size[0] + rand() * (spec.size[1] - spec.size[0])) * (1 + depth * spec.grow)
      if (!clear(spec.kind, x, y, size)) continue
      const tone = rand()
      out.push({
        id: out.length,
        kind: spec.kind,
        x,
        y,
        size,
        tone,
        depth,
        variant: spec.variants ? Math.floor(rand() * spec.variants) : 0,
        sy: 1 + (rand() * 2 - 1) * (spec.jitter ?? 0),
        sway: false,
      })
      const f = FOOT[spec.kind]
      if (f && (BIG.includes(spec.kind) || STRUCT.includes(spec.kind))) {
        solids.push({ x, y: y - size * f.lift, r: size * f.r })
      }
      return
    }
  }

  for (const spec of specs) {
    if (spec.mode === 'scatter') {
      for (let i = 0; i < spec.count; i++) add(spec, 0)
    } else {
      for (const side of [-1, 1] as const) for (let i = 0; i < spec.count; i++) add(spec, side)
    }
  }

  for (const spot of towerSpots.value) {
    const side: -1 | 1 = spot.x < W / 2 ? -1 : 1
    const gx = spot.x
    const gy = spot.y
    out.push({
      id: out.length,
      kind: 'torch',
      x: gx - side * (TOWER.s / 2 + 26),
      y: gy + 8,
      size: 20,
      tone: 0.5,
      depth: 0,
      variant: 0,
      sy: 1,
      sway: false,
    })
  }

  const swayable = out.filter((p) => SWAYING.includes(p.kind))
  const every = Math.max(1, Math.ceil(swayable.length / MAX_SWAYING))
  swayable.forEach((p, i) => {
    p.sway = i % every === 0
  })

  return out.sort((a, b) => (DECOR.includes(a.kind) ? 0 : 1) - (DECOR.includes(b.kind) ? 0 : 1) || baseOf(a) - baseOf(b))
})

const SPEED: Partial<Record<CritterKind, number>> = {
  rabbit: 22,
  wolf: 26,
  goat: 14,
  soldier: 16,
  jumper: 12,
  bird: 48,
}
const FLYING: CritterKind[] = ['bird', 'butterfly', 'firefly', 'arrow']
const BLOCKING: PropKind[] = ['peak', 'hill', 'tent', 'catapult', 'pond']
const BLOCKING_FROG: PropKind[] = ['peak', 'hill']

const critters = computed<CritterPlaced[]>(() => {
  const rand = rng(props.seed + 'critters')
  const out: CritterPlaced[] = []
  const solids = placed.value.filter((p) => BLOCKING.includes(p.kind))

  const fortress = props.scene.terrain === 'fortress'
  const wall = wallGeom.value

  const inCorridor = (sx: number, sy: number) =>
    sx > wall.x(-1, sy) + 46 && sx < wall.x(1, sy) - 46

  const walkable = (kind: CritterKind, x: number, y: number, dx: number, dy: number, size: number) => {
    const kinds = kind === 'frog' ? BLOCKING_FROG : BLOCKING
    for (let k = 0; k <= 5; k++) {
      const sx = x + (dx * k) / 5
      const sy = y + (dy * k) / 5
      if (fortress && !inCorridor(sx, sy)) return false
      if (towerSpots.value.some((s) => Math.hypot(sx - s.x, sy - (s.y - 15)) < 44 + size * 0.4)) return false
      if (nodeCenters.value.some((n) => Math.hypot(sx - n.x, sy - n.y) < 58 + size * 0.3)) return false
      for (const b of solids) {
        if (!kinds.includes(b.kind)) continue
        if (Math.hypot(sx - b.x, sy - (b.y - b.size * 0.3)) < b.size * 0.52 + size * 0.4) return false
      }
    }
    return true
  }

  for (const spec of props.scene.critters) {
    const flying = FLYING.includes(spec.kind)
    for (let i = 0; i < spec.count; i++) {
      if (spec.kind === 'arrow' && fortress) {
        const side: -1 | 1 = rand() < 0.5 ? -1 : 1
        const y = 60 + rand() * (props.height - 120)
        const gap = wall.x(1, y) - wall.x(-1, y) - 110
        const dur = 2.8 + rand() * 2
        out.push({
          id: out.length,
          kind: 'arrow',
          x: wall.x(side, y) - side * 16,
          y,
          size: spec.size[0] + rand() * (spec.size[1] - spec.size[0]),
          dx: Math.max(60, Math.min(110 + rand() * 60, gap)),
          dy: 0,
          dur,
          delay: -rand() * dur,
          face: (-side) as 1 | -1,
          flying: true,
        })
        continue
      }
      for (let attempt = 0; attempt < 10; attempt++) {
        const f = Math.pow(rand(), flying ? 2 : 1.6)
        const span = flying ? edgeWidth.value : width.value / 2
        const x = rand() < 0.5 ? f * span : width.value - f * span
        const y = rand() * props.height
        if (!flying && !trail.value.every((p) => Math.hypot(p.x - x, p.y - y) >= 34)) continue
        const sign: 1 | -1 = rand() < 0.5 ? -1 : 1
        const size = spec.size[0] + rand() * (spec.size[1] - spec.size[0])
        let dx: number
        let dy: number
        let dur: number
        if (spec.kind === 'bird') {
          dx = sign * (60 + rand() * 80)
          dy = (rand() - 0.5) * 90
        } else if (spec.kind === 'arrow') {
          dx = 140 + rand() * 140
          dy = 0
        } else if (spec.kind === 'butterfly') {
          dx = sign * (40 + rand() * 70)
          dy = 30 + rand() * 40
        } else if (spec.kind === 'firefly') {
          dx = sign * (20 + rand() * 30)
          dy = 14 + rand() * 24
        } else if (spec.kind === 'frog' || spec.kind === 'archer') {
          dx = 0
          dy = 0
        } else {
          dx = sign * (50 + rand() * 90)
          dy = (rand() - 0.5) * 30
        }
        if (!flying && !walkable(spec.kind, x, y, dx, dy, size)) continue
        const speed = SPEED[spec.kind]
        if (spec.kind === 'arrow') dur = 2.8 + rand() * 2
        else if (speed) dur = (2 * Math.abs(dx)) / speed
        else if (spec.kind === 'frog' || spec.kind === 'archer') dur = 6 + rand() * 4
        else dur = 8 + rand() * 6
        out.push({
          id: out.length,
          kind: spec.kind,
          x,
          y,
          size,
          dx,
          dy,
          dur,
          delay: -rand() * dur,
          face: dx !== 0 ? (Math.sign(dx) as 1 | -1) : sign,
          flying,
        })
        break
      }
    }
  }
  return out
})

type Item = { prop: Placed; critter?: undefined } | { critter: CritterPlaced; prop?: undefined }

const visible = (y: number) => y > band.value.top - 220 && y < band.value.bottom + 220

const items = computed<Item[]>(() => {
  const list: { rank: number; y: number; item: Item }[] = []
  for (const p of placed.value) {
    if (visible(p.y)) list.push({ rank: DECOR.includes(p.kind) ? 0 : 1, y: baseOf(p), item: { prop: p } })
  }
  for (const c of critters.value) {
    if (visible(c.y)) list.push({ rank: c.flying ? 2 : 1, y: c.y, item: { critter: c } })
  }
  return list.sort((a, b) => a.rank - b.rank || a.y - b.y).map((e) => e.item)
})

const patches = computed(() => {
  const rand = rng(props.seed + 'ground')
  return Array.from({ length: 60 }, () => {
    const w = 140 + rand() * 220
    return {
      left: rand() * width.value,
      top: rand() * props.height,
      w,
      h: w * (0.4 + rand() * 0.3),
      color: rand() < 0.5 ? props.scene.ground.patchA : props.scene.ground.patchB,
    }
  })
})

const glows = computed(() => placed.value.filter((p) => GLOWING.includes(p.kind) && visible(p.y)))
</script>

<template>
  <div ref="root" class="terrain" aria-hidden="true">
    <div class="ground">
      <i
        v-for="(p, i) in patches"
        :key="i"
        class="patch"
        :style="{
          left: p.left + 'px',
          top: p.top + 'px',
          width: p.w + 'px',
          height: p.h + 'px',
          background: `radial-gradient(closest-side, ${p.color}, transparent)`,
        }"
      />
    </div>
    <Wall v-if="wallLines.length" :lines="wallLines" :towers="towerSpots" :width="width" :height="height" />
    <template v-for="it in items" :key="it.prop ? 'p' + it.prop.id : 'c' + it.critter.id">
      <Prop
        v-if="it.prop"
        :kind="it.prop.kind"
        :variant="it.prop.variant"
        :sway="it.prop.sway"
        :style="{
          left: it.prop.x + 'px',
          top: it.prop.y + 'px',
          width: it.prop.size + 'px',
          '--tone': it.prop.tone,
          '--depth': it.prop.depth,
          '--sy': it.prop.sy,
        }"
      />
      <Critter
        v-else-if="it.critter"
        :kind="it.critter.kind"
        :style="{
          left: it.critter.x + 'px',
          top: it.critter.y + 'px',
          width: it.critter.size + 'px',
          '--dx': it.critter.dx + 'px',
          '--dy': it.critter.dy + 'px',
          '--dur': it.critter.dur + 's',
          '--delay': it.critter.delay + 's',
          '--face': it.critter.face,
        }"
      />
    </template>
    <div class="layer" :style="{ background: scene.tint }" />
    <div class="layer" :style="{ background: scene.wash }" />
    <div
      class="layer"
      :style="{
        background: `linear-gradient(to right, rgb(0 0 0 / ${scene.vignette}), transparent 22%, transparent 78%, rgb(0 0 0 / ${scene.vignette}))`,
      }"
    />
    <i
      v-for="(g, i) in glows"
      :key="'g' + i"
      class="glow"
      :style="{ left: g.x + 'px', top: g.y - g.size * (g.kind === 'torch' ? 1.2 : 0.5) + 'px' }"
    />
  </div>
</template>

<style scoped>
.terrain {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  background: var(--ground);
  isolation: isolate;
}
.ground,
.layer {
  position: absolute;
  inset: 0;
}
.patch {
  position: absolute;
  transform: translate(-50%, -50%);
  opacity: 0.6;
}
.glow {
  position: absolute;
  width: 170px;
  height: 170px;
  transform: translate(-50%, -50%);
  background: radial-gradient(closest-side, rgb(255 170 70 / 0.5), transparent);
  mix-blend-mode: screen;
}
</style>
