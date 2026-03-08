import axios, { type AxiosError, type AxiosResponse } from 'axios'
import { useUserStore } from '@/store/user'
import type { ApiResponse } from '@/types'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：注入 JWT Token
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error: AxiosError) => Promise.reject(error),
)

// 响应拦截器：统一处理错误
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    return response.data as unknown as AxiosResponse
  },
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
    }
    const message = error.response?.data?.message ?? error.message
    return Promise.reject(new Error(message))
  },
)

export default request
