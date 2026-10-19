<template>
  <section class="game">
    <div class="container game__grid">
      <div class="game__text">
        <h2>Ты — герой.<br>Твой прогресс — <span class="accent">сила</span>.</h2>
        <p>За каждую тренировку ты получаешь XP. Растёт уровень — открываются темы, режимы и достижения. Комбо и бонусы ускоряют прокачку.</p>
        <div class="level">
          <div class="level__top"><span>🔥 Уровень 5</span><span class="muted">114 / 250 XP</span></div>
          <div class="level__bar"><div class="level__fill" style="width: 46%"></div></div>
        </div>
      </div>
      <div class="achievements">
        <div class="ach"><div class="ach__icon">🏆</div><div><div class="ach__title">Печатный гений</div><div class="ach__text">100 WPM</div></div></div>
        <div class="ach"><div class="ach__icon">🛡️</div><div><div class="ach__title">Стальной палец</div><div class="ach__text">250 символов подряд без ошибок</div></div></div>
        <div class="ach"><div class="ach__icon">⚡</div><div><div class="ach__title">Скорость 300</div><div class="ach__text">300 CPM</div></div></div>
      </div>
    </div>
  </section>
</template>
<style scoped>
.game { padding: 80px 0; background: var(--bg-elev); }
.game__grid { display: grid; grid-template-columns: 1fr 1fr; gap: 56px; align-items: center; }
.game h2 { font-size: 36px; line-height: 1.15; }
.game .accent { color: var(--accent); }
.game p { margin-top: 16px; color: var(--text-muted); max-width: 420px; }
.level { margin-top: 32px; max-width: 380px; }
.level__top { display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 8px; }
.level__top .muted { color: var(--text-muted); }
.level__bar { height: 10px; background: rgba(255,255,255,.06); border-radius: 999px; overflow: hidden; }
.level__fill { height: 100%; background: linear-gradient(90deg, var(--accent), #ffcf6b); border-radius: 999px; box-shadow: 0 0 20px rgba(245,165,36,.6); }
.achievements { display: grid; gap: 14px; }
.ach { display: flex; align-items: center; gap: 16px; background: var(--bg-card); padding: 18px 20px; border-radius: var(--radius-md); border: 1px solid rgba(255,255,255,.04); }
.ach__icon { font-size: 22px; width: 44px; height: 44px; display: grid; place-items: center; background: rgba(245,165,36,.12); border-radius: var(--radius-sm); }
.ach__title { font-weight: 600; }
.ach__text { font-size: 13px; color: var(--text-muted); }
@media (max-width: 900px) { .game__grid { grid-template-columns: 1fr; } }
</style>