import { ref } from 'vue'

type Theme = 'dark' | 'light'
const KEY = 'hero.theme'

function readStored(): Theme {
  try {
    return localStorage.getItem(KEY) === 'light' ? 'light' : 'dark'
  } catch {
    return 'dark'
  }
}

const theme = ref<Theme>(readStored())
document.documentElement.dataset.theme = theme.value

export function useTheme() {
  function toggle() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    document.documentElement.dataset.theme = theme.value
    try {
      localStorage.setItem(KEY, theme.value)
    } catch {
      // storage is unavailable, theme just will not persist
    }
  }
  return { theme, toggle }
}
