<script setup lang="ts">
import { computed } from 'vue'
import type { Lang } from '@/api'
import { FINGERS, FINGER_COLORS, LAYOUTS, findTarget } from '@/typing/layouts'

const props = defineProps<{ lang: Lang; expected: string | null; message?: string; tab?: boolean }>()

const rows = computed(() => LAYOUTS[props.lang])
const target = computed(() => findTarget(props.lang, props.expected))
const shiftSide = computed(() => {
  const t = target.value
  if (!t?.shift) return null
  return t.finger <= 3 ? 'right' : 'left'
})
const hint = computed(() => {
  if (props.message) return props.message
  if (props.tab) return 'Tab (4 пробела)'
  const t = target.value
  if (!t) return ''
  return `${FINGERS[t.finger]}${t.shift ? ' + Shift' : ''}`
})

function label(n: string, s: string) {
  return /\p{L}/u.test(n) ? s : n
}
function isTarget(r: number, c: number) {
  const t = target.value
  return !!t && t.row === r && t.col === c
}
</script>

<template>
  <div class="kb" aria-hidden="true">
    <div v-for="(r, ri) in rows" :key="ri" class="row" :style="{ '--off': [0, 1, 1.3, 0][ri] }">
      <span v-if="ri === 3" class="key wide shift" :class="{ target: shiftSide === 'left' }" :style="{ '--fc': FINGER_COLORS[0] }">Shift</span>
      <span
        v-for="(k, ci) in r"
        :key="ci"
        class="key"
        :class="{ target: isTarget(ri, ci) }"
        :style="{ '--fc': FINGER_COLORS[k.f] }"
      >
        {{ label(k.n, k.s) }}
      </span>
      <span v-if="ri === 2" class="key wide" :class="{ target: expected === '\n' }" :style="{ '--fc': FINGER_COLORS[7] }">↵</span>
      <span v-if="ri === 3" class="key wide shift" :class="{ target: shiftSide === 'right' }" :style="{ '--fc': FINGER_COLORS[7] }">Shift</span>
    </div>
    <div class="row">
      <span class="key space" :class="{ target: expected === ' ' }" :style="{ '--fc': FINGER_COLORS[8] }" />
    </div>
    <div class="finger" :class="{ fix: !!message }">{{ message ? hint : hint ? `Нажми: ${hint}` : '' }}</div>
  </div>
</template>

<style scoped>
.kb {
  --k: calc(clamp(26px, calc((100vh - 520px) / 7), 38px) * max(var(--fit, 1), 0.85));
  display: grid;
  gap: 4px;
  justify-items: center;
  padding: 8px 8px 0;
}
.row {
  display: flex;
  gap: 5px;
  padding-left: calc(var(--off, 0) * var(--k));
}
.key {
  min-width: var(--k);
  height: var(--k);
  display: grid;
  place-items: center;
  border-radius: 8px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  box-shadow: inset 0 -4px 0 var(--fc);
  font-size: 14px;
  font-weight: 600;
  color: var(--muted);
  transition: transform 0.08s;
}
.key.wide {
  padding: 0 12px;
  font-size: 12px;
}
.key.space {
  width: calc(var(--k) * 7);
}
.key.target {
  background: var(--fc);
  color: #1b1d24;
  border-color: var(--fc);
  transform: translateY(2px) scale(1.08);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--fc) 40%, transparent);
}
.finger.fix {
  color: var(--err);
  font-weight: 700;
}
.finger {
  min-height: 20px;
  font-size: 14px;
  color: var(--muted);
}
@media (max-width: 760px) {
  .kb {
    --k: calc(24px * max(var(--fit, 1), 0.85));
  }
  .key {
    font-size: 11px;
    border-radius: 5px;
  }
  .key.wide {
    padding: 0 6px;
    font-size: 9px;
  }
}
</style>
