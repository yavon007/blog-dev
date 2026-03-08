import request from '@/utils/request'
import type { ApiResponse, Category, Tag, CreateCategoryPayload, CreateTagPayload } from '@/types'

export const categoryApi = {
  list() {
    return request.get<ApiResponse<Category[]>>('/categories')
  },
  create(payload: CreateCategoryPayload) {
    return request.post<ApiResponse<Category>>('/admin/categories', payload)
  },
  update(id: number, payload: CreateCategoryPayload) {
    return request.put<ApiResponse<Category>>(`/admin/categories/${id}`, payload)
  },
  delete(id: number) {
    return request.delete<ApiResponse<null>>(`/admin/categories/${id}`)
  },
}

export const tagApi = {
  list() {
    return request.get<ApiResponse<Tag[]>>('/tags')
  },
  create(payload: CreateTagPayload) {
    return request.post<ApiResponse<Tag>>('/admin/tags', payload)
  },
  update(id: number, payload: CreateTagPayload) {
    return request.put<ApiResponse<Tag>>(`/admin/tags/${id}`, payload)
  },
  delete(id: number) {
    return request.delete<ApiResponse<null>>(`/admin/tags/${id}`)
  },
}
