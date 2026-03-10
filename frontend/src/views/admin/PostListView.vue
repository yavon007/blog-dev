<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { postApi } from '@/api/post'
import type { Post } from '@/types'

const posts = ref<Post[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

async function fetchPosts() {
  loading.value = true
  try {
    const res = await postApi.adminList({ page: page.value, page_size: 15 }) as unknown as { data: { items: Post[]; total: number } }
    posts.value = res.data.items ?? []
    total.value = res.data.total ?? 0
  } finally {
    loading.value = false
  }
}

async function toggleStatus(post: Post) {
  const newStatus = post.status === 'published' ? 'draft' : 'published'
  await postApi.updateStatus(post.id, newStatus)
  post.status = newStatus
}

async function deletePost(id: number) {
  if (!confirm('确定删除这篇文章？')) return
  await postApi.delete(id)
  posts.value = posts.value.filter(p => p.id !== id)
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('zh-CN')
}

onMounted(fetchPosts)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">文章管理</h1>
      <RouterLink to="/admin/posts/new" class="btn-primary">
        <span class="i-carbon-add mr-1" />新建文章
      </RouterLink>
    </div>

    <div class="card overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-400">加载中...</div>
      <table v-else class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-700/50">
          <tr>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300">标题</th>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300 hidden md:table-cell">分类</th>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300">状态</th>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300 hidden sm:table-cell">日期</th>
            <th class="px-4 py-3 text-right font-medium text-gray-600 dark:text-gray-300">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
          <tr v-for="post in posts" :key="post.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/30">
            <td class="px-4 py-3 max-w-xs truncate font-medium text-gray-900 dark:text-gray-100">{{ post.title }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 hidden md:table-cell">{{ post.category_name || '-' }}</td>
            <td class="px-4 py-3">
              <span
                class="text-xs px-2 py-0.5 rounded-full"
                :class="post.status === 'published' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'"
              >
                {{ post.status === 'published' ? '已发布' : '草稿' }}
              </span>
            </td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 text-xs hidden sm:table-cell">
              {{ formatDate(post.updated_at) }}
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-2">
                <RouterLink :to="`/admin/posts/${post.id}/edit`" class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded" title="编辑">
                  <span class="i-carbon-edit text-base" />
                </RouterLink>
                <button class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded" :title="post.status === 'published' ? '下线' : '发布'" @click="toggleStatus(post)">
                  <span :class="post.status === 'published' ? 'i-carbon-cloud-offline' : 'i-carbon-cloud-upload'" class="text-base" />
                </button>
                <button class="p-1.5 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-500 rounded" title="删除" @click="deletePost(post.id)">
                  <span class="i-carbon-trash-can text-base" />
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="posts.length === 0">
            <td colspan="5" class="px-4 py-12 text-center text-gray-400">暂无文章</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div v-if="total > 15" class="flex justify-center gap-2 mt-4">
      <button :disabled="page <= 1" class="btn-secondary disabled:opacity-40" @click="page--; fetchPosts()">上一页</button>
      <span class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ page }} / {{ Math.ceil(total / 15) }}</span>
      <button :disabled="page >= Math.ceil(total / 15)" class="btn-secondary disabled:opacity-40" @click="page++; fetchPosts()">下一页</button>
    </div>
  </div>
</template>
