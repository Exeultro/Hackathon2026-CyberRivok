<template>
  <div class="kb">
    <div
      v-for="(row, rowIndex) in rows"
      :key="rowIndex"
      class="kb__row"
    >
      <span
        v-for="k in row"
        :key="k"
        class="kb__key"
        :class="[getKeyClass(k), { 'is-active': k === activeKey }]"
      >
        {{ k }}
      </span>
    </div>
    <p class="kb__legend">
      Нажми: <strong class="text-violet">левый указательный</strong>
    </p>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const rows = [
  ['Ё','1','2','3','4','5','6','7','8','9','0','-','='],
  ['Й','Ц','У','К','Е','Н','Г','Ш','Щ','З','Х','Ъ','\\'],
  ['Ф','Ы','В','А','П','Р','О','Л','Д','Ж','Э',']'],
  ['Shift','Я','Ч','С','М','И','Т','Ь','Б','Ю','.','Shift']
]

// Все символы, которые есть на клавиатуре (для отслеживания нажатий)
const allKeys = rows.flat()
const activeKey = ref('П') // По умолчанию подсвечена П

// Слушаем нажатия клавиш на клавиатуре компьютера
function handleKeyDown(e) {
  // Переводим в верхний регистр для русской раскладки
  const key = e.key.toUpperCase()
  
  // Если нажатая клавиша есть на нашей виртуальной клавиатуре — подсвечиваем её
  if (allKeys.includes(key)) {
    activeKey.value = key
  } else if (e.key === ' ') {
    // Если нажат пробел, можно снять подсветку (или сделать что-то ещё)
    activeKey.value = ''
  }
}

onMounted(() => window.addEventListener('keydown', handleKeyDown))
onUnmounted(() => window.removeEventListener('keydown', handleKeyDown))

// Распределяем цвета по пальцам (RGB стиль)
function getKeyClass(k) {
  // Левая рука (мизинец, безымянный, средний, указательный)
  const leftPinky = ['Ё','Й','Ф','Shift','1']
  const leftRing = ['Ц','Ы','2','Я']
  const leftMiddle = ['У','В','3','Ч']
  const leftIndex = ['К','Е','Н','Г','А','П','4','5','Р','О','С','М','И']

  // Правая рука (указательный, средний, безымянный, мизинец)
  const rightIndex = ['Ш','Щ','З','Х','Ъ','Д','Ж','Э','Т','Ь','Б','6','7']
  const rightMiddle = ['О','Л','8','Ю']
  const rightRing = ['Д','Ж','9','.']
  const rightPinky = ['Х','Ъ','\\',']','=','0','-','Shift']

  if (leftPinky.includes(k)) return 'key-pinky-l'
  if (leftRing.includes(k)) return 'key-ring-l'
  if (leftMiddle.includes(k)) return 'key-middle-l'
  if (leftIndex.includes(k)) return 'key-index-l'

  if (rightIndex.includes(k)) return 'key-index-r'
  if (rightMiddle.includes(k)) return 'key-middle-r'
  if (rightRing.includes(k)) return 'key-ring-r'
  if (rightPinky.includes(k)) return 'key-pinky-r'

  return ''
}
</script>

<style scoped>
.kb {
  display: inline-flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px;
  background: var(--bg-card);
  border-radius: var(--radius-md);
  box-shadow: 0 10px 40px rgba(0,0,0,0.3);
}

.kb__row {
  display: flex;
  gap: 6px;
  justify-content: center;
}

.kb__key {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 38px;
  height: 40px;
  padding: 0 6px;
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 500;
  border-radius: 6px;
  background: #2a2a32;
  color: #a0a0b0;
  border: 1px solid rgba(255,255,255,0.05);
  border-bottom-width: 3px; /* Толщина цветной подложки */
  transition: all 0.1s ease;
  position: relative;
}

/* Эффект при нажатии реальной клавиши */
.kb__key.is-active {
  background: #ffffff !important; /* Ярко-белый фон при нажатии */
  color: #1c1c22 !important;      /* Темный текст */
  transform: translateY(2px);     /* Кнопка "вдавливается" */
  border-bottom-width: 1px;       /* Подложка уменьшается */
  box-shadow: 0 0 20px currentColor; /* Свечение цвета пальца */
  z-index: 10;
}

/* --- Цвета подложек для каждого пальца (RGB стиль) --- */
/* Левая рука */
.key-pinky-l { border-bottom-color: #ff007f; } /* Розовый */
.key-ring-l  { border-bottom-color: #ff00ff; } /* Маджента */
.key-middle-l{ border-bottom-color: #8a2be2; } /* Фиолетовый */
.key-index-l { border-bottom-color: #00d2ff; } /* Голубой */

/* Правая рука */
.key-index-r { border-bottom-color: #00ff88; } /* Зеленый */
.key-middle-r{ border-bottom-color: #adff2f; } /* Лаймовый */
.key-ring-r  { border-bottom-color: #ffd700; } /* Желтый */
.key-pinky-r { border-bottom-color: #ff4500; } /* Оранжево-красный */

/* Красивое свечение для активной клавиши, зависит от цвета её подложки */
.kb__key.is-active.key-pinky-l { color: #ff007f; }
.kb__key.is-active.key-ring-l  { color: #ff00ff; }
.kb__key.is-active.key-middle-l{ color: #8a2be2; }
.kb__key.is-active.key-index-l { color: #00d2ff; }
.kb__key.is-active.key-index-r { color: #00ff88; }
.kb__key.is-active.key-middle-r{ color: #adff2f; }
.kb__key.is-active.key-ring-r  { color: #ffd700; }
.kb__key.is-active.key-pinky-r { color: #ff4500; }

/* Легенда снизу */
.kb__legend {
  text-align: center;
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 12px;
}

.text-violet {
  color: #8a2be2;
  text-shadow: 0 0 10px rgba(138,43,226,0.5);
}
</style>