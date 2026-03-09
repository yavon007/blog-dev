import request from '@/utils/request'
import type { ApiResponse, PagedData, Post, PostListFilter, CreatePostPayload } from '@/types'

export interface ArchiveItem {
  year: number
  month: number
  count: number
}

export const postApi = {
  // 公开接口
  list(filter: PostListFilter = {}) {
    return request.get<ApiResponse<PagedData<Post>>>('/posts', { params: filter })
  },

  getBySlug(slug: string) {
    return request.get<ApiResponse<Post>>(`/posts/${slug}`)
  },

  getArchive() {
    return request.get<ApiResponse<ArchiveItem[]>>('/posts/archive')
  },

  listByYearMonth(year: number, month: number, page = 1, pageSize = 10) {
    return request.get<ApiResponse<PagedData<Post>>>(`/posts/archive/${year}/${month}`, {
      params: { page, page_size: pageSize },
    })
  },

  // 管理接口
  adminList(filter: PostListFilter = {}) {
    return request.get<ApiResponse<PagedData<Post>>>('/admin/posts', { params: filter })
  },

  adminGetById(id: number) {
    return request.get<ApiResponse<Post>>(`/admin/posts/${id}`)
  },

  create(payload: CreatePostPayload) {
    return request.post<ApiResponse<Post>>('/admin/posts', payload)
  },

  update(id: number, payload: CreatePostPayload) {
    return request.put<ApiResponse<Post>>(`/admin/posts/${id}`, payload)
  },

  updateStatus(id: number, status: 'draft' | 'published') {
    return request.patch<ApiResponse<null>>(`/admin/posts/${id}/status`, { status })
  },

  delete(id: number) {
    return request.delete<ApiResponse<null>>(`/admin/posts/${id}`)
  },
}
