import axios, { type AxiosError, type AxiosResponse } from 'axios'
import { useUserStore } from '@/store/user'
import type { ApiResponse } from '@/types'

// 自定义错误类，保留响应数据
export class ApiError extends Error {
  status: number
  code: number
  data?: unknown

  constructor(message: string, status: number, code: number, data?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.data = data
  }
}

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
    const { status, data } = error.response ?? {}
    const message = data?.message ?? error.message
    const code = data?.code ?? status ?? 0
    return Promise.reject(new ApiError(message, status ?? 0, code, data?.data))
  },
)

export default request
