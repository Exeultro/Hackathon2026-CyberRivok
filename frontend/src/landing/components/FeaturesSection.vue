<script setup>
import FeatureCard from './FeatureCard.vue'
const features = [
  { icon: '📝', title: 'Слова', text: 'Случайные слова от 2 до 10 букв — разминка для пальцев.' },
  { icon: '📖', title: 'Предложения', text: 'Настоящие литературные отрывки: от технической прозы до стихов.' },
  { icon: '💻', title: 'Код', text: 'Фрагменты на Python, JavaScript, HTML/CSS с автоскобками.' },
  { icon: '🔣', title: 'Символы', text: 'Буквы, цифры и спецсимволы — для настоящих хардкорщиков.' },
]
</script>
<template>
  <section id="features" class="section">
    <div class="container">
      <h2 class="section__title">4 режима тренировки</h2>
      <p class="section__lead">От простых слов до кода на скорость</p>
      <div class="features">
        <FeatureCard v-for="f in features" :key="f.title" v-bind="f" />
      </div>
    </div>
  </section>
</template>
<style scoped>
.section { padding: 80px 0; }
.section__title { font-size: 32px; font-weight: 700; text-align: center; }
.section__lead { text-align: center; color: var(--text-muted); margin-top: 8px; margin-bottom: 40px; }
.features { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px; }
</style>