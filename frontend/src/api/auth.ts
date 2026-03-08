import request from '@/utils/request'
import type { ApiResponse, TokenPair } from '@/types'

export const authApi = {
  login(email: string, password: string) {
    return request.post<ApiResponse<TokenPair>>('/admin/auth/login', { email, password })
  },

  refresh(refreshToken: string) {
    return request.post<ApiResponse<TokenPair>>('/admin/auth/refresh', { refresh_token: refreshToken })
  },

  logout(refreshToken: string) {
    return request.post<ApiResponse<null>>('/admin/auth/logout', { refresh_token: refreshToken })
  },
}
