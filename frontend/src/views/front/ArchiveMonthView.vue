<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { postApi } from '@/api/post'
import type { Post } from '@/types'

const route = useRoute()
const router = useRouter()

const posts = ref<Post[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(true)
const year = ref(0)
const month = ref(0)

async function fetchPosts() {
  loading.value = true
  try {
    const res = await postApi.listByYearMonth(year.value, month.value, page.value, 10) as unknown as { data: { items: Post[]; total: number } }
    posts.value = res.data.items ?? []
    total.value = res.data.total ?? 0
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })
}

onMounted(() => {
  year.value = parseInt(route.params.year as string)
  month.value = parseInt(route.params.month as string)
  fetchPosts()
})

watch(() => route.params, (params) => {
  year.value = parseInt(params.year as string)
  month.value = parseInt(params.month as string)
  page.value = 1
  fetchPosts()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-4 mb-6">
      <button class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg" @click="router.push('/archive')">
        <span class="i-carbon-arrow-left text-xl" />
      </button>
      <h1 class="text-2xl font-bold">{{ year }}年{{ month }}月</h1>
      <span class="text-sm text-gray-500 dark:text-gray-400">共 {{ total }} 篇文章</span>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="card p-6 animate-pulse">
        <div class="h-5 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-3" />
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
      </div>
    </div>

    <div v-else-if="posts.length === 0" class="text-center py-16 text-gray-500 dark:text-gray-400">
      <p>该月份暂无文章</p>
    </div>

    <div v-else class="space-y-4">
      <RouterLink
        v-for="post in posts"
        :key="post.id"
        :to="`/post/${post.slug}`"
        class="card p-6 hover:shadow-md transition-shadow block"
      >
        <h2 class="text-lg font-bold text-gray-900 dark:text-gray-100 mb-2 hover:text-primary-600 dark:hover:text-primary-400">
          {{ post.title }}
        </h2>
        <div class="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
          <time v-if="post.published_at">{{ formatDate(post.published_at) }}</time>
          <span v-if="post.category_name" class="text-primary-600 dark:text-primary-400">{{ post.category_name }}</span>
        </div>
        <p v-if="post.summary" class="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">{{ post.summary }}</p>
      </RouterLink>
    </div>

    <!-- 分页 -->
    <div v-if="total > 10" class="flex justify-center gap-2 mt-8">
      <button
        :disabled="page <= 1"
        class="btn-secondary disabled:opacity-40"
        @click="page--; fetchPosts()"
      >
        上一页
      </button>
      <span class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
        第 {{ page }} 页 / 共 {{ Math.ceil(total / 10) }} 页
      </span>
      <button
        :disabled="page >= Math.ceil(total / 10)"
        class="btn-secondary disabled:opacity-40"
        @click="page++; fetchPosts()"
      >
        下一页
      </button>
    </div>
  </div>
</template>
