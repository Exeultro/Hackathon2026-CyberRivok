<script setup lang="ts">
const items = [
  { to: '/map', label: 'Карта', icon: 'M9 4 3 6v14l6-2 6 2 6-2V4l-6 2-6-2zM9 4v14M15 6v14' },
  {
    to: '/practice',
    label: 'Тренировка',
    icon: 'M2 8a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zM6 10h.01M10 10h.01M14 10h.01M18 10h.01M7 14h10',
  },
  { to: '/stats', label: 'Статистика', icon: 'M4 20V10M10 20V4M16 20v-7M22 20H2' },
  {
    to: '/leaderboard',
    label: 'Лидеры',
    icon: 'M8 21h8M12 17v4M7 4h10v5a5 5 0 0 1-10 0zM7 6H4v1a3 3 0 0 0 3 3M17 6h3v1a3 3 0 0 1-3 3',
  },
  { to: '/profile', label: 'Профиль', icon: 'M20 21a8 8 0 0 0-16 0M12 13a4 4 0 1 0 0-8 4 4 0 0 0 0 8z' },
  { to: '/settings', label: 'Настройки', icon: 'M4 6h10M18 6h2M4 12h2M10 12h10M4 18h12M14 4v4M8 10v4M18 16v4' },
]
</script>

<template>
  <nav class="nav" aria-label="Разделы">
    <RouterLink v-for="i in items" :key="i.to" :to="i.to" class="item">
      <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path :d="i.icon" />
      </svg>
      <span>{{ i.label }}</span>
    </RouterLink>
  </nav>
</template>

<style scoped>
.nav {
  position: fixed;
  top: var(--head-h);
  left: 0;
  bottom: 0;
  width: var(--nav-w);
  z-index: 15;
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px 10px;
  background: var(--surface);
  border-right: 1px solid var(--border);
}
.item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 10px 12px;
  border-radius: 12px;
  border: 2px solid transparent;
  color: var(--muted);
  text-decoration: none;
  font-weight: 700;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.item:hover {
  background: var(--surface-2);
}
.item.router-link-active {
  color: var(--accent);
  border-color: var(--accent);
  background: color-mix(in srgb, var(--accent) 12%, transparent);
}
@media (max-width: 900px) {
  .item {
    justify-content: center;
  }
  .item span {
    display: none;
  }
}
@media (max-width: 600px) {
  .nav {
    top: auto;
    width: 100%;
    height: var(--nav-b);
    flex-direction: row;
    justify-content: space-around;
    padding: 6px;
    border-right: 0;
    border-top: 1px solid var(--border);
  }
  .item {
    padding: 8px 10px;
  }
}
</style>
