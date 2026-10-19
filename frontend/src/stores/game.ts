import { computed, ref, watch } from 'vue'
import type { Ref } from 'vue'
import { defineStore } from 'pinia'
import { api } from '@/api'
import type { GameMap, Lang, MapId, Progress, StatsSummary, User } from '@/api'

function persisted<T extends string | boolean>(key: string, fallback: T) {
  let initial = fallback
  try {
    const raw = localStorage.getItem(key)
    if (raw !== null) initial = (typeof fallback === 'boolean' ? raw === '1' : raw) as T
  } catch {
    // storage is unavailable, use the fallback
  }
  const state = ref(initial) as Ref<T>
  watch(
    () => state.value,
    (v) => {
      try {
        localStorage.setItem(key, typeof v === 'boolean' ? (v ? '1' : '0') : String(v))
      } catch {
        // setting just will not persist
      }
    },
  )
  return state
}

export const useGameStore = defineStore('game', () => {
  const maps = ref<GameMap[]>([])
  const progress = ref<Progress | null>(null)
  const summary = ref<StatsSummary | null>(null)
  const user = ref<User | null>(null)
  const currentMapId = ref<MapId>('journey')
  const lang = persisted<Lang>('hero.lang', 'ru')
  const sound = persisted<boolean>('hero.sound', true)
  const showKeyboard = persisted<boolean>('hero.keyboard', true)
  const realMode = persisted<boolean>('hero.realmode', false)
  const unit = persisted<'char' | 'word'>('hero.unit', 'char')
  const autoPair = persisted<boolean>('hero.autopair', true)
  const loaded = ref(false)
  const loadError = ref(false)

  async function load() {
    try {
      const [m, p, s, u] = await Promise.all([api.getMaps(), api.getProgress(), api.getSummary(), api.getMe()])
      maps.value = m
      progress.value = p
      summary.value = s
      user.value = u
      loaded.value = true
      loadError.value = false
    } catch {
      loadError.value = true
    }
  }

  const xpPercent = computed(() =>
    progress.value ? Math.round((progress.value.xp_in_level / progress.value.xp_to_next) * 100) : 0,
  )

  return { maps, progress, summary, user, currentMapId, lang, sound, showKeyboard, realMode, unit, autoPair, loaded, loadError, load, xpPercent }
})
