<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import { authApi } from '@/api/auth'

const userStore = useUserStore()
const appStore = useAppStore()
const route = useRoute()

const navItems = [
  { to: '/admin', label: '仪表盘', icon: 'i-carbon-dashboard' },
  { to: '/admin/posts', label: '文章', icon: 'i-carbon-document' },
  { to: '/admin/categories', label: '分类', icon: 'i-carbon-folder' },
  { to: '/admin/tags', label: '标签', icon: 'i-carbon-tag' },
  { to: '/admin/comments', label: '评论', icon: 'i-carbon-chat' },
]

// 判断是否激活：仪表盘需要精确匹配，其他前缀匹配
function isActive(item: { to: string }) {
  if (item.to === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(item.to)
}

async function handleLogout() {
  try {
    await authApi.logout(userStore.refreshToken)
  } catch {
    // ignore
  }
  userStore.logout()
}
</script>

<template>
  <div class="min-h-screen flex bg-gray-100 dark:bg-gray-950">
    <!-- 侧边栏 -->
    <aside
      class="fixed inset-y-0 left-0 z-50 flex flex-col w-60 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700 transition-transform"
      :class="appStore.isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo -->
      <div class="h-16 flex items-center px-6 border-b border-gray-200 dark:border-gray-700">
        <RouterLink to="/admin" class="text-lg font-bold text-primary-600 dark:text-primary-400">
          📝 Blog Admin
        </RouterLink>
      </div>

      <!-- 导航 -->
      <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
          :class="isActive(item) ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-400' : 'text-gray-600 dark:text-gray-300'"
        >
          <span :class="item.icon" class="text-lg flex-shrink-0" />
          {{ item.label }}
        </RouterLink>
      </nav>

      <!-- 底部操作 -->
      <div class="px-3 py-4 border-t border-gray-200 dark:border-gray-700 space-y-1">
        <RouterLink
          to="/"
          target="_blank"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <span class="i-carbon-launch text-lg" />
          查看前台
        </RouterLink>
        <button
          class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
          @click="handleLogout"
        >
          <span class="i-carbon-logout text-lg" />
          退出登录
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col" :class="appStore.isSidebarOpen ? 'ml-60' : 'ml-0'">
      <!-- 顶部栏 -->
      <header class="sticky top-0 z-40 h-16 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 flex items-center px-6 gap-4">
        <button
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          @click="appStore.toggleSidebar()"
        >
          <span class="i-carbon-menu text-xl" />
        </button>

        <div class="flex-1" />

        <button
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          @click="appStore.toggleDark()"
        >
          <span v-if="appStore.isDark" class="i-carbon-sun text-xl" />
          <span v-else class="i-carbon-moon text-xl" />
        </button>
      </header>

      <!-- 页面内容 -->
      <main class="flex-1 p-6 overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>
