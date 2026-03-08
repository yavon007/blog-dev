<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { authApi } from '@/api/auth'
import type { TokenPair } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    const res = await authApi.login(form.value.email, form.value.password) as unknown as { data: TokenPair }
    userStore.setTokens(res.data)
    const redirect = (route.query.redirect as string) ?? '/admin'
    router.push(redirect)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '登录失败，请检查账号密码'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">📝 Blog Admin</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">管理员登录</p>
      </div>

      <div class="card p-8">
        <div v-if="error" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 rounded-lg text-sm">
          {{ error }}
        </div>

        <form class="space-y-4" @submit.prevent="handleLogin">
          <div>
            <label class="block text-sm font-medium mb-1.5">邮箱</label>
            <input
              v-model="form.email"
              type="email"
              required
              autocomplete="email"
              class="input-base"
              placeholder="admin@example.com"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1.5">密码</label>
            <input
              v-model="form.password"
              type="password"
              required
              autocomplete="current-password"
              class="input-base"
              placeholder="••••••••"
            />
          </div>
          <button
            type="submit"
            :disabled="loading"
            class="btn-primary w-full justify-center disabled:opacity-60"
          >
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
