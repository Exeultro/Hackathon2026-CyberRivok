let ctx: AudioContext | null = null
let master: GainNode | null = null

function audio() {
  ctx ??= new AudioContext()
  if (ctx.state === 'suspended') void ctx.resume()
  if (!master) {
    master = ctx.createGain()
    master.gain.value = 0.9
    master.connect(ctx.destination)
  }
  return { ctx, out: master }
}

interface ToneOptions {
  type?: OscillatorType
  vol?: number
  slideTo?: number
  lowpass?: number
}

function tone(freq: number, delay: number, dur: number, o: ToneOptions = {}) {
  try {
    const { ctx, out } = audio()
    const t = ctx.currentTime + delay
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.type = o.type ?? 'sine'
    osc.frequency.setValueAtTime(freq, t)
    if (o.slideTo) osc.frequency.exponentialRampToValueAtTime(o.slideTo, t + dur)
    gain.gain.setValueAtTime(0.0001, t)
    gain.gain.exponentialRampToValueAtTime(o.vol ?? 0.08, t + 0.008)
    gain.gain.exponentialRampToValueAtTime(0.0001, t + dur)
    let node: AudioNode = osc
    if (o.lowpass) {
      const f = ctx.createBiquadFilter()
      f.type = 'lowpass'
      f.frequency.value = o.lowpass
      node.connect(f)
      node = f
    }
    node.connect(gain).connect(out)
    osc.start(t)
    osc.stop(t + dur + 0.02)
  } catch {
    // audio is unavailable, ignore
  }
}

function click(vol: number) {
  try {
    const { ctx, out } = audio()
    const len = Math.floor(ctx.sampleRate * 0.018)
    const buf = ctx.createBuffer(1, len, ctx.sampleRate)
    const data = buf.getChannelData(0)
    for (let i = 0; i < len; i++) data[i] = (Math.random() * 2 - 1) * (1 - i / len) ** 2
    const src = ctx.createBufferSource()
    src.buffer = buf
    const band = ctx.createBiquadFilter()
    band.type = 'bandpass'
    band.frequency.value = 2400 + Math.random() * 900
    band.Q.value = 0.8
    const gain = ctx.createGain()
    gain.gain.value = vol
    src.connect(band).connect(gain).connect(out)
    src.start()
  } catch {
    // audio is unavailable, ignore
  }
}

const COMBO_NOTES: Record<number, number[]> = {
  2: [523.25, 659.25, 783.99],
  3: [587.33, 740.37, 880, 1174.66],
  5: [659.25, 830.61, 987.77, 1318.51, 1567.98],
}

export const sfx = {
  key() {
    click(0.16)
    tone(150 + Math.random() * 20, 0, 0.04, { vol: 0.05 })
  },
  back() {
    click(0.1)
    tone(110, 0, 0.05, { vol: 0.04 })
  },
  error() {
    tone(210, 0, 0.17, { type: 'sawtooth', vol: 0.06, slideTo: 120, lowpass: 900 })
    tone(105, 0, 0.17, { type: 'sine', vol: 0.07, slideTo: 70 })
  },
  combo(mult: number) {
    const notes = COMBO_NOTES[mult] ?? COMBO_NOTES[2]
    notes.forEach((f, i) => {
      tone(f, i * 0.055, 0.35, { type: 'triangle', vol: 0.06 })
      tone(f * 2, i * 0.055, 0.25, { type: 'sine', vol: 0.02 })
    })
  },
  star(i: number) {
    const f = 784 * 2 ** (i / 6)
    tone(f, 0, 0.55, { type: 'sine', vol: 0.07 })
    tone(f * 2.76, 0, 0.3, { type: 'sine', vol: 0.025 })
  },
  win() {
    ;[523.25, 659.25, 783.99, 1046.5].forEach((f, i) => tone(f, i * 0.09, 0.5, { type: 'triangle', vol: 0.07 }))
  },
  fail() {
    tone(392, 0, 0.35, { type: 'triangle', vol: 0.06 })
    tone(311.13, 0.22, 0.5, { type: 'triangle', vol: 0.06 })
  },
}
