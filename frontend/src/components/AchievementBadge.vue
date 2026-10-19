<script setup lang="ts">
defineProps<{ icon: string; unlocked: boolean; size?: number }>()

const PATHS: Record<string, string> = {
  keyboard: 'M2 8a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zM6 10h.01M10 10h.01M14 10h.01M18 10h.01M7 14h10',
  bolt: 'M13 2 4 14h7l-1 8 10-13h-7z',
  trophy: 'M8 21h8M12 17v4M7 4h10v5a5 5 0 0 1-10 0zM7 6H4v1a3 3 0 0 0 3 3M17 6h3v1a3 3 0 0 1-3 3',
  shield: 'M12 3 4 6v6c0 5 3.5 8 8 9 4.5-1 8-4 8-9V6zM9 12l2 2 4-4',
  target: 'M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18zM12 16a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM12 12h.01',
  clock: 'M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18zM12 7v5l3 2',
  flame: 'M12 2c1 4 5 6 5 11a5 5 0 0 1-10 0c0-2 1-3 2-4 0 2 1 3 2 3 0-4-1-6 1-10z',
  code: 'M8 7 3 12l5 5M16 7l5 5-5 5M14 4l-4 16',
  globe: 'M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18zM3 12h18M12 3c3 3 3 15 0 18M12 3c-3 3-3 15 0 18',
}
const STAR = 'M12 17.27 18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z'
</script>

<template>
  <span class="badge" :class="{ off: !unlocked }" :style="{ '--s': (size ?? 56) + 'px' }">
    <svg viewBox="0 0 24 24" width="55%" height="55%" aria-hidden="true">
      <path
        :d="icon === 'star' ? STAR : (PATHS[icon] ?? PATHS.target)"
        :fill="icon === 'star' ? 'currentColor' : 'none'"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </span>
</template>

<style scoped>
.badge {
  flex: none;
  width: var(--s);
  height: var(--s);
  border-radius: 50%;
  display: grid;
  place-items: center;
  color: var(--accent-ink);
  background: var(--accent);
  box-shadow: 0 4px 0 color-mix(in srgb, var(--accent) 58%, #000);
}
.badge.off {
  background: var(--surface-2);
  color: var(--muted);
  box-shadow: 0 4px 0 var(--border);
}
</style>
