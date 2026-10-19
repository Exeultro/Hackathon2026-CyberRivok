<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import {
  CategoryScale,
  Chart,
  Filler,
  LineController,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip,
} from 'chart.js'
import { useTheme } from '@/useTheme'

Chart.register(CategoryScale, LinearScale, LineController, LineElement, PointElement, Filler, Tooltip)

const props = defineProps<{
  labels: string[]
  values: number[]
  title: string
  unit: string
  color: string
  yMin?: number
  yMax?: number
}>()

const { theme } = useTheme()
const canvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

function token(name: string) {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

function build() {
  chart?.destroy()
  if (!canvas.value || !props.values.length) return
  const color = token(props.color)
  const grid = token('--border')
  const muted = token('--muted')
  const surface = token('--surface')
  chart = new Chart(canvas.value, {
    type: 'line',
    data: {
      labels: props.labels,
      datasets: [
        {
          label: props.title,
          data: props.values,
          borderColor: color,
          backgroundColor: color + '26',
          borderWidth: 2,
          pointRadius: 4,
          pointHoverRadius: 7,
          pointBackgroundColor: color,
          pointBorderColor: surface,
          pointBorderWidth: 2,
          tension: 0.3,
          fill: true,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      animation: { duration: 250 },
      interaction: { mode: 'index', intersect: false },
      plugins: {
        legend: { display: false },
        tooltip: {
          displayColors: false,
          callbacks: { label: (c) => `${props.title}: ${c.parsed.y}${props.unit}` },
        },
      },
      scales: {
        x: { grid: { display: false }, border: { color: grid }, ticks: { color: muted, maxRotation: 0, autoSkip: true } },
        y: {
          min: props.yMin,
          max: props.yMax,
          grid: { color: grid },
          border: { display: false },
          ticks: { color: muted, callback: (v) => `${v}${props.unit}` },
        },
      },
    },
  })
}

onMounted(build)
onBeforeUnmount(() => chart?.destroy())
watch([() => props.labels, () => props.values, theme], build)
</script>

<template>
  <div class="chart">
    <canvas v-show="values.length" ref="canvas" role="img" :aria-label="title" />
    <p v-if="!values.length" class="empty">Пока нет данных. Пройди пару тренировок.</p>
  </div>
</template>

<style scoped>
.chart {
  position: relative;
  height: 220px;
}
.empty {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: var(--muted);
  margin: 0;
}
</style>
