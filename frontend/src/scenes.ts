import type { MapId } from '@/api'
import type { PathShape } from '@/geometry'

export type PropKind =
  | 'oak'
  | 'pine'
  | 'bush'
  | 'flowers'
  | 'tuft'
  | 'mushroom'
  | 'stump'
  | 'log'
  | 'pond'
  | 'fence'
  | 'signpost'
  | 'hill'
  | 'haystack'
  | 'rock'
  | 'cairn'
  | 'peak'
  | 'torch'
  | 'banner'
  | 'tent'
  | 'campfire'
  | 'barrel'
  | 'catapult'

export type CritterKind =
  | 'rabbit'
  | 'wolf'
  | 'frog'
  | 'butterfly'
  | 'firefly'
  | 'bird'
  | 'goat'
  | 'soldier'
  | 'jumper'
  | 'archer'
  | 'arrow'

export interface CritterSpec {
  kind: CritterKind
  count: number
  size: [number, number]
}

export interface PropSpec {
  kind: PropKind
  count: number
  size: [number, number]
  grow: number
  mode?: 'edge' | 'scatter'
  variants?: number
  jitter?: number
}

export type WindKind = 'leaves' | 'petals' | 'dust' | 'embers'

export interface Scene {
  label: string
  terrain: 'forest' | 'hills' | 'mountains' | 'fortress'
  palette: Record<string, string>
  ground: { base: string; patchA: string; patchB: string }
  light: { shx: number; shy: number; sho: number; flip: 1 | -1 }
  tint: string
  wash: string
  vignette: number
  path: PathShape
  wind: { kind: WindKind; count: number }
  props: PropSpec[]
  critters: CritterSpec[]
}

const BASE_PALETTE = {
  f1: '#9be26b',
  f2: '#55b24c',
  f3: '#2f8244',
  r1: '#d3d4cb',
  r2: '#a5a79f',
  r3: '#777a73',
  s1: '#ffffff',
  s2: '#c9d6e8',
  w1: '#a9aebd',
  w2: '#7d849a',
  w3: '#4f556a',
  y1: '#f3d96e',
  y2: '#cfa72a',
  trunk: '#7a4f2a',
  edge: '#1d4a2c',
  trail: '#e8d5a0',
}

