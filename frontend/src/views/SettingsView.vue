<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api'
import { useGameStore } from '@/stores/game'
import { useTheme } from '@/useTheme'
import { sfx } from '@/useSound'
import Switch from '@/components/Switch.vue'

const game = useGameStore()
const { theme, toggle } = useTheme()

const nickname = ref('')
const saved = ref(false)
const confirming = ref(false)
const resetDone = ref(false)

onMounted(async () => {
  if (!game.loaded) await game.load()
  nickname.value = game.user?.nickname ?? ''
})

const valid = computed(() => nickname.value.trim().length >= 2 && nickname.value.trim().length <= 24)
const changed = computed(() => nickname.value.trim() !== game.user?.nickname)

async function saveNick() {
  if (!valid.value || !changed.value) return
  game.user = await api.updateMe({ nickname: nickname.value })
  nickname.value = game.user.nickname
  saved.value = true
  window.setTimeout(() => (saved.value = false), 2000)
}

function setTheme(want: 'dark' | 'light') {
  if (theme.value !== want) toggle()
}

async function reset() {
  await api.reset?.()
  await game.load()
  nickname.value = game.user?.nickname ?? ''
  confirming.value = false
  resetDone.value = true
}
</script>

<template>
  <main class="settings">
    <h1>Настройки</h1>

    <section class="panel">
      <h2>Профиль</h2>
      <label class="field">
        <span>Никнейм</span>
        <div class="inline">
          <input v-model="nickname" maxlength="24" autocomplete="off" @keydown.enter="saveNick" />
          <button class="btn btn-primary" :disabled="!valid || !changed" @click="saveNick">Сохранить</button>
        </div>
        <small v-if="!valid" class="err">От 2 до 24 символов</small>
        <small v-else-if="saved" class="ok">Сохранено</small>
      </label>
    </section>

    <section class="panel">
      <h2>Набор</h2>
      <div class="row">
        <span class="text"><b>Язык набора</b><small>Тексты и подсказка клавиатуры</small></span>
        <div class="seg" role="group" aria-label="Язык набора">
          <button :aria-pressed="game.lang === 'ru'" @click="game.lang = 'ru'">RU</button>
          <button :aria-pressed="game.lang === 'en'" @click="game.lang = 'en'">EN</button>
        </div>
      </div>
      <div class="row">
        <span class="text"><b>Подсветка текста</b><small>Выделять только символ или всё текущее слово</small></span>
        <div class="seg" role="group" aria-label="Подсветка текста">
          <button :aria-pressed="game.unit === 'char'" @click="game.unit = 'char'">Символ</button>
          <button :aria-pressed="game.unit === 'word'" @click="game.unit = 'word'">Слово</button>
        </div>
      </div>
      <Switch v-model="game.realMode" label="Реальный режим" hint="Ошибочная буква впечатывается красным, исправляй Backspace" />
      <Switch v-model="game.autoPair" label="Автоскобки в коде" hint="Скобки и кавычки закрываются сами, как в редакторе" />
      <Switch v-model="game.showKeyboard" label="Клавиатура с пальцами" hint="Подсказка, какой палец нажимает клавишу" />
    </section>

    <section class="panel">
      <h2>Звук</h2>
      <Switch v-model="game.sound" label="Звуковые эффекты" hint="Клики, ошибки, комбо и награды" />
      <button class="btn btn-ghost test" @click="sfx.combo(3)">Проверить звук</button>
    </section>

    <section class="panel">
      <h2>Оформление</h2>
      <div class="row">
        <span class="text"><b>Тема</b></span>
        <div class="seg" role="group" aria-label="Тема">
          <button :aria-pressed="theme === 'dark'" @click="setTheme('dark')">Тёмная</button>
          <button :aria-pressed="theme === 'light'" @click="setTheme('light')">Светлая</button>
        </div>
      </div>
    </section>

    <section class="panel danger">
      <h2>Данные</h2>
      <p class="muted">Сбросит уровень, звёзды, достижения и всю статистику. Настройки останутся.</p>
      <template v-if="!confirming">
        <button class="btn btn-ghost warn" @click="confirming = true">Сбросить прогресс</button>
      </template>
      <div v-else class="confirm">
        <span>Точно сбросить? Это нельзя отменить.</span>
        <button class="btn btn-danger" @click="reset">Да, сбросить</button>
        <button class="btn btn-ghost" @click="confirming = false">Отмена</button>
      </div>
      <p v-if="resetDone" class="ok">Прогресс сброшен.</p>
    </section>
  </main>
</template>

<style scoped>
.settings {
  max-width: 680px;
  margin: 0 auto;
  padding: 16px;
}
h1 {
  margin: 0 0 14px;
}
h2 {
  margin: 0 0 12px;
  font-size: 14px;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--muted);
}
.panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 12px;
  display: grid;
  gap: 14px;
}
.panel h2 {
  margin: 0;
}
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}
.text {
  display: grid;
}
small,
.muted {
  color: var(--muted);
}
.muted {
  margin: 0;
}
.seg {
  display: inline-flex;
  flex: none;
  background: var(--surface-2);
  border-radius: 12px;
  padding: 3px;
  gap: 3px;
}
.seg button {
  border: 0;
  background: transparent;
  color: var(--muted);
  padding: 7px 14px;
  border-radius: 9px;
  cursor: pointer;
  font-weight: 700;
}
.seg button[aria-pressed='true'] {
  background: var(--accent);
  color: var(--accent-ink);
}
.field {
  display: grid;
  gap: 6px;
}
.inline {
  display: flex;
  gap: 8px;
}
input {
  flex: 1;
  min-width: 0;
  font: inherit;
  padding: 10px 12px;
  border-radius: 10px;
  border: 2px solid var(--border);
  background: var(--bg);
  color: var(--text);
}
.btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.err {
  color: var(--err);
}
.ok {
  color: var(--ok);
  margin: 0;
}
.test {
  justify-self: start;
}
.warn {
  justify-self: start;
  color: var(--err);
  border-color: var(--err);
}
.btn-danger {
  background: var(--err);
  color: #fff;
}
.confirm {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}
.confirm span {
  flex: 1 1 100%;
}
</style>
