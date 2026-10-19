import { onBeforeUnmount, onMounted, ref } from 'vue'
import type { Ref } from 'vue'

export function useViewBand(root: Ref<HTMLElement | null>, margin = 350, step = 150) {
  const band = ref({ top: -Infinity, bottom: Infinity })
  let raf = 0

  function update() {
    raf = 0
    const el = root.value
    if (!el) return
    const r = el.getBoundingClientRect()
    const top = Math.floor((-r.top - margin) / step) * step
    const bottom = Math.ceil((-r.top + window.innerHeight + margin) / step) * step
    if (top !== band.value.top || bottom !== band.value.bottom) band.value = { top, bottom }
  }

  function schedule() {
    if (!raf) raf = requestAnimationFrame(update)
  }

  onMounted(() => {
    update()
    window.addEventListener('scroll', schedule, { passive: true })
    window.addEventListener('resize', schedule)
  })
  onBeforeUnmount(() => {
    window.removeEventListener('scroll', schedule)
    window.removeEventListener('resize', schedule)
    if (raf) cancelAnimationFrame(raf)
  })

  return band
}
