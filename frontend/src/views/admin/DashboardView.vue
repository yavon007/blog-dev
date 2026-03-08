<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { postApi } from '@/api/post'
import { commentApi } from '@/api/comment'

const stats = ref({ posts: 0, comments: 0, pendingComments: 0 })

onMounted(async () => {
  try {
    const [postsRes, commentsRes, pendingRes] = await Promise.all([
      postApi.adminList({ page_size: 1 }) as unknown as { data: { total: number } },
      commentApi.adminList({ page_size: 1 }) as unknown as { data: { total: number } },
      commentApi.adminList({ status: 'pending', page_size: 1 }) as unknown as { data: { total: number } },
    ])
    stats.value = {
      posts: postsRes.data.total,
      comments: commentsRes.data.total,
      pendingComments: pendingRes.data.total,
    }
  } catch { /* ignore */ }
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">仪表盘</h1>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
      <div class="card p-6">
        <div class="flex items-center gap-3">
          <span class="i-carbon-document text-3xl text-primary-500" />
          <div>
            <p class="text-2xl font-bold">{{ stats.posts }}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">文章总数</p>
          </div>
        </div>
      </div>
      <div class="card p-6">
        <div class="flex items-center gap-3">
          <span class="i-carbon-chat text-3xl text-green-500" />
          <div>
            <p class="text-2xl font-bold">{{ stats.comments }}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">评论总数</p>
          </div>
        </div>
      </div>
      <div class="card p-6">
        <div class="flex items-center gap-3">
          <span class="i-carbon-warning text-3xl text-orange-500" />
          <div>
            <p class="text-2xl font-bold">{{ stats.pendingComments }}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">待审核评论</p>
          </div>
        </div>
      </div>
    </div>

    <div class="flex gap-4">
      <RouterLink to="/admin/posts/new" class="btn-primary">
        <span class="i-carbon-add mr-1" />新建文章
      </RouterLink>
      <RouterLink v-if="stats.pendingComments > 0" to="/admin/comments?status=pending" class="btn-secondary">
        审核评论 ({{ stats.pendingComments }})
      </RouterLink>
    </div>
  </div>
</template>
