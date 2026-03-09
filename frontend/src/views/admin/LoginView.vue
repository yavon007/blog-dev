<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { authApi } from '@/api/auth'
import request, { ApiError } from '@/utils/request'
import type { TokenPair, CaptchaResponse } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({ email: '', password: '', captcha: '' })
const loading = ref(false)
const error = ref('')
const requireCaptcha = ref(false)
const captchaId = ref('')
const captchaImg = ref('')

async function loadCaptcha() {
  try {
    // 验证码需要传 email 参数
    const res = await request.get<unknown, { data: CaptchaResponse }>(`/auth/captcha?email=${encodeURIComponent(form.value.email)}`)
    captchaId.value = res.data.id
    captchaImg.value = res.data.image
  } catch {
    error.value = '无法加载验证码，请刷新重试'
  }
}

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    const payload: { email: string; password: string; captcha_id?: string; captcha_code?: string } = {
      email: form.value.email,
      password: form.value.password,
    }
    if (requireCaptcha.value) {
      payload.captcha_id = captchaId.value
      payload.captcha_code = form.value.captcha
    }
    const res = await authApi.login(payload) as unknown as { data: TokenPair }
    userStore.setTokens(res.data)
    const redirect = (route.query.redirect as string) ?? '/admin'
    router.push(redirect)
  } catch (e) {
    // 检查是否需要验证码 (HTTP 428)
    if (e instanceof ApiError && e.status === 428) {
      const data = e.data as { captcha_required?: boolean } | undefined
      if (data?.captcha_required) {
        requireCaptcha.value = true
        form.value.captcha = ''
        await loadCaptcha()
        error.value = '登录失败次数过多，请输入验证码'
        return
      }
    }
    // 其他错误
    if (e instanceof ApiError && e.status === 401) {
      error.value = '邮箱或密码错误'
    } else {
      error.value = e instanceof Error ? e.message : '登录失败，请检查账号密码'
    }
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

          <!-- 验证码区域 -->
          <div v-if="requireCaptcha" class="animate-fade-in">
            <label class="block text-sm font-medium mb-1.5">验证码</label>
            <div class="flex items-center gap-3">
              <input
                v-model="form.captcha"
                type="text"
                required
                maxlength="6"
                autocomplete="off"
                class="input-base flex-1"
                placeholder="请输入6位验证码"
              />
              <img
                v-if="captchaImg"
                :src="captchaImg"
                @click="loadCaptcha"
                class="h-10 w-28 rounded cursor-pointer hover:opacity-80 transition-opacity border border-gray-200 dark:border-gray-700"
                alt="图形验证码，点击刷新"
                title="点击刷新验证码"
              />
            </div>
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

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-in;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
