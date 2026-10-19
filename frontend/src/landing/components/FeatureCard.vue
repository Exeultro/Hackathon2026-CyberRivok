<script setup>
defineProps({ icon: String, title: String, text: String })
</script>
<template>
  <article class="card">
    <div class="card__icon">{{ icon }}</div>
    <h3 class="card__title">{{ title }}</h3>
    <p class="card__text">{{ text }}</p>
  </article>
</template>
<style scoped>
.card { background: var(--bg-card); border-radius: var(--radius-md); padding: 24px; border: 1px solid rgba(255,255,255,.04); transition: .2s; }
.card:hover { background: var(--bg-card-hover); transform: translateY(-2px); border-color: rgba(245,165,36,.3); }
.card__icon { font-size: 26px; width: 48px; height: 48px; display: grid; place-items: center; border-radius: var(--radius-sm); background: rgba(245,165,36,.1); margin-bottom: 16px; }
.card__title { font-size: 17px; margin-bottom: 8px; }
.card__text { font-size: 14px; color: var(--text-muted); }
</style>