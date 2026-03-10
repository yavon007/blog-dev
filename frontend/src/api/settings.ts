import request from '@/utils/request'

export const settingsApi = {
  getSeo: () => request.get('/api/v1/seo/meta'),
  updateSeo: (data: {
    site_title: string
    site_description: string
    default_meta_title: string
    default_meta_description: string
    og_image_url?: string
  }) => request.put('/api/v1/admin/seo/meta', data),
}
