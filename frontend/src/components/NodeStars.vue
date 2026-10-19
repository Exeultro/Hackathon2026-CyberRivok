<script setup lang="ts">
import StarIcon from './StarIcon.vue'

defineProps<{ count: number }>()

const RADIUS = 66
const SLOTS = [
  { angle: -42, size: 27 },
  { angle: 0, size: 35 },
  { angle: 42, size: 27 },
]

function place(angle: number) {
  const rad = (angle * Math.PI) / 180
  return { x: RADIUS * Math.sin(rad), y: -RADIUS * Math.cos(rad) }
}
</script>

<template>
  <div class="arc" :aria-label="`Звёзд: ${count} из 3`" role="img">
    <span
      v-for="(s, i) in SLOTS"
      :key="i"
      class="star"
      :style="{
        transform: `translate(calc(-50% + ${place(s.angle).x}px), calc(-50% + ${place(s.angle).y}px)) rotate(${s.angle / 2}deg)`,
      }"
    >
      <StarIcon :filled="i < count" :size="s.size" />
    </span>
  </div>
</template>

<style scoped>
.arc {
  position: absolute;
  left: 50%;
  top: 50%;
  width: 0;
  height: 0;
  pointer-events: none;
}
.star {
  position: absolute;
  left: 0;
  top: 0;
  display: block;
  line-height: 0;
  filter: drop-shadow(0 2px 0 rgb(0 0 0 / 0.2));
}
</style>
