import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usePreferredDark } from '@vueuse/core'

export const useAppStore = defineStore('app', () => {
  const prefersDark = usePreferredDark()
  const isDark = ref<boolean>(
    localStorage.getItem('blog_theme') === 'dark' ||
    (localStorage.getItem('blog_theme') === null && prefersDark.value),
  )
  const isSidebarOpen = ref<boolean>(true)

  function toggleDark() {
    isDark.value = !isDark.value
    localStorage.setItem('blog_theme', isDark.value ? 'dark' : 'light')
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  function toggleSidebar() {
    isSidebarOpen.value = !isSidebarOpen.value
  }

  // 初始化时应用主题
  function initTheme() {
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  return {
    isDark,
    isSidebarOpen,
    toggleDark,
    toggleSidebar,
    initTheme,
  }
})
