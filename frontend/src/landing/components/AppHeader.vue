<script setup>
import { ref } from 'vue'
const lang = ref('ru')
</script>

<template>
  <header class="header">
    <div class="container header__inner">
      <RouterLink to="/" class="logo">
        <span class="logo__mark">⌨</span>
        <span class="logo__text">Печать героя</span>
      </RouterLink>
      <nav class="header__nav">
        <div class="lang-switch">
          <button :class="{ active: lang === 'ru' }" @click="lang = 'ru'">RU</button>
          <button :class="{ active: lang === 'en' }" @click="lang = 'en'">EN</button>
        </div>
        <RouterLink to="/map" class="btn btn--accent">Играть</RouterLink>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.header { position: sticky; top: 0; z-index: 50; background: rgba(28, 28, 34, .85); backdrop-filter: blur(12px); border-bottom: 1px solid rgba(255,255,255,.06); }
.header__inner { display: flex; align-items: center; justify-content: space-between; height: 64px; }
.logo { display: flex; align-items: center; gap: 10px; font-weight: 700; font-size: 18px; color: var(--accent); }
.logo__mark { display: grid; place-items: center; width: 34px; height: 34px; border-radius: var(--radius-sm); background: rgba(245,165,36,.12); color: var(--accent); font-size: 18px; }
.header__nav { display: flex; align-items: center; gap: 16px; }
.lang-switch { display: inline-flex; background: var(--bg-card); border-radius: var(--radius-sm); padding: 3px; gap: 2px; }
.lang-switch button { padding: 4px 12px; font-size: 12px; font-weight: 600; border-radius: 6px; color: var(--text-muted); transition: .15s; }
.lang-switch button.active { background: var(--accent); color: #1c1c22; }
</style>