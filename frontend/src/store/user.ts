import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { TokenPair } from '@/types'

const TOKEN_KEY = 'blog_access_token'
const REFRESH_KEY = 'blog_refresh_token'

export const useUserStore = defineStore('user', () => {
  const router = useRouter()

  const token = ref<string>(localStorage.getItem(TOKEN_KEY) ?? '')
  const refreshToken = ref<string>(localStorage.getItem(REFRESH_KEY) ?? '')

  const isAuthenticated = computed(() => Boolean(token.value))

  function setTokens(pair: TokenPair) {
    token.value = pair.access_token
    refreshToken.value = pair.refresh_token
    localStorage.setItem(TOKEN_KEY, pair.access_token)
    localStorage.setItem(REFRESH_KEY, pair.refresh_token)
  }

  function logout() {
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    router.push('/login')
  }

  return {
    token,
    refreshToken,
    isAuthenticated,
    setTokens,
    logout,
  }
})
