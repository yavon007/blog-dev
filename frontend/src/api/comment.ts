import request from '@/utils/request'
import type { ApiResponse, PagedData, Comment } from '@/types'

export const commentApi = {
  // 公开接口
  list(slug: string, params: { page?: number; page_size?: number } = {}) {
    return request.get<ApiResponse<PagedData<Comment>>>(`/posts/${slug}/comments`, { params })
  },

  create(slug: string, payload: { author_name: string; author_email: string; body: string; parent_comment_id?: number | null }) {
    return request.post<ApiResponse<Comment>>(`/posts/${slug}/comments`, payload)
  },

  // 管理接口
  adminList(params: { post_id?: number; status?: string; page?: number; page_size?: number } = {}) {
    return request.get<ApiResponse<PagedData<Comment>>>('/admin/comments', { params })
  },

  updateStatus(id: number, status: 'approved' | 'rejected') {
    return request.patch<ApiResponse<null>>(`/admin/comments/${id}`, { status })
  },

  delete(id: number) {
    return request.delete<ApiResponse<null>>(`/admin/comments/${id}`)
  },
}