export const SCENES: Record<MapId, Scene> = {
  journey: {
    label: 'Утро',
    terrain: 'forest',
    palette: { ...BASE_PALETTE },
    ground: { base: '#86c85a', patchA: '#9bd96c', patchB: '#74b84c' },
    light: { shx: 12, shy: 4, sho: 0.25, flip: 1 },
    tint: 'rgb(255 200 120 / 0.07)',
    wash: 'linear-gradient(to right, rgb(255 224 170 / 0.28), transparent 55%)',
    vignette: 0.12,
    path: { amp: 150, freq: 1.15, phase: 0.4, amp2: 0, freq2: 1, phase2: 0 },
    wind: { kind: 'leaves', count: 48 },
    critters: [
      { kind: 'rabbit', count: 8, size: [34, 44] },
      { kind: 'wolf', count: 3, size: [60, 76] },
    ],
    props: [
      { kind: 'tuft', count: 50, size: [14, 22], grow: 0, mode: 'scatter' },
      { kind: 'flowers', count: 34, size: [22, 32], grow: 0, mode: 'scatter' },
      { kind: 'mushroom', count: 10, size: [16, 22], grow: 0, mode: 'scatter' },
      { kind: 'stump', count: 6, size: [26, 34], grow: 0, mode: 'scatter' },
      { kind: 'log', count: 5, size: [44, 56], grow: 0, mode: 'scatter' },
      { kind: 'rock', count: 8, size: [20, 30], grow: 0, mode: 'scatter' },
      { kind: 'bush', count: 16, size: [30, 44], grow: 0.3, mode: 'scatter' },
      { kind: 'pond', count: 2, size: [90, 120], grow: 0, mode: 'scatter' },
      { kind: 'signpost', count: 2, size: [26, 30], grow: 0, mode: 'scatter' },
      { kind: 'pine', count: 45, size: [40, 64], grow: 0.8 },
      { kind: 'oak', count: 50, size: [50, 78], grow: 0.9 },
    ],
  },
  chase: {
    label: 'День',
    terrain: 'hills',
    palette: {
      ...BASE_PALETTE,
      f1: '#d8e874',
      f2: '#a0c94e',
      f3: '#6b9f3a',
      edge: '#356b2c',
    },
    ground: { base: '#9bd75f', patchA: '#aee470', patchB: '#86c24f' },
    light: { shx: 3, shy: 8, sho: 0.22, flip: 1 },
    tint: 'transparent',
    wash: 'transparent',
    vignette: 0.1,
    path: { amp: 195, freq: 2.0, phase: 0.2, amp2: 0, freq2: 1, phase2: 0 },
    wind: { kind: 'petals', count: 52 },
    critters: [
      { kind: 'frog', count: 9, size: [24, 32] },
      { kind: 'butterfly', count: 9, size: [22, 30] },
      { kind: 'firefly', count: 16, size: [9, 12] },
    ],
    props: [
      { kind: 'tuft', count: 60, size: [14, 22], grow: 0, mode: 'scatter' },
      { kind: 'flowers', count: 44, size: [22, 32], grow: 0, mode: 'scatter' },
      { kind: 'rock', count: 12, size: [20, 32], grow: 0, mode: 'scatter' },
      { kind: 'bush', count: 14, size: [30, 44], grow: 0.3, mode: 'scatter' },
      { kind: 'haystack', count: 6, size: [38, 50], grow: 0, mode: 'scatter' },
      { kind: 'fence', count: 6, size: [56, 70], grow: 0, mode: 'scatter' },
      { kind: 'signpost', count: 2, size: [26, 30], grow: 0, mode: 'scatter' },
      { kind: 'pond', count: 3, size: [100, 140], grow: 0, mode: 'scatter' },
      { kind: 'oak', count: 8, size: [44, 60], grow: 0.4, mode: 'scatter' },
      { kind: 'hill', count: 26, size: [100, 170], grow: 0.8, variants: 3, jitter: 0.25 },
    ],
  },
  search: {
    label: 'Вечер',
    terrain: 'mountains',
    palette: {
      ...BASE_PALETTE,
      f1: '#5aa07c',
      f2: '#3a7561',
      f3: '#24493f',
      r1: '#e0d7c6',
      r2: '#b3a893',
      r3: '#7c725f',
      s1: '#fff4ea',
      s2: '#c9b6dc',
      edge: '#3a2a3d',
      trail: '#efe2c4',
    },
    ground: { base: '#b3a98a', patchA: '#c4bb9d', patchB: '#9f9575' },
    light: { shx: -16, shy: 5, sho: 0.3, flip: -1 },
    tint: 'rgb(255 140 70 / 0.12)',
    wash: 'linear-gradient(to left, rgb(255 140 80 / 0.25), transparent 55%)',
    vignette: 0.15,
    path: { amp: 120, freq: 0.9, phase: 1.0, amp2: 75, freq2: 2.6, phase2: 0.5 },
    wind: { kind: 'dust', count: 42 },
    critters: [
      { kind: 'bird', count: 8, size: [30, 40] },
      { kind: 'goat', count: 5, size: [46, 58] },
    ],
    props: [
      { kind: 'tuft', count: 22, size: [14, 20], grow: 0, mode: 'scatter' },
      { kind: 'rock', count: 22, size: [20, 36], grow: 0.2, mode: 'scatter' },
      { kind: 'cairn', count: 5, size: [22, 28], grow: 0, mode: 'scatter' },
      { kind: 'pine', count: 12, size: [32, 48], grow: 0.3, mode: 'scatter' },
      { kind: 'signpost', count: 2, size: [26, 30], grow: 0, mode: 'scatter' },
      { kind: 'rock', count: 22, size: [24, 44], grow: 0.8 },
      { kind: 'pine', count: 26, size: [34, 50], grow: 0.6 },
      { kind: 'peak', count: 30, size: [96, 160], grow: 0.9, variants: 2, jitter: 0.3 },
    ],
  },
  siege: {
    label: 'Ночь',
    terrain: 'fortress',
    palette: {
      ...BASE_PALETTE,
      edge: '#0b0f1e',
      r1: '#a8aab4',
      r2: '#7e808c',
      r3: '#575964',
      f1: '#5d9a6c',
      f2: '#3a7050',
      f3: '#244634',
      trail: '#9a99b2',
    },
    ground: { base: '#7a7468', patchA: '#8a8478', patchB: '#686258' },
    light: { shx: 0, shy: 6, sho: 0.22, flip: 1 },
    tint: 'rgb(22 34 98 / 0.5)',
    wash: 'transparent',
    vignette: 0.3,
    path: { amp: 165, freq: 1.55, phase: 2.2, amp2: 45, freq2: 0.6, phase2: 1.0 },
    wind: { kind: 'embers', count: 40 },
    critters: [
      { kind: 'soldier', count: 6, size: [24, 32] },
      { kind: 'jumper', count: 5, size: [24, 32] },
      { kind: 'archer', count: 6, size: [26, 34] },
      { kind: 'arrow', count: 9, size: [28, 40] },
    ],
    props: [
      { kind: 'tuft', count: 14, size: [14, 20], grow: 0, mode: 'scatter' },
      { kind: 'rock', count: 16, size: [18, 32], grow: 0, mode: 'scatter' },
      { kind: 'barrel', count: 12, size: [16, 22], grow: 0, mode: 'scatter' },
      { kind: 'banner', count: 6, size: [24, 32], grow: 0, mode: 'scatter' },
      { kind: 'tent', count: 8, size: [50, 66], grow: 0, mode: 'scatter' },
      { kind: 'campfire', count: 5, size: [24, 30], grow: 0, mode: 'scatter' },
      { kind: 'catapult', count: 3, size: [66, 88], grow: 0, mode: 'scatter' },
    ],
  },
}

export function sceneVars(s: Scene) {
  const vars: Record<string, string> = {
    '--ground': s.ground.base,
    '--shx': String(s.light.shx),
    '--shy': String(s.light.shy),
    '--sho': String(s.light.sho),
    '--flip': String(s.light.flip),
  }
  for (const [k, v] of Object.entries(s.palette)) vars[`--${k}`] = v
  return vars
}
