<template>
  <footer class="footer">
    <div class="container">
      <div class="cta">
        <h2>Готов стать героем?</h2>
        <p>Бесплатно. Без установки. Прогресс сохраняется между устройствами.</p>
        <RouterLink to="/map" class="btn btn--accent btn--lg">Начать играть</RouterLink>
      </div>
      <div class="footer__bottom">
      <span>© 2026 «Печать героя»</span>
        <nav><a href="#">Документация</a><a href="#">GitHub</a><a href="#">Контакты</a></nav>
      </div>
    </div>
  </footer>
</template>
<style scoped>
.footer { padding: 80px 0 32px; }
.cta { text-align: center; padding: 56px 24px; background: radial-gradient(circle at top, rgba(245,165,36,.15), transparent 70%), var(--bg-card); border-radius: var(--radius-lg); margin-bottom: 48px; }
.cta h2 { font-size: 32px; }
.cta p { color: var(--text-muted); margin: 12px 0 28px; }
.footer__bottom { display: flex; justify-content: space-between; align-items: center; font-size: 13px; color: var(--text-muted); padding-top: 24px; border-top: 1px solid rgba(255,255,255,.06); }
.footer__bottom nav { display: flex; gap: 20px; }
.footer__bottom a:hover { color: var(--accent); }
</style>