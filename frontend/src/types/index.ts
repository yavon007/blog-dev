// API 类型定义
export interface ApiResponse<T = unknown> {
    code: number
    data: T
    message: string
}

export interface PagedData<T> {
    items: T[]
    total: number
    page: number
    size: number
}

export interface Post {
    id: number
    title: string
    slug: string
    summary: string
    content_md: string
    content_html_cached: string
    cover_url: string
    status: 'draft' | 'published'
    published_at: string | null
    category_id: number | null
    category_name: string
    author_id: number
    tags: Tag[]
    created_at: string
    updated_at: string
}

export interface Tag {
    id: number
    name: string
    slug: string
    post_count?: number
    created_at?: string
}

export interface Category {
    id: number
    name: string
    slug: string
    description: string
    post_count?: number
    created_at?: string
}

export interface Comment {
    id: number
    post_id: number
    parent_comment_id: number | null
    author_name: string
    author_email: string
    body: string
    status: 'pending' | 'approved' | 'rejected'
    created_at: string
    updated_at: string
}

export interface TokenPair {
    access_token: string
    refresh_token: string
    expires_in: number
}

export interface CaptchaResponse {
    id: string
    image: string // data:image/png;base64,...
    expires_in: number
}

export interface LoginResponse extends TokenPair {
    captcha_required?: boolean
    failures?: {
        ip: number
        email: number
    }
}

export interface PostListFilter {
    page?: number
    page_size?: number
    category?: string
    tag?: string
    q?: string
    status?: string
}

export interface CreatePostPayload {
    title: string
    slug: string
    summary: string
    content_md: string
    cover_url: string
    status: 'draft' | 'published'
    category_id: number | null
    tag_ids: number[]
}

export interface CreateCategoryPayload {
    name: string
    slug: string
    description?: string
}

export interface CreateTagPayload {
    name: string
    slug: string
}
