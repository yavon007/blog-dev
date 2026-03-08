<script setup lang="ts">
import { useAppStore } from '@/store/app'
import { RouterLink, RouterView } from 'vue-router'

const appStore = useAppStore()
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 text-gray-900 dark:text-gray-100 transition-colors">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-50 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 shadow-sm">
      <div class="max-w-4xl mx-auto px-4 h-16 flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold text-primary-600 dark:text-primary-400 hover:text-primary-700">
          📝 Blog
        </RouterLink>

        <nav class="hidden md:flex items-center gap-6">
          <RouterLink to="/" class="text-sm font-medium hover:text-primary-600 dark:hover:text-primary-400 transition-colors">
            文章
          </RouterLink>
          <RouterLink to="/tags" class="text-sm font-medium hover:text-primary-600 dark:hover:text-primary-400 transition-colors">
            标签
          </RouterLink>
        </nav>

        <div class="flex items-center gap-3">
          <!-- 主题切换 -->
          <button
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            :aria-label="appStore.isDark ? '切换浅色模式' : '切换深色模式'"
            @click="appStore.toggleDark()"
          >
            <span v-if="appStore.isDark" class="i-carbon-sun text-xl" />
            <span v-else class="i-carbon-moon text-xl" />
          </button>

          <RouterLink
            to="/admin"
            class="hidden md:block text-sm text-gray-500 hover:text-primary-600 dark:text-gray-400 dark:hover:text-primary-400 transition-colors"
          >
            管理
          </RouterLink>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="max-w-4xl mx-auto px-4 py-8">
      <RouterView />
    </main>

    <!-- 页脚 -->
    <footer class="border-t border-gray-200 dark:border-gray-700 mt-16">
      <div class="max-w-4xl mx-auto px-4 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
        <p>© {{ new Date().getFullYear() }} Blog. Built with Vue 3 + Go.</p>
      </div>
    </footer>
  </div>
</template>
