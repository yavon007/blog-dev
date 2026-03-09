import request from '@/utils/request'
import type { ApiResponse, TokenPair, CaptchaResponse, LoginResponse } from '@/types'

export interface LoginPayload {
  email: string
  password: string
  captcha_id?: string
  captcha_code?: string
}

export const authApi = {
  login(payload: LoginPayload) {
    return request.post<ApiResponse<LoginResponse>>('/auth/login', payload)
  },

  getCaptcha() {
    return request.get<ApiResponse<CaptchaResponse>>('/auth/captcha')
  },

  refresh(refreshToken: string) {
    return request.post<ApiResponse<TokenPair>>('/auth/refresh', { refresh_token: refreshToken })
  },

  logout(refreshToken: string) {
    return request.post<ApiResponse<null>>('/admin/auth/logout', { refresh_token: refreshToken })
  },
}
