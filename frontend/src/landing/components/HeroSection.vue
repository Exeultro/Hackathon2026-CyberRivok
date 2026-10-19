<script setup>
import TypingDemo from './TypingDemo.vue'
import KeyboardHint from './KeyboardHint.vue'
</script>
<template>
  <section class="hero">
    <div class="container hero__grid">
      <div class="hero__left">
      <span class="hero__badge">⚡ Твой путь в мир слепой печати</span>
        <h1 class="hero__title">Стань <span class="accent">героем</span><br> слепой печати</h1>
        <p class="hero__subtitle">Тренажёр с уровнями, достижениями, комбо и таблицей лидеров. Печатай быстрее, прокачивай героя и возвращайся каждый день.</p>
        <div class="hero__actions">
          <RouterLink to="/map" class="btn btn--accent btn--lg">Начать тренировку →</RouterLink>
          <a href="#features" class="btn btn--ghost btn--lg">Как это работает</a>
        </div>
        <div class="hero__meta"><span>🎮 4 режима</span><span>🏆 15 достижений</span><span>🗺 4 карты</span></div>
      </div>
      <div class="hero__right">
        <TypingDemo />
        <div class="hero__keyboard"><KeyboardHint /></div>
      </div>
    </div>
  </section>
</template>
<style scoped>
.hero { padding: 72px 0 80px; position: relative; overflow: hidden; }
.hero::before { content: ''; position: absolute; top: -200px; right: -200px; width: 600px; height: 600px; background: radial-gradient(circle, rgba(245,165,36,.18), transparent 60%); pointer-events: none; }
.hero__grid { display: grid; grid-template-columns: 1fr 1.1fr; gap: 56px; align-items: center; position: relative; }
.hero__badge { display: inline-block; padding: 6px 12px; border-radius: 999px; background: rgba(245,165,36,.12); color: var(--accent); font-size: 13px; font-weight: 600; margin-bottom: 20px; }
.hero__title { font-size: clamp(36px, 5vw, 60px); font-weight: 800; line-height: 1.05; letter-spacing: -.02em; }
.hero__title .accent { color: var(--accent); text-shadow: 0 0 40px rgba(245,165,36,.5); }
.hero__subtitle { margin-top: 20px; font-size: 17px; color: var(--text-muted); max-width: 480px; }
.hero__actions { display: flex; gap: 12px; margin-top: 32px; }
.hero__meta { display: flex; gap: 20px; margin-top: 28px; font-size: 13px; color: var(--text-muted); }
.hero__right { display: flex; flex-direction: column; gap: 20px; }
.hero__keyboard { display: flex; justify-content: center; }
@media (max-width: 900px) { .hero__grid { grid-template-columns: 1fr; } }
</style>