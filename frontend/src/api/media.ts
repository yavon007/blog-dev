import request from '@/utils/request'
import type { ApiResponse, PagedData, MediaFile, UploadResult } from '@/types'

export const mediaApi = {
  list(page = 1, pageSize = 20, mimeType?: string) {
    const params: Record<string, unknown> = { page, page_size: pageSize }
    if (mimeType) params.mime_type = mimeType
    return request.get<ApiResponse<PagedData<MediaFile>>>('/admin/media', { params })
  },

  upload(file: File, altText?: string) {
    const formData = new FormData()
    formData.append('file', file)
    if (altText) formData.append('alt_text', altText)
    return request.post<ApiResponse<UploadResult>>('/admin/media', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  delete(id: number) {
    return request.delete<ApiResponse<null>>(`/admin/media/${id}`)
  },
}
