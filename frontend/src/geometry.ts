export interface Pt {
  x: number
  y: number
}

export const PATH_W = 600
export const STEP = 190
export const TOP = 130

export interface PathShape {
  amp: number
  freq: number
  phase: number
  amp2: number
  freq2: number
  phase2: number
}

export function nodePoint(i: number, shape: PathShape): Pt {
  const x =
    PATH_W / 2 +
    shape.amp * Math.sin(i * shape.freq + shape.phase) +
    shape.amp2 * Math.sin(i * shape.freq2 + shape.phase2)
  return { x, y: TOP + i * STEP }
}

export function pathHeight(count: number) {
  return TOP + count * STEP - 40
}

export function layout(sceneW: number) {
  const colW = Math.max(280, Math.min(640, sceneW) - 32)
  return { colW, colLeft: (sceneW - colW) / 2, s: colW / PATH_W }
}

export function controls(a: Pt, b: Pt): [Pt, Pt] {
  const mid = (a.y + b.y) / 2
  return [
    { x: a.x, y: mid },
    { x: b.x, y: mid },
  ]
}

export function cubic(a: Pt, c1: Pt, c2: Pt, b: Pt, t: number): Pt {
  const u = 1 - t
  return {
    x: u * u * u * a.x + 3 * u * u * t * c1.x + 3 * u * t * t * c2.x + t * t * t * b.x,
    y: u * u * u * a.y + 3 * u * u * t * c1.y + 3 * u * t * t * c2.y + t * t * t * b.y,
  }
}

export function trailSamples(count: number, sceneW: number, shape: PathShape): Pt[] {
  const { colLeft, s } = layout(sceneW)
  const pts = Array.from({ length: count }, (_, i) => {
    const p = nodePoint(i, shape)
    return { x: colLeft + p.x * s, y: p.y }
  })
  const out: Pt[] = pts.slice(0, 1)
  for (let i = 1; i < pts.length; i++) {
    const [c1, c2] = controls(pts[i - 1], pts[i])
    for (let k = 1; k <= 14; k++) out.push(cubic(pts[i - 1], c1, c2, pts[i], k / 14))
  }
  return out
}

export const TOWER = { s: 58, rise: { x: -14, y: -37 }, gap: 56 }
