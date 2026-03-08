<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { postApi } from '@/api/post'
import { categoryApi } from '@/api/taxonomy'
import type { Post, Category } from '@/types'

const route = useRoute()
const router = useRouter()

const posts = ref<Post[]>([])
const categories = ref<Category[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

const query = ref((route.query.q as string) ?? '')
const selectedCategory = ref((route.query.category as string) ?? '')

async function fetchPosts() {
  loading.value = true
  try {
    const res = await postApi.list({
      page: page.value,
      page_size: 10,
      q: query.value || undefined,
      category: selectedCategory.value || undefined,
    }) as unknown as { data: { items: Post[]; total: number } }
    posts.value = res.data.items ?? []
    total.value = res.data.total ?? 0
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const res = await categoryApi.list() as unknown as { data: Category[] }
    categories.value = res.data ?? []
  } catch {
    categories.value = []
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(() => {
  fetchPosts()
  fetchCategories()
})

watch([page, query, selectedCategory], () => {
  router.replace({ query: { q: query.value || undefined, category: selectedCategory.value || undefined } })
  fetchPosts()
})
</script>

<template>
  <div>
    <!-- 搜索和筛选 -->
    <div class="flex flex-col sm:flex-row gap-3 mb-8">
      <input
        v-model="query"
        type="search"
        placeholder="搜索文章..."
        class="input-base flex-1"
        @keyup.enter="page = 1"
      />
      <select v-model="selectedCategory" class="input-base sm:w-40" @change="page = 1">
        <option value="">所有分类</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.slug">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="card p-6 animate-pulse">
        <div class="h-5 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-3" />
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2" />
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-2/3" />
      </div>
    </div>

    <!-- 文章列表 -->
    <div v-else-if="posts.length > 0" class="space-y-6">
      <article
        v-for="post in posts"
        :key="post.id"
        class="card p-6 hover:shadow-md transition-shadow"
      >
        <!-- 封面图 -->
        <img
          v-if="post.cover_url"
          :src="post.cover_url"
          :alt="post.title"
          class="w-full h-48 object-cover rounded-lg mb-4"
        />

        <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400 mb-2">
          <span v-if="post.category_name" class="bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 px-2 py-0.5 rounded-full">
            {{ post.category_name }}
          </span>
          <time :datetime="post.published_at ?? post.created_at">
            {{ formatDate(post.published_at ?? post.created_at) }}
          </time>
        </div>

        <RouterLink :to="`/post/${post.slug}`" class="block group">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100 group-hover:text-primary-600 dark:group-hover:text-primary-400 mb-2 transition-colors">
            {{ post.title }}
          </h2>
        </RouterLink>

        <p v-if="post.summary" class="text-gray-600 dark:text-gray-400 text-sm leading-relaxed mb-3 line-clamp-2">
          {{ post.summary }}
        </p>

        <div class="flex items-center gap-2 flex-wrap">
          <span
            v-for="tag in post.tags"
            :key="tag.id"
            class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-700 rounded-full text-gray-600 dark:text-gray-300"
          >
            #{{ tag.name }}
          </span>
        </div>
      </article>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-16 text-gray-500 dark:text-gray-400">
      <span class="i-carbon-document-blank text-6xl block mx-auto mb-4" />
      <p>暂无文章</p>
    </div>

    <!-- 分页 -->
    <div v-if="total > 10" class="flex justify-center gap-2 mt-8">
      <button
        :disabled="page <= 1"
        class="btn-secondary disabled:opacity-40"
        @click="page--"
      >
        上一页
      </button>
      <span class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
        第 {{ page }} 页 / 共 {{ Math.ceil(total / 10) }} 页
      </span>
      <button
        :disabled="page >= Math.ceil(total / 10)"
        class="btn-secondary disabled:opacity-40"
        @click="page++"
      >
        下一页
      </button>
    </div>
  </div>
</template>
