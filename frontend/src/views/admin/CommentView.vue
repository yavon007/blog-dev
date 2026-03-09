<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { commentApi } from '@/api/comment'
import type { Comment } from '@/types'

const comments = ref<Comment[]>([])
const total = ref(0)
const page = ref(1)
const statusFilter = ref('')
const loading = ref(false)

async function fetchComments() {
  loading.value = true
  try {
    const res = await commentApi.adminList({
      page: page.value,
      page_size: 20,
      status: statusFilter.value || undefined,
    }) as unknown as { data: { items: Comment[]; total: number } }
    comments.value = res.data.items ?? []
    total.value = res.data.total ?? 0
  } finally {
    loading.value = false
  }
}

async function approve(id: number) {
  await commentApi.updateStatus(id, 'approved')
  const c = comments.value.find(c => c.id === id)
  if (c) c.status = 'approved'
}

async function reject(id: number) {
  await commentApi.updateStatus(id, 'rejected')
  const c = comments.value.find(c => c.id === id)
  if (c) c.status = 'rejected'
}

async function deleteComment(id: number) {
  if (!confirm('确定删除此评论？')) return
  await commentApi.delete(id)
  comments.value = comments.value.filter(c => c.id !== id)
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('zh-CN')
}

const statusLabel: Record<string, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
}

const statusClass: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  approved: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  rejected: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
}

onMounted(fetchComments)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">评论管理</h1>
      <select v-model="statusFilter" class="input-base w-36" @change="page = 1; fetchComments()">
        <option value="">全部状态</option>
        <option value="pending">待审核</option>
        <option value="approved">已通过</option>
        <option value="rejected">已拒绝</option>
      </select>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400">加载中...</div>

    <div v-else class="space-y-3">
      <div v-for="comment in comments" :key="comment.id" class="card p-4">
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <span class="font-medium text-sm">{{ comment.author_name }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ comment.author_email }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(comment.created_at) }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full" :class="statusClass[comment.status]">
                {{ statusLabel[comment.status] }}
              </span>
            </div>
            <p class="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-line">{{ comment.body }}</p>
          </div>

          <div class="flex items-center gap-1 flex-shrink-0">
            <button
              v-if="comment.status !== 'approved'"
              class="p-1.5 hover:bg-green-50 dark:hover:bg-green-900/20 text-green-600 rounded"
              title="通过"
              @click="approve(comment.id)"
            >
              <span class="i-carbon-checkmark text-base" />
            </button>
            <button
              v-if="comment.status !== 'rejected'"
              class="p-1.5 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 text-yellow-600 rounded"
              title="拒绝"
              @click="reject(comment.id)"
            >
              <span class="i-carbon-close text-base" />
            </button>
            <button
              class="p-1.5 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-500 rounded"
              title="删除"
              @click="deleteComment(comment.id)"
            >
              <span class="i-carbon-trash-can text-base" />
            </button>
          </div>
        </div>
      </div>

      <div v-if="comments.length === 0" class="text-center py-12 text-gray-400">
        暂无评论
      </div>
    </div>

    <div v-if="total > 20" class="flex justify-center gap-2 mt-4">
      <button :disabled="page <= 1" class="btn-secondary disabled:opacity-40" @click="page--; fetchComments()">上一页</button>
      <span class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ page }} / {{ Math.ceil(total / 20) }}</span>
      <button :disabled="page >= Math.ceil(total / 20)" class="btn-secondary disabled:opacity-40" @click="page++; fetchComments()">下一页</button>
    </div>
  </div>
</template>
