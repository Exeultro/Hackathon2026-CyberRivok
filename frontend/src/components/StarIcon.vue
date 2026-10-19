<script setup lang="ts">
defineProps<{ filled: boolean; size?: number }>()
</script>

<template>
  <svg :width="size ?? 24" :height="size ?? 24" viewBox="0 0 24 24" aria-hidden="true">
    <path
      d="M12 17.27 18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"
      :fill="filled ? 'var(--accent)' : 'rgb(0 0 0 / 0.28)'"
    />
  </svg>
</template>
